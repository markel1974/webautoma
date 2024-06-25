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

import "fmt"

type colorCode string

const (
	/*
		red     colorCode = "\033[01;31m"
		green   colorCode = "\033[01;32m"
		yellow  colorCode = "\033[01;33m"
		blue    colorCode = "\033[01;34m"
		magenta colorCode = "\033[01;35m"
		cyan    colorCode = "\033[01;36m"
		white   colorCode = "\033[01;37m"

		brightRed     colorCode = "\033[22;31m"
		brightGreen   colorCode = "\033[22;32m"
		brightYellow  colorCode = "\033[22;33m"
		brightBlue    colorCode = "\033[22;34m"
		brightMagenta colorCode = "\033[22;35m"
		brightCyan    colorCode = "\033[22;36m"

		black  colorCode = "\033[22;30m"
		gray   colorCode = "\033[22;37m"
	*/
	normal colorCode = "\033[0m"
)

const (
	fgBlack   colorCode = "\033[30m"
	fgRed     colorCode = "\033[31m"
	fgGreen   colorCode = "\033[32m"
	fgYellow  colorCode = "\033[33m"
	fgBlue    colorCode = "\033[34m"
	fgMagenta colorCode = "\033[35m"
	fgCyan    colorCode = "\033[36m"
	fgWhite   colorCode = "\033[37m"

	fgBrightBlack   colorCode = "\033[30;1m"
	fgBrightRed     colorCode = "\033[31;1m"
	fgBrightGreen   colorCode = "\033[32;1m"
	fgBrightYellow  colorCode = "\033[33;1m"
	fgBrightBlue    colorCode = "\033[34;1m"
	fgBrightMagenta colorCode = "\033[35;1m"
	fgBrightCyan    colorCode = "\033[36;1m"
	fgBrightWhite   colorCode = "\033[37;1"
)

const (
	bgBlack   colorCode = "\033[40m"
	bgRed     colorCode = "\033[41m"
	bgGreen   colorCode = "\033[42m"
	bgYellow  colorCode = "\033[43m"
	bgBlue    colorCode = "\033[44m"
	bgMagenta colorCode = "\033[45m"
	bgCyan    colorCode = "\033[46m"
	bgWhite   colorCode = "\033[47m"

	bgBrightBlack   colorCode = "\033[40;1m"
	bgBrightRed     colorCode = "\033[41;1m"
	bgBrightGreen   colorCode = "\033[42;1m"
	bgBrightYellow  colorCode = "\033[43;1m"
	bgBrightBlue    colorCode = "\033[44;1m"
	bgBrightMagenta colorCode = "\033[45;1m"
	bgBrightCyan    colorCode = "\033[46;1m"
	bgBrightWhite   colorCode = "\033[47;1"
)

const (
	fg8bit colorCode = "\033[38:5:"
)

const (
	bg8bit colorCode = "\033[48:5:"
)

func Colorize(text string, fg colorCode) string {
	terminator := normal
	return fmt.Sprintf("%s%s%s", string(fg), text, string(terminator))
}

func ColorizeWithBackground(text string, fg colorCode, bg colorCode) string {
	terminator := normal
	return fmt.Sprintf("%s%s%s%s", string(fg), string(bg), text, string(terminator))
}

func Colorize8(text string, fg int) string {
	terminator := normal
	return fmt.Sprintf("%s%d%s%s", fg8bit, fg, text, string(terminator))
}

func Colorize8WithBackground(text string, fg int, bg int) string {
	terminator := normal
	out := fmt.Sprintf("%s%d%s%d%s%s", fg8bit, fg, bg8bit, bg, text, string(terminator))
	fmt.Println("TEST", out)
	return out
}
