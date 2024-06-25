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
)

func Create(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "stats"
	root.Short = "System Stats"
	root.Long = "System Stats"
	root.Run = func(cmd *cli.Command, pid int, args []string) {}

	t.AddCommand(root, CreateProfileCPUStart(t))
	t.AddCommand(root, CreateProfileCPUStop(t))
	t.AddCommand(root, CreateProfileMemory(t))
	t.AddCommand(root, CreateMemoryStatus(t))
	t.AddCommand(root, CreateMemoryPlot(t))
	t.AddCommand(root, CreateCPUStatus(t))
	t.AddCommand(root, CreateFdStatus(t))
	t.AddCommand(root, CreateCPUUsage(t))

	return root
}
