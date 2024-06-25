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

package interfaces

type KeyType int

const (
	KeyTypeKey KeyType = iota
	KeyTypeCursor
	KeyTypeCtrl
	KeyTypeEnter
	KeyTypeTab
	KeyTypeBackspace
	KeyTypeCancel
)

const (
	CursorUpDef    CursorCodeDef = iota
	CursorDownDef  CursorCodeDef = iota
	CursorLeftDef  CursorCodeDef = iota
	CursorRightDef CursorCodeDef = iota
)

type KeyData struct {
	Type KeyType
	Key  rune
}

func NewKeyData(t KeyType, key rune) *KeyData {
	return &KeyData{
		Type: t,
		Key:  key,
	}
}

func (e *KeyData) GetType() KeyType {
	return e.Type
}
