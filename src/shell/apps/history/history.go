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

func Create(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "history"
	root.Short = "History"
	root.Long = "History"
	root.Aliases = []string{"h"}
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		idx := -1
		if len(args) > 0 {
			if i, err := strconv.Atoi(args[0]); err == nil {
				idx = i
			}
		}
		if idx > -1 {
			r.History(interfaces.HistoryActionExec, idx)
		} else {
			r.History(interfaces.HistoryActionList, -1)
		}
	}

	t.AddCommand(root, CreateHistoryClear(t))
	t.AddCommand(root, CreateHistoryExec(t))
	t.AddCommand(root, CreateHistoryList(t))

	return root
}
