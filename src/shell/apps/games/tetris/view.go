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

package tetris

import (
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"strings"
)

const (
	srcBackground = `
		WWWWWWWWWWWW WWWWWW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WWWWWWWWWWWW
	`

	boardXOffset    = 3
	boardYOffset    = 2
	nextMinoXOffset = 16
	nextMinoYOffset = 2
	blankColor      = interfaces.ColorBlackDef
)

var (
	background = strings.Split(srcBackground, "\n")

	colorMapping = map[rune]interfaces.ColorDef{
		'k': interfaces.ColorBlackDef,
		'K': interfaces.ColorBrightBlackDef,
		'r': interfaces.ColorRedDef,
		'R': interfaces.ColorBrightRedDef,
		'g': interfaces.ColorGreenDef,
		'G': interfaces.ColorBrightGreenDef,
		'y': interfaces.ColorYellowDef,
		'Y': interfaces.ColorBrightYellowDef,
		'b': interfaces.ColorBlueDef,
		'B': interfaces.ColorBrightBlueDef,
		'm': interfaces.ColorMagentaDef,
		'M': interfaces.ColorBrightMagentaDef,
		'c': interfaces.ColorCyanDef,
		'C': interfaces.ColorBrightCyanDef,
		'w': interfaces.ColorWhiteDef,
		'W': interfaces.ColorWhiteDef,
	}
)

var (
	availableColors = []interfaces.ColorDef{
		interfaces.ColorRedDef,
		interfaces.ColorGreenDef,
		interfaces.ColorYellowDef,
		interfaces.ColorBlueDef,
		interfaces.ColorMagentaDef,
		interfaces.ColorCyanDef,
		interfaces.ColorWhiteDef,
		interfaces.ColorBrightRedDef,
		interfaces.ColorBrightGreenDef,
		interfaces.ColorBrightYellowDef,
		interfaces.ColorBrightBlueDef,
		interfaces.ColorBrightMagentaDef,
		interfaces.ColorBrightCyanDef,
		interfaces.ColorBrightWhiteDef,
		interfaces.ColorBrightBlackDef,
	}
)

func colorByChar(ch rune) interfaces.ColorDef {
	return colorMapping[ch]
}

// charByColor
func _(color interfaces.ColorDef) rune {
	for ch, currentColor := range colorMapping {
		if currentColor == color {
			return ch
		}
	}
	return '.'
}
