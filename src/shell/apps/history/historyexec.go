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

package history

import (
	"github.com/markel1974/webautoma/src/shell/apps/commandcreator"
	"github.com/markel1974/webautoma/src/shell/cli"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"strconv"
)

func CreateHistoryExec(t commandcreator.ICreator) *cli.Command {
	execute := t.CreateCommand()
	execute.Use = "exec"
	execute.Short = "Exec"
	execute.Long = "Exec"
	execute.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		if len(args) > 0 {
			if idx, err := strconv.Atoi(args[0]); err == nil {
				r.History(interfaces.HistoryActionExec, idx)
			}
		}
	}
	return execute
}
