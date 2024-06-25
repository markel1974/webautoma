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

func CreateCPUStatus(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "cpu"
	root.Short = "CPUs status"
	root.Long = "CPUs status"
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn("")
		r.WriteLn(fmt.Sprintf("Number of logical CPUs: %d", runtime.NumCPU()))
		r.WriteLn(fmt.Sprintf("Maximum number of CPUs that can be executing simultaneously: %d", runtime.GOMAXPROCS(0)))
		r.WriteLn(fmt.Sprintf("Number of goroutines that currently exist: %d", runtime.NumGoroutine()))
		r.WriteLn(fmt.Sprintf("Number of cgo calls made by the current process: %d", runtime.NumCgoCall()))
	}
	return root
}
