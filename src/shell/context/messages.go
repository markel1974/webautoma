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

package context

type MessageType int

const (
	MessageTypeRead  MessageType = iota
	MessageTypeTimer MessageType = iota
	MessageTypePaint MessageType = iota
	MessageTypeQuit  MessageType = iota
)

type iMessage interface {
	getType() MessageType
	postEvent(chan iMessage)
}

type MessageRead struct {
	data []byte
}

func newMessageRead(data []byte, n int) iMessage {
	if n > len(data) {
		n = len(data) - 1
	}
	x := data[:n]
	return &MessageRead{data: x}
}
func (m *MessageRead) getType() MessageType {
	return MessageTypeRead
}
func (m *MessageRead) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}

type MessageTimer struct {
	pid      int
	tid      int
	interval int
}

func newMessageTimer(pid int, interval int) *MessageTimer {
	return &MessageTimer{pid: pid, interval: interval}
}
func (m *MessageTimer) getType() MessageType {
	return MessageTypeTimer
}
func (m *MessageTimer) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}

type MessageQuit struct {
}

func newMessageQuit() *MessageQuit {
	return &MessageQuit{}
}
func (m *MessageQuit) getType() MessageType {
	return MessageTypeQuit
}
func (m *MessageQuit) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}

type MessagePaint struct {
}

func newMessagePaint() *MessagePaint {
	return &MessagePaint{}
}
func (m *MessagePaint) getType() MessageType {
	return MessageTypePaint
}
func (m *MessagePaint) postEvent(ch chan iMessage) {
	go func() { ch <- m }()
}
