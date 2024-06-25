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

package vt100

import (
	"fmt"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"io"
	"log"
)

//var controlSequenceIntroducer = []byte{ 27, 91 }

var (
	escMoveCursorLeftDef    = []byte{27, 91, 68}
	escMoveCursorRightDef   = []byte{27, 91, 67}
	escMoveCursorTopLeftDef = []byte{27, 91, 'H'}
	escSaveCursorDef        = []byte{27, '7'}
	escRestoreCursorDef     = []byte{27, '8'}

	//	escClearLineDef =   	  []byte{ 27, 91, '2', 'K' }
)

type VT100 struct {
	z        io.Writer
	keyFunc  interfaces.KeyFunc
	debug    bool
	enterKey rune
}

func NewVt100(z io.Writer, debug bool) *VT100 {
	return &VT100{
		z:        z,
		debug:    debug,
		enterKey: 13,
	}
}

func (l *VT100) Write(text string) (int, error) {
	return l.z.Write([]byte(text))
}

func (l *VT100) WriteColor(text string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) (int, error) {
	return l.z.Write([]byte(l.Colorize(text, int(fg), int(bg), mode)))
}

func (l *VT100) Colorize(text string, f int, b int, mode interfaces.ColorMode) string {
	switch mode {
	case interfaces.ModeNormal:
		fg := interfaces.ColorDef(f)
		bg := interfaces.ColorDef(b)
		if fg == interfaces.ColorNoneDef && bg == interfaces.ColorNoneDef {
			return text
		}
		fgColor := l.colorFgAdapter(fg)
		if bg == interfaces.ColorNoneDef {
			return Colorize(text, fgColor)
		}
		bgColor := l.colorBgAdapter(bg)
		return ColorizeWithBackground(text, fgColor, bgColor)

	case interfaces.Mode8bit:
		if f == -1 && b == -1 {
			return text
		}
		if b == -1 {
			return Colorize8(text, f)
		}
		return Colorize8WithBackground(text, f, b)

	default:
		return text
	}
}

func (l *VT100) SaveCursor() (int, error) {
	return l.z.Write(escSaveCursorDef)
}

func (l *VT100) RestoreCursor() (int, error) {
	return l.z.Write(escRestoreCursorDef)
}

func (l *VT100) MoveCursorLeft() (int, error) {
	return l.z.Write(escMoveCursorLeftDef)
}

func (l *VT100) MoveCursorRight() (int, error) {
	return l.z.Write(escMoveCursorRightDef)
}

func (l *VT100) MoveCursorTopLeft() (int, error) {
	return l.z.Write(escMoveCursorTopLeftDef)
}

func (l *VT100) ClearLine(_ string) (int, error) {
	//l.ResetBuffer()
	//l.current = []rune(line)
	//l.pos = len(l.current)

	// ClearLine sends the VT100 code for erasing the line followed by a carriage return
	// to move the cursor back to the beginning of the line
	clearLine := "\x1B[2K" + "\r"
	return l.z.Write([]byte(clearLine))
}

func (l *VT100) ClearScreen() (int, error) {
	clearScreen := "\x1B[2J" + "\r"
	return l.z.Write([]byte(clearScreen))
}

func (l *VT100) SetKeyFunc(e interfaces.KeyFunc) {
	l.keyFunc = e
}

func (l *VT100) SetSize(w int, h int) {
	if l.debug {
		log.Println("Screen size", w, h)
	}
}

func (l *VT100) SetEnterKey(key rune) {
	l.enterKey = key
}

func (l *VT100) colorFgAdapter(c interfaces.ColorDef) colorCode {
	var color colorCode

	switch c {
	case interfaces.ColorNoneDef:
		color = normal
	case interfaces.ColorRedDef:
		color = fgRed
	case interfaces.ColorGreenDef:
		color = fgGreen
	case interfaces.ColorYellowDef:
		color = fgYellow
	case interfaces.ColorBlueDef:
		color = fgBlue
	case interfaces.ColorMagentaDef:
		color = fgMagenta
	case interfaces.ColorCyanDef:
		color = fgCyan
	case interfaces.ColorWhiteDef:
		color = fgWhite
	case interfaces.ColorBrightRedDef:
		color = fgBrightRed
	case interfaces.ColorBrightGreenDef:
		color = fgBrightGreen
	case interfaces.ColorBrightYellowDef:
		color = fgBrightYellow
	case interfaces.ColorBrightBlueDef:
		color = fgBrightBlue
	case interfaces.ColorBrightMagentaDef:
		color = fgBrightMagenta
	case interfaces.ColorBrightCyanDef:
		color = fgBrightCyan
	case interfaces.ColorBrightWhiteDef:
		color = fgBrightWhite
	case interfaces.ColorBrightBlackDef:
		color = fgBrightBlack
	case interfaces.ColorBlackDef:
		color = fgBlack
	case interfaces.ColorGrayDef:
		color = fgBrightBlack
	case interfaces.ColorNormalDef:
		color = normal
	default:
		color = normal
	}

	return color
}

func (l *VT100) colorBgAdapter(c interfaces.ColorDef) colorCode {
	var color colorCode

	switch c {
	case interfaces.ColorNoneDef:
		color = normal
	case interfaces.ColorRedDef:
		color = bgRed
	case interfaces.ColorGreenDef:
		color = bgGreen
	case interfaces.ColorYellowDef:
		color = bgYellow
	case interfaces.ColorBlueDef:
		color = bgBlue
	case interfaces.ColorMagentaDef:
		color = bgMagenta
	case interfaces.ColorCyanDef:
		color = bgCyan
	case interfaces.ColorWhiteDef:
		color = bgWhite
	case interfaces.ColorBrightRedDef:
		color = bgBrightRed
	case interfaces.ColorBrightGreenDef:
		color = bgBrightGreen
	case interfaces.ColorBrightYellowDef:
		color = bgBrightYellow
	case interfaces.ColorBrightBlueDef:
		color = bgBrightBlue
	case interfaces.ColorBrightMagentaDef:
		color = bgBrightMagenta
	case interfaces.ColorBrightCyanDef:
		color = bgBrightCyan
	case interfaces.ColorBrightWhiteDef:
		color = bgBrightWhite
	case interfaces.ColorBrightBlackDef:
		color = bgBrightBlack
	case interfaces.ColorBlackDef:
		color = bgBlack
	case interfaces.ColorGrayDef:
		color = bgBrightBlack
	case interfaces.ColorNormalDef:
		color = normal
	default:
		color = normal
	}

	return color
}

//func (l * VT100) RetrieveBuffer() string {
//	out := string(l.current)
//	l.doReset()
//	return out
//}

/*
func (l * VT100) RetrieveTab() (string, bool) {
	var tab []rune
	found := false
	if l.pos == len(l.current) {
		tab = l.current[:l.pos]
		found = true
	}
	return string(tab), found
}
*/

func (l *VT100) Scan(data []byte) {
	escape := false
	var escapeSequence []byte
	var escapeParameter byte
	var escapeIntermediate byte

	if len(data) <= 0 {
		return
	}

	if l.debug {
		fmt.Printf("RECEIVED %s\r\n", string(data))
	}

	//UTF8
	sequence := []rune(string(data))

	for pos, key := range sequence {
		if escape {
			switch len(escapeSequence) {
			case 0:
				escape = false
			case 1:
				if key == 91 {
					escapeSequence = append(escapeSequence, byte(key))
				} else {
					escape = false
				}

			case 2, 3, 4:
				if key >= 0x30 && key <= 0x3F {
					escapeParameter = byte(key)
					escapeSequence = append(escapeSequence, byte(key))
				} else if key >= 0x20 && key <= 0x2F {
					escapeIntermediate = byte(key)
					escapeSequence = append(escapeSequence, byte(key))
				} else if key >= 0x40 && key <= 0x7E {
					escapeSequence = append(escapeSequence, byte(key))
					l.doEscape(escapeParameter, escapeIntermediate, byte(key))
					escape = false
				} else {
					escape = false
				}
			default:
				escape = false
			}
		} else {
			if l.debug {
				log.Printf("Key Pressed %d %d\r\n", key, pos)
			}

			if key == 9 {
				l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeTab, '\t'))
			} else {
				switch key {
				case l.enterKey:
					l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeEnter, '\n'))
				case 3:
					l.doCtrl(key)
				case 4:
					l.doCtrl(key)
				case 8:
					if l.keyFunc != nil {
						l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeBackspace, 8))
					}
				case 27:
					escape = true
					escapeSequence = nil
					escapeIntermediate = 0
					escapeParameter = 0
					escapeSequence = append(escapeSequence, byte(key))
				case 127:
					if l.keyFunc != nil {
						l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeBackspace, 8))
					}
				default:
					if l.keyFunc != nil {
						l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeKey, key))
					}
				}
			}
		}
	}
}

func (l *VT100) doEscape(parameter byte, intermediate byte, final byte) bool {
	if l.debug {
		log.Printf("Escape sequence %d %d %d\r\n", parameter, intermediate, final)
	}

	switch final {
	case 65:
		l.doMoveCursorUp()
	case 66:
		l.doMoveCursorDown()
	case 67:
		l.doMoveCursorRight()
	case 68:
		l.doMoveCursorLeft()
	case 126:
		if parameter == 51 {
			if l.keyFunc != nil {
				l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeCancel, 127))
			}
		}
	}

	return false
}

func (l *VT100) doCtrl(key rune) {
	if l.keyFunc != nil {
		l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeCtrl, key))
	}
}

func (l *VT100) doMoveCursorRight() {
	if l.keyFunc != nil {
		l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeCursor, rune(interfaces.CursorRightDef)))
	}
}

func (l *VT100) doMoveCursorLeft() {
	if l.keyFunc != nil {
		l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeCursor, rune(interfaces.CursorLeftDef)))
	}
}

func (l *VT100) doMoveCursorUp() {
	if l.keyFunc != nil {
		l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeCursor, rune(interfaces.CursorUpDef)))
	}
}

func (l *VT100) doMoveCursorDown() {
	if l.keyFunc != nil {
		l.keyFunc(interfaces.NewKeyData(interfaces.KeyTypeCursor, rune(interfaces.CursorDownDef)))
	}
}
