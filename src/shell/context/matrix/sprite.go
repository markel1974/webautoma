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

package matrix

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"strings"
)

type spriteAttributes struct {
	fg  interfaces.ColorDef
	bg  interfaces.ColorDef
	cur rune
}

func (s *spriteAttributes) GetFg() interfaces.ColorDef {
	return s.fg
}
func (s *spriteAttributes) GetBg() interfaces.ColorDef {
	return s.bg
}
func (s *spriteAttributes) GetCur() rune {
	return s.cur
}

type Sprite struct {
	attr [][]spriteAttributes
	size Size
	base rune
}

func NewSprite(raw string, base rune) *Sprite {
	if base == -1 {
		base = '▓'
	}
	s := &Sprite{
		base: base,
	}

	s.setup(raw)
	return s
}

func (s *Sprite) setup(raw string) {
	var fg = interfaces.ColorNoneDef
	var bg = interfaces.ColorNoneDef

	var maxW = 0
	var maxH = 0
	var y = 0

	for _, data := range strings.Split(raw, "\n") {
		if len(data) > 0 {
			var line []spriteAttributes
			for x, c := range data {
				var attr = spriteAttributes{fg: interfaces.ColorNoneDef, bg: interfaces.ColorNoneDef, cur: ' '}
				var valid = true
				switch c {
				case 'w':
					bg = interfaces.ColorWhiteDef
				case 'r':
					bg = interfaces.ColorRedDef
				case 'g':
					bg = interfaces.ColorGreenDef
				case 'y':
					bg = interfaces.ColorYellowDef
				case 'b':
					bg = interfaces.ColorBlueDef
				case 'm':
					bg = interfaces.ColorMagentaDef
				case 'c':
					bg = interfaces.ColorCyanDef
				case 'k':
					bg = interfaces.ColorBlackDef
				case 'W':
					fg = interfaces.ColorWhiteDef
				case 'R':
					fg = interfaces.ColorRedDef
				case 'G':
					fg = interfaces.ColorGreenDef
				case 'Y':
					fg = interfaces.ColorYellowDef
				case 'B':
					fg = interfaces.ColorBlueDef
				case 'M':
					fg = interfaces.ColorMagentaDef
				case 'C':
					fg = interfaces.ColorCyanDef
				case 'K':
					fg = interfaces.ColorBlackDef
				default:
					valid = false
				}
				if valid {
					attr.fg = fg
					attr.bg = bg
					attr.cur = s.base
				}
				line = append(line, attr)
				if x > maxW {
					maxW = x
				}
			}
			s.attr = append(s.attr, line)
			if y > maxH {
				maxH = y
			}
			y++
		}
	}

	s.size.w = float64(maxW + 1)
	s.size.h = float64(maxH + 1)
}
