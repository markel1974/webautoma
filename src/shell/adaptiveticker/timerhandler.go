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

type TimerMode int

const (
	TimerModeContinuous TimerMode = iota
	TimerModeCounter    TimerMode = iota
)

type TimerHandler struct {
	id        int
	target    chan *TimerHandler
	Event     interface{}
	first     int64
	interval  int64
	mode      TimerMode
	deadline  int64
	min       int64
	counter   int64
	loopCount int64
	removed   bool
}

func rounding(val int64, r int64) int64 {
	if val < r {
		return r
	}
	return (val / r) * r
}

func NewTimerHandler(target chan *TimerHandler, event interface{}, first int64, interval int64, loopCount int64, min int64) *TimerHandler {
	var mode TimerMode

	if loopCount <= 0 {
		mode = TimerModeContinuous
	} else {
		mode = TimerModeCounter
	}

	t := &TimerHandler{
		target:    target,
		Event:     event,
		mode:      mode,
		min:       min,
		loopCount: loopCount,
		first:     rounding(first, min),
		interval:  rounding(interval, min),
		deadline:  0,
		counter:   0,
		removed:   false,
	}

	return t
}

func (t *TimerHandler) Prepare(now int64) {
	interval := t.interval
	if t.counter == 0 {
		interval = t.first
	}
	t.deadline = rounding(now+interval, t.min)
	t.counter++
}

func (t *TimerHandler) SetId(id int) {
	t.id = id
}

func (t *TimerHandler) Unset() {
	t.removed = true
}

func (t *TimerHandler) IsUsable() bool {
	var ret bool
	switch t.mode {
	case TimerModeCounter:
		if t.counter == t.loopCount {
			ret = false
		}
	case TimerModeContinuous:
		ret = true
	default:
		ret = false
	}
	return ret
}

func (t *TimerHandler) PostEvent() {
	go func() {
		t.target <- t
	}()
}
