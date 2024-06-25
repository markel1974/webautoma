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
	"github.com/markel1974/webautoma/src/shell/apps/commandcreator"
	"github.com/markel1974/webautoma/src/shell/cli"
	"os"
	"runtime/pprof"
)

func CreateProfileCPUStart(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "startcpuprofile"
	root.Short = "Start cpu profiling"
	root.Long = "Start cpu profiling"
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()

		r.WriteLn("")

		if len(args) <= 0 {
			r.WriteLn("could not create cpu profile: " + "missing filename")
			return
		}

		f, err := os.Create(args[0])
		if err != nil {
			r.WriteLn("could not create CPU profile: " + err.Error())
			return
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			r.WriteLn("could not start CPU profile: " + err.Error())
			return
		}

		r.WriteLn("Cpu Profiling started")
	}

	return root
}
