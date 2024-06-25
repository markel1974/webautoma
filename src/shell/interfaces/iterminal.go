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

package interfaces

type CursorCodeDef rune
type ColorDef int
type ColorMode int

const (
	ModeNormal ColorMode = iota
	Mode8bit   ColorMode = iota
)

const (
	ColorNoneDef          ColorDef = iota
	ColorRedDef           ColorDef = iota
	ColorGreenDef         ColorDef = iota
	ColorYellowDef        ColorDef = iota
	ColorBlueDef          ColorDef = iota
	ColorMagentaDef       ColorDef = iota
	ColorCyanDef          ColorDef = iota
	ColorWhiteDef         ColorDef = iota
	ColorBlackDef         ColorDef = iota
	ColorBrightRedDef     ColorDef = iota
	ColorBrightGreenDef   ColorDef = iota
	ColorBrightYellowDef  ColorDef = iota
	ColorBrightBlueDef    ColorDef = iota
	ColorBrightMagentaDef ColorDef = iota
	ColorBrightCyanDef    ColorDef = iota
	ColorBrightWhiteDef   ColorDef = iota
	ColorBrightBlackDef   ColorDef = iota
	ColorGrayDef          ColorDef = iota
	ColorNormalDef        ColorDef = iota
)

type KeyFunc func(event *KeyData)

type ITerminal interface {
	SetKeyFunc(k KeyFunc)

	SetEnterKey(key rune)

	Colorize(text string, fg int, bg int, mode ColorMode) string

	WriteColor(text string, fg ColorDef, bg ColorDef, mode ColorMode) (int, error)

	Write(text string) (int, error)

	SaveCursor() (int, error)

	RestoreCursor() (int, error)

	MoveCursorLeft() (int, error)

	MoveCursorRight() (int, error)

	MoveCursorTopLeft() (int, error)

	ClearLine(line string) (int, error)

	ClearScreen() (int, error)

	SetSize(w int, h int)

	Scan(data []byte)
}
