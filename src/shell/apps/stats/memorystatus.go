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

package stats

import (
	"fmt"
	"github.com/markel1974/webautoma/src/shell/apps/commandcreator"
	"github.com/markel1974/webautoma/src/shell/cli"
	"runtime"
)

func CreateMemoryStatus(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "rt"
	root.Short = "Runtime Status"
	root.Long = "Runtime Status"
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		var m runtime.MemStats

		runtime.ReadMemStats(&m)
		r.WriteLn("")
		r.WriteLn(fmt.Sprintf("Allocated memory in heap objects: %.3f MB", bToMb(m.Alloc)))
		r.WriteLn(fmt.Sprintf("Total memory allocated for heap objects: %.3f MB", bToMb(m.TotalAlloc)))
		r.WriteLn(fmt.Sprintf("Total memory obtained from the OS: %.3f MB", bToMb(m.Sys)))
		r.WriteLn(fmt.Sprintf("Number of completed GC cycles: %d", m.NumGC))
	}
	return root
}
