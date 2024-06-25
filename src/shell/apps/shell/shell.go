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

package shell

import (
	"fmt"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"log"
	"unicode"
)

const (
	stateUndefined        = iota
	stateUsernameRequired = iota
	statePasswordRequired = iota
	stateAuthenticated    = iota
)

const (
	usernamePrompt   = "Username: "
	passwordPrompt   = "Password: "
	maxPasswordRetry = 3
)

type ExecSuggestionType func(in string, count int) bool

type ExecCommandType func(command string) bool

type Shell struct {
	current  []rune
	pos      int
	echo     bool
	history  *HistoryHandler
	tabData  string
	tabFound bool
	tabCount int
	terminal interfaces.ITerminal

	defaultPrompt   string
	prompt          string
	currentUsername string
	passwordRetry   int
	state           int
	auth            interfaces.IAuthenticator
	ExecSuggestion  ExecSuggestionType
	ExecCommand     ExecCommandType
}

func NewShell(auth interfaces.IAuthenticator, terminal interfaces.ITerminal, prompt string, autosave bool) *Shell {
	c := &Shell{
		history:       NewHistoryHandler(128, autosave),
		echo:          true,
		terminal:      terminal,
		auth:          auth,
		defaultPrompt: prompt,
		passwordRetry: 0,
		state:         stateUndefined,
	}

	switch c.auth.GetAuthenticationMode() {
	case interfaces.AuthModeFull:
		c.setUsernameRequiredState()
	case interfaces.AuthModePassword:
		c.setPasswordRequiredState()
	case interfaces.AuthModeNone:
		c.setAuthenticatedState()
	default:
		c.setAuthenticatedState()
	}

	return c
}

func (c *Shell) KeyEvent(event *interfaces.KeyData) bool {
	ret := false
	switch event.GetType() {
	case interfaces.KeyTypeEnter:
		ret = c.enterPressed()
	case interfaces.KeyTypeTab:
		c.tabPressed()
	case interfaces.KeyTypeCancel:
		c.textCancel()
	case interfaces.KeyTypeBackspace:
		c.textBackspace()
	case interfaces.KeyTypeKey:
		c.keyPressed(event.Key)
	case interfaces.KeyTypeCursor:
		c.cursorPressed(interfaces.CursorCodeDef(event.Key))
	}

	return ret
}

func (c *Shell) ClearHistory() {
	c.history.Clear()
}

func (c *Shell) GetHistoryAtPos(idx int) (string, bool) {
	return c.history.GetHistoryAtPos(idx)
}

func (c *Shell) GetHistory() string {
	out := ""
	for n, x := range c.history.GetHistory() {
		out += "\r\n"
		out += fmt.Sprintf("%d: %s", n, x)
	}
	return out
}

func (c *Shell) SetHistoryDefault(data string) {
	c.history.SetDefault(data)
}

func (c *Shell) DoNext() {
	c.resetBuffer()
	_, _ = c.terminal.WriteColor("\r\n", interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	_, _ = c.terminal.WriteColor(c.prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

func (c *Shell) DoRedraw(line string) {
	c.current = []rune(line)
	c.pos = len(c.current)
	_, _ = c.terminal.ClearLine(line)
	_, _ = c.terminal.WriteColor(c.prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	_, _ = c.terminal.WriteColor(line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

func (c *Shell) cursorPressed(code interfaces.CursorCodeDef) {
	switch code {
	case interfaces.CursorUpDef:
		if data, valid := c.history.GetHistoryPrev(); valid {
			c.DoRedraw(data)
		}
	case interfaces.CursorDownDef:
		if data, valid := c.history.GetHistoryNext(); valid {
			c.DoRedraw(data)
		}

	case interfaces.CursorLeftDef:
		if c.pos > 0 {
			c.pos--
			_, _ = c.terminal.MoveCursorLeft()
		}
	case interfaces.CursorRightDef:
		if c.pos >= 0 && c.pos < len(c.current) {
			c.pos++
			_, _ = c.terminal.MoveCursorRight()
		}
	}
}

func (c *Shell) enterPressed() bool {
	buffer := string(c.current)
	quit := false

	if len(buffer) > 0 {
		switch c.state {
		case stateUsernameRequired:
			c.passwordRetry = 0
			c.currentUsername = buffer
			c.setPasswordRequiredState()

		case statePasswordRequired:
			if c.auth.IsAuthenticated(c.currentUsername, buffer) {
				c.setAuthenticatedState()
			} else {
				c.passwordRetry++
				if c.passwordRetry >= maxPasswordRetry {
					_, _ = c.terminal.WriteColor("\r\nUnauthorized\r\n", interfaces.ColorRedDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
					quit = true
				} else {
					_, _ = c.terminal.WriteColor("\r\nLogin incorrect\r\n", interfaces.ColorRedDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
				}
			}

		case stateAuthenticated:
			c.history.AddToHistory(buffer)
			c.history.SetDefault("")

			if c.ExecCommand != nil {
				c.ExecCommand(buffer)
			}
			//c.quit = c.execCommand(buffer)

		default:
			quit = true
		}
	}

	c.DoNext()

	return quit
}

func (c *Shell) tabPressed() {
	if c.state == stateAuthenticated {
		if c.tabCount == 0 {
			c.tabFound = false
			c.tabData = ""

			if c.pos == len(c.current) {
				c.tabData = string(c.current[:c.pos])
				c.tabFound = true
			}
		}
		if c.tabFound {

			if c.ExecSuggestion != nil {
				c.ExecSuggestion(c.tabData, c.tabCount)
			}
			//c.terminalExecSuggestion(c.tabData, c.tabCount)
		}
		c.tabCount++
	}
}

func (c *Shell) keyPressed(key rune) {
	if unicode.IsPrint(key) {
		if c.pos < 0 {
			log.Println("doTextInsert: negative pos", c.pos)
		} else if c.pos == len(c.current) {
			if c.echo {
				_, _ = c.terminal.WriteColor(string(key), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			}
			c.current = append(c.current, key)
			c.pos++
		} else if c.pos < len(c.current) {
			if c.echo {
				_, _ = c.terminal.WriteColor(string(key), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
				_, _ = c.terminal.SaveCursor()
				_, _ = c.terminal.WriteColor(string(c.current[c.pos:]), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			}
			_, _ = c.terminal.RestoreCursor()

			c.current = insertAtPos(c.current, key, c.pos)
			c.pos++
		} else {
			log.Println("terminalKeyPressed: invalid pos", c.pos)
		}
	}

	if c.state == stateAuthenticated {
		c.history.SetDefault(string(c.current))
		c.tabCount = 0
	}
}

func (c *Shell) setUsernameRequiredState() {
	c.echo = true
	c.prompt = usernamePrompt
	c.history.SetEnabled(false)
	c.state = stateUsernameRequired
}

func (c *Shell) setPasswordRequiredState() {
	c.echo = false
	c.prompt = passwordPrompt
	c.history.SetEnabled(false)
	c.state = statePasswordRequired
}

func (c *Shell) setAuthenticatedState() {
	c.echo = true
	c.prompt = c.defaultPrompt
	c.history.SetEnabled(true)
	c.state = stateAuthenticated
}

func (c *Shell) textBackspace() {
	if c.pos > 0 {
		c.pos--
		c.current = removeAtPos(c.current, c.pos)

		if c.echo {
			_, _ = c.terminal.MoveCursorLeft()
			_, _ = c.terminal.SaveCursor()
			_, _ = c.terminal.WriteColor(string(c.current[c.pos:]), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			_, _ = c.terminal.WriteColor(string(' '), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			_, _ = c.terminal.RestoreCursor()
		}
	}
}

func (c *Shell) textCancel() {
	if c.pos >= 0 {
		c.current = removeAtPos(c.current, c.pos)

		if c.echo {
			_, _ = c.terminal.SaveCursor()
			_, _ = c.terminal.WriteColor(string(c.current[c.pos:]), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			_, _ = c.terminal.WriteColor(string(' '), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			_, _ = c.terminal.RestoreCursor()
		}
	}
}

func (c *Shell) resetBuffer() {
	c.current = nil
	c.pos = 0
}
