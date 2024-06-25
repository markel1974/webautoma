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

type IContext interface {
	Write(data string)
	WriteLn(data string)
	WriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)
	WriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)
	GetScreenSize() (int, int)
	SetContext(pid int, ctx interface{}) bool
	IsActive(pid int) bool
	CreateTimer(pid int, first int, interval int, count int) bool
	StopTimer(pid int, tid int) bool
	PaintRequest(pid int) bool
	SetCaption(pid int, caption string) bool
	SetBasePath(arg string)
	SetSelectionMode(int)
	SetSelectionOptions(option rune, value float64) bool
	SetSelectionModeNext()
	SetSelectionModePrevious()
	Deactivate(pid int) bool
	DeactivateAll(name string) int
	History(verb HistoryAction, idx int)
	ClearScreen()
	TaskList() string
	SaveTasks(name string) bool
	RestoreTasks(name string) bool
	ListTasks() []string
	SetExit()
	SetFg(pid int) bool
}
