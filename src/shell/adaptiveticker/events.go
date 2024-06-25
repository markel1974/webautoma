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

type eventType int

const (
	eventTypeCreate eventType = iota
	eventTypeRemove eventType = iota
	eventTypeExpire eventType = iota
	eventTypeQuit   eventType = iota
)

type iEvent interface {
	GetType() eventType
	PostEvent(ch chan iEvent)
}

type createEvent struct {
	handler *TimerHandler
}

func newCreateEvent(handler *TimerHandler) *createEvent {
	return &createEvent{handler: handler}
}
func (b *createEvent) GetType() eventType {
	return eventTypeCreate
}
func (b *createEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}

type removeEvent struct {
	tids []int
}

func newRemoveEvent(tids []int) *removeEvent {
	return &removeEvent{tids: tids}
}
func (b *removeEvent) GetType() eventType {
	return eventTypeRemove
}
func (b *removeEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}

type expireEvent struct {
}

func newExpireEvent() *expireEvent {
	return &expireEvent{}
}
func (b *expireEvent) GetType() eventType {
	return eventTypeExpire
}
func (b *expireEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}

type quitEvent struct {
}

func newQuitEvent() *quitEvent {
	return &quitEvent{}
}
func (b *quitEvent) GetType() eventType {
	return eventTypeQuit
}
func (b *quitEvent) PostEvent(ch chan iEvent) {
	go func() { ch <- b }()
}
