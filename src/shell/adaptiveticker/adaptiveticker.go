/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package adaptiveticker

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

const (
	maxEventsPerSecond        = 10
	adaptiveTickerMaxQueueLen = 1024
)

type AdaptiveTicker struct {
	container   list.List
	timer       *time.Timer
	runDeadline int64
	minInterval int64
	ids         *Ids
	index       map[int64]*TimerGroupHandler

	lock     sync.Mutex
	messages chan iEvent
	quit     bool

	statsEpsSecond    int64
	statsExpiredTotal int64
	statsExpired      int64
}

func NewAdaptiveTicker() *AdaptiveTicker {
	a := &AdaptiveTicker{
		timer:       nil,
		runDeadline: 0,
		ids:         NewIds(8192),
		minInterval: 1000 / maxEventsPerSecond,
		index:       make(map[int64]*TimerGroupHandler),
		messages:    make(chan iEvent, adaptiveTickerMaxQueueLen),
		quit:        false,
	}
	a.container.Init()

	a.eventLoop()

	return a
}

func (a *AdaptiveTicker) Create(target chan *TimerHandler, event interface{}, first int64, interval int64, count int64) int {
	var current = NewTimerHandler(target, event, first, interval, count, a.minInterval)

	a.lock.Lock()
	if !a.quit {
		a.ids.Set(current)
	}
	a.lock.Unlock()

	if current.id != UnknownId {
		//fmt.Println("create - trying go send message ....", current.id, len(a.messages))
		//debug.PrintStack()
		var event = newCreateEvent(current)
		event.PostEvent(a.messages)
		//fmt.Println("create - message sent", current.id, len(a.messages))
	}

	return current.id
}

func (a *AdaptiveTicker) Remove(tids []int) bool {
	var removed []int

	a.lock.Lock()
	if !a.quit {
		for _, tid := range tids {
			if a.ids.Unset(tid) {
				removed = append(removed, tid)
			}
		}
	}
	a.lock.Unlock()

	if len(removed) > 0 {
		var event = newRemoveEvent(removed)
		event.PostEvent(a.messages)
	}

	return len(tids) == len(removed)
}

func (a *AdaptiveTicker) Quit() {
	var emit = false
	a.lock.Lock()
	if !a.quit {
		a.quit = true
		emit = true
	}
	a.lock.Unlock()

	if emit {
		var event = newQuitEvent()
		event.PostEvent(a.messages)
	}
}

func (a *AdaptiveTicker) doExpire(now int64) []*TimerHandler {
	var expired []*TimerHandler

	a.runDeadline = 0

	var next *list.Element
	for e := a.container.Front(); e != nil; e = next {
		next = e.Next()
		group := e.Value.(*TimerGroupHandler)
		if group.deadline <= now {
			for _, current := range group.container {
				if !current.removed {
					expired = append(expired, current)
				}
			}

			delete(a.index, group.deadline)
			a.container.Remove(e)
		} else {
			break
		}
	}

	for _, current := range expired {
		if current.IsUsable() {
			a.doAdd(now, current)
		} else {
			a.lock.Lock()
			a.ids.Unset(current.id)
			a.lock.Unlock()
		}
	}

	return expired
}

func (a *AdaptiveTicker) doHeadRun(now int64) {
	if a.container.Len() == 0 {
		return
	}

	if a.runDeadline > 0 {
		return
	}

	var target = a.container.Front().Value.(*TimerGroupHandler)
	a.runDeadline = target.deadline

	interval := a.runDeadline - now
	if interval < a.minInterval {
		interval = a.minInterval
	}

	d := time.Duration(interval) * time.Millisecond

	if a.timer != nil {
		if started := a.timer.Reset(d); !started {
			a.timer = nil
		}
	}

	if a.timer == nil {
		a.timer = time.NewTimer(d)
		go func(t *time.Timer) {
			select {
			case <-t.C:
				var event = newExpireEvent()
				event.PostEvent(a.messages)
			}
		}(a.timer)
	}
}

func (a *AdaptiveTicker) doAdd(now int64, event *TimerHandler) {
	event.Prepare(now)

	if group, ok := a.index[event.deadline]; ok {
		group.Add(event)
		//fmt.Println("ADDING", event.id, event.interval, event.deadline, "- ITERATION", 0)
		return
	}

	var newGroup = NewGroupEvent(event)

	a.index[newGroup.deadline] = newGroup

	if a.container.Len() == 0 {
		a.container.PushFront(newGroup)
		a.runDeadline = 0
		//fmt.Println("ADDING", event.id, event.interval, event.deadline, "^ ITERATION", 0)
		return
	}

	var head = a.container.Front()
	var g = head.Value.(*TimerGroupHandler)
	if newGroup.deadline < g.deadline {
		if a.runDeadline > 0 {
			if newGroup.deadline < a.runDeadline {
				if diff := a.runDeadline - newGroup.deadline; diff > a.minInterval {
					a.runDeadline = 0
				}
			}
		}

		a.container.InsertBefore(newGroup, head)
		//fmt.Println("ADDING", event.id, event.interval, event.deadline, ". ITERATION", 0)
		return
	}
	//var count = 0
	for elem := head.Next(); elem != nil; elem = elem.Next() {
		//count ++
		g = elem.Value.(*TimerGroupHandler)
		if newGroup.deadline < g.deadline {
			a.container.InsertBefore(newGroup, elem)
			//fmt.Println("ADDING", event.id, event.interval, event.deadline, "+ ITERATION", count)
			return
		}
	}

	a.container.PushBack(newGroup)
	//fmt.Println("ADDING", event.id, event.interval, event.deadline, "$ ITERATION", count)
}

func (a *AdaptiveTicker) doQuit() {
	if a.timer != nil {
		a.timer.Stop()
		a.timer = nil
	}
	a.runDeadline = 0
	a.container.Init()
}

func (a *AdaptiveTicker) eventLoop() {
	go func() {
		for {
			select {
			case msg := <-a.messages:
				switch msg.GetType() {
				case eventTypeCreate:
					a.createEventHandler(msg.(*createEvent))
				case eventTypeRemove:
					a.removeEventHandler()
				case eventTypeExpire:
					a.expireEventHandler()
				case eventTypeQuit:
					a.quitEventHandler()
					return
				}
			}
		}
	}()
}

func (a *AdaptiveTicker) createEventHandler(event *createEvent) {
	var now = getEpochMs()
	a.doAdd(now, event.handler)
	a.doHeadRun(now)
}

func (a *AdaptiveTicker) removeEventHandler() {
	var now = getEpochMs()
	a.doHeadRun(now)
}

func (a *AdaptiveTicker) expireEventHandler() {
	var now = getEpochMs()
	var expired = a.doExpire(now)
	a.doHeadRun(now)

	//a.computeEps(now, int64(len(expired)))

	for _, current := range expired {
		current.PostEvent()
	}
}

func (a *AdaptiveTicker) quitEventHandler() {
	//var now = getEpochMs()
	a.doQuit()
}

func (a *AdaptiveTicker) computeEps(now int64, expired int64) {
	var epsSecond = now / 1000
	if epsSecond == a.statsEpsSecond {
		a.statsExpiredTotal += expired
		a.statsExpired++
	} else {
		fmt.Println("expired", a.statsExpired, "total", a.statsExpiredTotal)
		a.statsEpsSecond = epsSecond
		a.statsExpired = 1
		a.statsExpiredTotal = expired
	}
}

//func (a* AdaptiveTicker) doPrintContainer(now int64) {
//	var counter = 0
//	for e := a.container.Front(); e != nil; e = e.Next() {
//		if counter > 0 {
//			fmt.Print(", ")
//		}
//		group := e.Value.(*TimerGroupHandler)
//		for _, current := range group.container {
//			fmt.Print(current.id, " (", current.deadline, ", ", current.deadline- now, ")")
//		}
//		counter ++
//	}
//	fmt.Println()
//}

/*

import (
	"time"
	"container/list"
	"fmt"
)

var __fps int64 = 0
var __fpsSecond int64 = 0
var __expired int64 = 0


type AdaptiveTicker struct {
	container    list.List
	timer        *time.Timer
	runDeadline  int64
	session      *Session
	minInterval  int64
	ids          *Ids
	index        map[int64]*TimerGroupHandler
}

func NewAdaptiveTicker() * AdaptiveTicker {
	a := &AdaptiveTicker{
		timer:       nil,
		runDeadline: 0,
		ids:         NewIds(8192),
		minInterval: 100,
		index:       make(map[int64]*TimerGroupHandler),
	}

	a.session = NewSession(nil, []func(){ a.doHeadRun })
	a.container.Init()
	return a
}

func (a* AdaptiveTicker) Create(target chan *TimerHandler, event interface{}, first int64, interval int64, count int64) int {
	a.session.Acquire()

	var current = NewTimerHandler(target, event, first, interval, count, a.minInterval)

	a.ids.Set(current)
	if current.id != UnknownId {
		a.doAdd(current)
	}

	a.session.Release()

	return current.id
}

func (a* AdaptiveTicker) Remove(tids []int) bool {
	var count = 0

	a.session.Acquire()

	for _, tid := range tids {
		if a.ids.Unset(tid) {
			count++
		}
	}

	a.session.Release()

	return len(tids) == count
}

func (a* AdaptiveTicker) doExpire() {
	//MUST BE CALLED FROM A GO FUNCTION!
	var expired []*TimerHandler

	a.session.Acquire()

	a.runDeadline = 0

	//a.doPrintContainer()

	var next *list.Element
	for e := a.container.Front(); e != nil; e = next {
		next = e.Next()
		group := e.Value.(*TimerGroupHandler)
		if group.deadline <= a.session.Now() {
			for _, current := range group.container {
				if !current.removed {
					expired = append(expired, current)
				}
			}
			delete(a.index, group.deadline)
			a.container.Remove(e)
		} else {
			break
		}
	}

	for _, current := range expired {
		if current.IsUsable() {
			a.doAdd(current)
		} else {
			a.ids.Unset(current.id)
		}
	}

	a.computeEps(a.session.now, int64(len(expired)))

	a.session.Release()

	for _, current := range expired {
		current.target <- current
	}
}

func (a* AdaptiveTicker) Stop() {
	a.session.Acquire()

	if a.timer != nil {
		a.timer.Stop()
		a.timer = nil
	}
	a.runDeadline = 0
	a.container.Init()

	a.session.Release()
}

func (a* AdaptiveTicker) doHeadRun() {
	if a.container.Len() == 0 {
		return
	}

	if a.runDeadline > 0 {
		return
	}

	var target = a.container.Front().Value.(*TimerGroupHandler)
	a.runDeadline = target.deadline

	interval := a.runDeadline - a.session.Now()
	if interval < a.minInterval {
		interval = a.minInterval
	}

	d := time.Duration(interval) * time.Millisecond

	if a.timer != nil {
		if started := a.timer.Reset(d); !started {
			a.timer = nil
		}
	}

	if a.timer == nil {
		a.timer = time.NewTimer(d)

		go func(t * time.Timer) {
			select {
				case <-t.C: a.doExpire()
			}
		}(a.timer)
	}
}

func (a* AdaptiveTicker) doAdd(event * TimerHandler) {
	event.Prepare(a.session.Now())

	if group, ok := a.index[event.deadline]; ok {
		group.Add(event)
		//fmt.Println("ADDING", event.id, event.interval, event.deadline, "- ITERATION", 0)
		return
	}

	var newGroup = NewGroupEvent(event)

	a.index[newGroup.deadline] = newGroup

	if a.container.Len() == 0 {
		a.container.PushFront(newGroup)
		a.runDeadline = 0
		//fmt.Println("ADDING", event.id, event.interval, event.deadline, "^ ITERATION", 0)
		return
	}

	var head = a.container.Front()
	var g = head.Value.(*TimerGroupHandler)
	if newGroup.deadline < g.deadline {
		if a.runDeadline > 0 {
			if newGroup.deadline < a.runDeadline {
				if diff := a.runDeadline - newGroup.deadline; diff > a.minInterval {
					a.runDeadline = 0
				}
			}
		}

		a.container.InsertBefore(newGroup, head)
		//fmt.Println("ADDING", event.id, event.interval, event.deadline, ". ITERATION", 0)
		return
	}

	//var count = 0

	for elem := head.Next(); elem != nil; elem = elem.Next() {
		//count ++
		g = elem.Value.(*TimerGroupHandler)
		if newGroup.deadline < g.deadline {
			a.container.InsertBefore(newGroup, elem)
			//fmt.Println("ADDING", event.id, event.interval, event.deadline, "+ ITERATION", count)
			return
		}
	}

	a.container.PushBack(newGroup)
	//fmt.Println("ADDING", event.id, event.interval, event.deadline, "$ ITERATION", count)
}

func (a* AdaptiveTicker) computeEps(now int64, expired int64) {
	var fpsSecond  = now / 1000
	if fpsSecond  == __fpsSecond {
		__fps ++
		__expired += expired
	} else {
		fmt.Println("EPS", __fps, fpsSecond - __fpsSecond, __expired)
		__fpsSecond = fpsSecond
		__fps = 1
		__expired = 0
	}
}

//func (a* AdaptiveTicker) doPrintContainer() {
//	var counter = 0
//	for e := a.container.Front(); e != nil; e = e.Next() {
//		if counter > 0 {
//			fmt.Print(", ")
//		}
//		group := e.Value.(*TimerGroupHandler)
//		for _, current := range group.container {
//			fmt.Print(current.id, " (", current.deadline, ", ", current.deadline-a.session.Now(), ")")
//		}
//		counter ++
//	}
//	fmt.Println()
//}

*/
