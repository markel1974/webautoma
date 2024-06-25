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

import "container/list"

const UnknownId = -1

type IIds interface {
	SetId(int)
	Unset()
}

type Ids struct {
	slots       []bool
	currentSlot int
	kv          map[int]*list.Element
	ll          *list.List
}

func NewIds(max int) *Ids {
	a := &Ids{
		slots:       make([]bool, max),
		kv:          make(map[int]*list.Element),
		ll:          list.New(),
		currentSlot: 0,
	}
	return a
}

func (a *Ids) Set(obj IIds) int {
	var id = UnknownId

	if a.currentSlot >= len(a.slots) {
		a.currentSlot = 0
	}

	_, exists := a.kv[a.currentSlot]
	if !exists {
		id = a.currentSlot
	} else {
		for slot := 0; slot < len(a.slots); slot++ {
			if !a.slots[slot] {
				id = slot
				break
			}
		}
	}

	if id != UnknownId {
		a.slots[id] = true
		element := a.ll.PushBack(obj)
		a.kv[id] = element
	}

	obj.SetId(id)

	a.currentSlot++

	return id
}

func (a *Ids) Get(id int) (IIds, bool) {
	var obj IIds = nil
	var e, ok = a.kv[id]
	if ok {
		obj = e.Value.(IIds)
	}
	return obj, ok
}

func (a *Ids) All() *list.List {
	return a.ll
}

func (a *Ids) Unset(id int) bool {
	if id < 0 {
		return false
	}
	if id >= len(a.slots) {
		return false
	}
	var element, ok = a.kv[id]
	if !ok {
		return false
	}

	var obj = element.Value.(IIds)
	obj.Unset()

	a.ll.Remove(element)
	delete(a.kv, id)
	a.slots[id] = false

	return true
}
