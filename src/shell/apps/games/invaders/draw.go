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

package invaders

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
)

func drawColor(surface interfaces.ISurface, x, y int, fg, bg interfaces.ColorDef, data string) {
	for _, c := range data {
		surface.DrawColor(y, x, c, fg, bg, interfaces.ModeNormal)
		x++
	}
}
