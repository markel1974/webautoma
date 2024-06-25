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

package games

import (
	"github.com/markel1974/webautoma/src/shell/apps/commandcreator"
	"github.com/markel1974/webautoma/src/shell/apps/games/tetris"
	"github.com/markel1974/webautoma/src/shell/cli"
	"github.com/markel1974/webautoma/src/shell/interfaces"
)

func CreateTetris(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "tetris"
	root.Short = "Tetris"
	root.Long = "Tetris"
	root.Activate = true
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		//w, h := r.GetScreenSize()

		t := tetris.New(10, 18)
		//s.SetSize(h, w)
		r.SetContext(pid, t)
		r.CreateTimer(pid, 0, 300, -1)
	}
	root.ReadEvent = func(cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
		t := ctx.(*tetris.Tetris)
		switch key {
		case 'a':
			t.MoveLeft()
		case 'd':
			t.MoveRight()
		case 'w':
			t.RotateRight()
		case 's':
			t.MoveDown()
		case ' ':
			t.Drop()
		case '1':
			//r := cmd.GetRootContext()
			//w, h := r.GetScreenSize()
			t.Init(10, 18)
		}
	}
	root.TimerEvent = func(cmd *cli.Command, pid int, tid int, ctx interface{}, interval int) {
		r := cmd.GetRootContext()
		t := ctx.(*tetris.Tetris)
		t.ApplyGravity()
		r.PaintRequest(pid)
	}
	root.PaintEvent = func(cmd *cli.Command, pid int, ctx interface{}, surface interfaces.ISurface) {
		t := ctx.(*tetris.Tetris)

		//rows, columns := surface.GetSize()

		//h, w := t.GetSize()
		//if h != rows || w != columns {
		//	fmt.Println(h, w)
		//	t.Init(rows, columns)
		//}

		t.Draw(surface)
	}

	return root
}
