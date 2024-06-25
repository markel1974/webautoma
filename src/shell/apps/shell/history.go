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

package shell

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var historySaveLock sync.Mutex

type HistoryHandler struct {
	Queue    []string
	queuePos int
	max      uint
	def      string
	enabled  bool
	autosave bool
}

func NewHistoryHandler(max uint, autosave bool) *HistoryHandler {
	h := &HistoryHandler{
		max:      max,
		enabled:  true,
		autosave: autosave,
	}
	h.Clear()

	if h.autosave {
		h.restore()
	}
	return h
}

func (h *HistoryHandler) save() {
	if out, err := json.Marshal(h); err == nil {
		historySaveLock.Lock()
		ioutil.WriteFile("history.json", out, 0644)
		historySaveLock.Unlock()
	}
}

func (h *HistoryHandler) restore() {
	historySaveLock.Lock()
	body, err := ioutil.ReadFile("history.json")
	historySaveLock.Unlock()

	if err == nil {
		h.Clear()
		json.Unmarshal(body, h)
		h.queuePos = len(h.Queue) - 1
		h.def = ""
	}
}

func (h *HistoryHandler) Clear() {
	h.Queue = nil
	h.Queue = append(h.Queue, "")
	h.queuePos = 0
	h.def = ""
}

func (h *HistoryHandler) SetEnabled(enabled bool) {
	h.enabled = enabled
}

func (h *HistoryHandler) AddToHistory(data string) {
	if h.enabled {
		h.Queue[len(h.Queue)-1] = data

		h.Queue = append(h.Queue, "")
		if len(h.Queue) > int(h.max) {
			h.Queue = h.Queue[1:]
		}

		h.queuePos = len(h.Queue) - 1

		if h.autosave {
			h.save()
		}
	}
}

func (h *HistoryHandler) getHistoryAtIndex(idx int) string {
	var out string
	if h.enabled {
		if idx >= 0 && idx < len(h.Queue)-1 {
			out = h.Queue[idx]
		}
	}
	return out
}

func (h *HistoryHandler) GetHistoryPrev() (string, bool) {
	if !h.enabled {
		return "", false
	}

	if h.queuePos == 0 {
		return "", false
	}

	h.queuePos--
	data := h.getHistoryAtIndex(h.queuePos)

	return data, true
}

func (h *HistoryHandler) GetHistory() []string {
	var out []string
	l := len(h.Queue) - 1
	if l > 0 {
		for x := 0; x < l; x++ {
			out = append(out, h.Queue[x])
		}
	}
	return out
}

func (h *HistoryHandler) GetHistoryAtPos(pos int) (string, bool) {
	var out string
	found := false
	l := len(h.Queue) - 1
	if l > 0 && pos < l {
		out = h.Queue[pos]
		if len(out) > 0 {
			found = true
		}
	}
	return out, found
}

func (h *HistoryHandler) GetHistoryNext() (string, bool) {
	if !h.enabled {
		return "", false
	}

	var data string
	maxPos := len(h.Queue) - 1

	if h.queuePos == maxPos {
		return "", false
	}

	if h.queuePos < maxPos-1 {
		h.queuePos++
		data = h.getHistoryAtIndex(h.queuePos)
	} else {
		h.queuePos++
		data = h.def
	}

	return data, true
}

func (h *HistoryHandler) SetDefault(def string) {
	h.def = def
}
