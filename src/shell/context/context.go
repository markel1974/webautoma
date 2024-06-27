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

package context

import (
	"github.com/markel1974/webautoma/src/shell/adaptiveticker"
	"github.com/markel1974/webautoma/src/shell/apps"
	"github.com/markel1974/webautoma/src/shell/apps/shell"
	"github.com/markel1974/webautoma/src/shell/cli"
	"github.com/markel1974/webautoma/src/shell/interfaces"
	"github.com/markel1974/webautoma/src/shell/terminal"
	"io"
)

const (
	contextMaQueueLen = 1024
)

type Context struct {
	Exit        bool
	ticker      *adaptiveticker.AdaptiveTicker
	reader      io.Reader
	writer      io.Writer
	factory     *terminal.EquipmentFactory
	template    *cli.Command
	terminal    interfaces.ITerminal
	auth        interfaces.IAuthenticator
	defaultApp  *shell.Shell
	enterKey    rune
	tasks       *TaskManager
	messageChan chan iMessage
	timersChan  chan *adaptiveticker.TimerHandler
	prompt      string
	autosave    bool
}

func NewContext(ticker *adaptiveticker.AdaptiveTicker, reader io.Reader, writer io.Writer, auth interfaces.IAuthenticator, factory *terminal.EquipmentFactory, template *cli.Command, prompt string, autosave bool) *Context {
	ctx := &Context{
		ticker:      ticker,
		reader:      reader,
		writer:      writer,
		auth:        auth,
		factory:     factory,
		template:    template,
		Exit:        false,
		prompt:      prompt,
		enterKey:    -1,
		messageChan: make(chan iMessage, contextMaQueueLen),
		timersChan:  make(chan *adaptiveticker.TimerHandler, contextMaQueueLen),
		tasks:       nil,
		autosave:    autosave,
	}
	return ctx
}

func (c *Context) Setup(terminalType string, debug bool) {
	if len(terminalType) == 0 {
		terminalType = "VT100"
	}
	c.terminal = c.factory.Create(terminalType, c.writer, debug)
	c.terminal.SetKeyFunc(c.keyHandler)
	if c.enterKey > -1 {
		c.terminal.SetEnterKey(c.enterKey)
	}

	template := apps.NewTemplate(c, c.writer)
	root := template.Run(c.template)

	c.tasks = NewTaskManager(c.ticker, c.timersChan, root)

	c.defaultApp = shell.NewShell(c.auth, c.terminal, c.prompt, c.autosave)
	c.defaultApp.ExecCommand = c.execCommand
	c.defaultApp.ExecSuggestion = c.execSuggestion
}

func (c *Context) SetScreenSize(width int, height int) {
	c.terminal.SetSize(width, height)
	c.tasks.SetScreenSize(width, height)
}

func (c *Context) keyHandler(event *interfaces.KeyData) {
	if event.GetType() == interfaces.KeyTypeCtrl {
		c.ctrlPressed(event.Key)
		return
	}

	if fgPid := c.tasks.GetForegroundPid(); fgPid != adaptiveticker.UnknownId {
		c.tasks.ExecRead(fgPid, int(event.GetType()), event.Key)
		return
	}

	quit := c.defaultApp.KeyEvent(event)
	if quit {
		c.Exit = true
	}
}

func (c *Context) SetEnterKey(key rune) {
	c.enterKey = key
}

func (c *Context) Close() {
}

func (c *Context) Exec(async bool) {
	go func() {
		readBuffer := make([]byte, 16)
		for {
			if n, err := c.reader.Read(readBuffer); err == nil {
				if n > 0 {
					var readEvent = newMessageRead(readBuffer, n)
					readEvent.postEvent(c.messageChan)
					//d := string(readBuffer)
					//fmt.Println(d)
				}
			} else {
				var quitEvent = newMessageQuit()
				quitEvent.postEvent(c.messageChan)
				return
			}
			/*
				if n, err := c.reader.Read(readBuffer); err == nil {
					if n > 0 {
						var readEvent = newMessageRead(readBuffer, n)
						readEvent.postEvent(c.messageChan)
					}
				} else {
					var quitEvent = newMessageQuit()
					quitEvent.postEvent(c.messageChan)
					return
				}
			*/
		}
	}()

	if async {
		go func() {
			c.eventLoop()
		}()
	} else {
		c.eventLoop()
	}
}

func (c *Context) execCommand(line string) bool {
	return c.tasks.Execute(line, nil)
}

func (c *Context) ctrlPressed(key rune) {
	switch key {
	case 3:
		c.tasks.SetSelectionDisabled()
		c.tasks.KillForeground()
		c.defaultApp.DoNext()
	case 4:
		c.tasks.ExecActivate()
	}
}

func (c *Context) execSuggestion(in string, count int) bool {
	ret := false
	data, suggestions, found := c.tasks.GetSuggestion(in, count)
	if found {
		sLen := len(suggestions)
		if sLen > 0 {
			idx := count % sLen
			if idx < sLen {
				complete := suggestions[idx]
				if len(complete) > len(data) {
					tabLine := in + complete[len(data):]
					c.defaultApp.DoRedraw(tabLine)
					c.defaultApp.SetHistoryDefault(tabLine)
					ret = true
				}
			}
		}
	}
	return ret
}

func (c *Context) eventLoop() {
	_, _ = c.terminal.WriteColor("Admin Console Ready", interfaces.ColorBlueDef, interfaces.ColorRedDef, interfaces.ModeNormal)

	c.defaultApp.DoNext()

	for {
		select {
		case m := <-c.messageChan:
			c.messageEventHandler(m)
		case t := <-c.timersChan:
			c.messageEventHandler(t.Event.(iMessage))
		}

		if c.Exit {
			c.shutdown()
			return
		}
	}
}

func (c *Context) messageEventHandler(m iMessage) {
	if m != nil {

		switch m.getType() {
		case MessageTypeRead:
			if mm, ok := m.(*MessageRead); ok {
				c.terminal.Scan(mm.data)
			}

		case MessageTypeTimer:
			if mt, ok := m.(*MessageTimer); ok {
				c.tasks.ExecTimer(mt.pid, mt.tid, mt.interval)
			}

		case MessageTypePaint:
			if _, ok := m.(*MessagePaint); ok {
				c.tasks.ExecPaint(c.terminal)
			}

		case MessageTypeQuit:
			if _, ok := m.(*MessageQuit); ok {
				c.Exit = true
			}
		}
	}
}

func (c *Context) shutdown() {
	c.tasks.KillAll("")
}

//CLI INTERFACE

func (c *Context) SetFg(pid int) bool {
	return c.tasks.SetFg(pid)
}

func (c *Context) TaskList() string {
	return c.tasks.List()
}

func (c *Context) Write(data string) {
	_, _ = c.terminal.Write(data)
}
func (c *Context) WriteLn(data string) {
	_, _ = c.terminal.Write(data + "\r\n")
}

func (c *Context) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	_, _ = c.terminal.WriteColor(data, fg, bg, mode)
}

func (c *Context) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	_, _ = c.terminal.WriteColor(data, fg, bg, mode)
	_, _ = c.terminal.Write("\r\n")
}

func (c *Context) ClearScreen() {
	_, _ = c.terminal.ClearScreen()
}

func (c *Context) SaveTasks(name string) bool {
	return c.tasks.SaveTasks(name)
}

func (c *Context) RestoreTasks(name string) bool {
	return c.tasks.RestoreTasks(name)
}

func (c *Context) ListTasks() []string {
	return c.tasks.ListTasks()
}

func (c *Context) SetContext(pid int, ctx interface{}) bool {
	return c.tasks.SetContext(pid, ctx)
}

func (c *Context) SetCaption(pid int, caption string) bool {
	return c.tasks.SetCaption(pid, caption)
}

func (c *Context) PaintRequest(_ int) bool {
	return c.tasks.PaintRequest()
}

func (c *Context) CreateTimer(pid int, first int, interval int, count int) bool {
	return c.tasks.CreateTimer(pid, first, interval, count)
}

func (c *Context) StopTimer(pid int, tid int) bool {
	return c.tasks.StopTimer(pid, tid)
}

func (c *Context) IsActive(pid int) bool {
	return c.tasks.IsActive(pid)
}

func (c *Context) Deactivate(pid int) bool {
	return c.tasks.Kill(pid)
}

func (c *Context) DeactivateAll(name string) int {
	return c.tasks.KillAll(name)
}

func (c *Context) GetScreenSize() (int, int) {
	return c.tasks.GetScreenSize()
}

func (c *Context) SetExit() {
	c.Exit = true
}

func (c *Context) SetBasePath(arg string) {
	c.tasks.SetBasePath(arg)
}

func (c *Context) SetSelectionMode(pid int) {
	c.tasks.SetSelectionMode(pid)
}

func (c *Context) SetSelectionModeNext() {
	c.tasks.SetSelectionModePrevious()
}

func (c *Context) SetSelectionModePrevious() {
	c.tasks.SetSelectionModePrevious()
}

func (c *Context) SetSelectionOptions(option rune, value float64) bool {
	return c.tasks.SetSelectionOptions(option, value)
}

func (c *Context) History(verb interfaces.HistoryAction, idx int) {
	switch verb {
	case interfaces.HistoryActionClear:
		c.defaultApp.ClearHistory()
	case interfaces.HistoryActionExec:
		if arg, found := c.defaultApp.GetHistoryAtPos(idx); found {
			c.execCommand(arg)
		}
	case interfaces.HistoryActionList:
		_, _ = c.terminal.Write(c.defaultApp.GetHistory())
	}
}
