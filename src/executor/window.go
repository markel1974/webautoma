package executor

import (
	"errors"
	"fmt"
	"github.com/markel1974/webautoma/src/wd/base"
	"strings"
	"time"
)

type Window struct {
	wd      base.IWebDriver
	varName string
	name    string
	handle  string
	open    bool
	timeout int
}

func NewWindow(wd base.IWebDriver, command ConfigCommand) *Window {
	w := &Window{
		wd:      wd,
		varName: "${" + command.WindowHandleName + "}",
		name:    command.WindowHandleName,
		handle:  "",
		open:    command.OpensWindow,
		timeout: command.WindowTimeout,
	}
	if w.timeout > 0 {
		time.Sleep(time.Millisecond * time.Duration(w.timeout))
	}
	return w
}

func (w *Window) Store() error {
	var err error
	if w.handle, err = w.wd.CurrentWindowHandle(); err != nil {
		return err
	}
	return nil
}

func (w *Window) GetVarName() string {
	return w.varName
}

func (w *Window) GetName() string {
	return w.name
}

func (w *Window) GetHandle() string {
	return w.handle
}

type Windows struct {
	last      *Window
	container map[string]*Window
}

func NewWindows() *Windows {
	return &Windows{
		last:      nil,
		container: map[string]*Window{},
	}
}

func (w *Windows) Add(wd base.IWebDriver, command ConfigCommand) int {
	win := NewWindow(wd, command)
	w.last = win
	return win.timeout
}

func (w *Windows) Store(target string) error {
	//TODO
	//"target": "root",
	if w.last == nil {
		return fmt.Errorf("nil stack")
	}
	if err := w.last.Store(); err != nil {
		return err
	}
	w.container[w.last.GetVarName()] = w.last
	w.last = nil
	return nil
}

func (w *Windows) GetHandle(target string) (string, error) {
	//"handle=${win9266}",
	v := strings.Split(target, "=")
	action := v[0]
	variable := v[1]
	if action != "handle" {
		return "", errors.New("unsupported action")
	}
	handle := variable
	if nw, ok := w.container[variable]; ok {
		handle = nw.GetHandle()
	}
	return handle, nil
}
