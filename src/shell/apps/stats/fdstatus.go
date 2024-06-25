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
	"os"
	"os/exec"
	"strings"
)

func CreateFdStatus(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "fd"
	root.Short = "FD Status"
	root.Long = "FD Status"
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn("")
		out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			r.WriteLn(fmt.Sprintf("Opened fd: %d", len(lines)-1))
		}
	}
	return root
}
