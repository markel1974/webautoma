package executor

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/markel1974/webautoma/src/wd/base"
)

// Window represents a browser window and contains metadata such as handle, name, state, and the associated web driver.
type Window struct {
	wd      base.IWebDriver
	varName string
	name    string
	handle  string
	open    bool
	timeout int
}

// NewWindow creates a new Window instance using the provided WebDriver and ConfigCommand settings.
// It initializes properties such as window handle, name, and timeout, and optionally delays based on the timeout duration.
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

// Store retrieves and stores the current window handle in the Window struct. Returns an error if the operation fails.
func (w *Window) Store() error {
	var err error
	if w.handle, err = w.wd.CurrentWindowHandle(); err != nil {
		return err
	}
	return nil
}

// GetVarName returns the variable name associated with the Window instance.
func (w *Window) GetVarName() string {
	return w.varName
}

// GetName returns the name of the window as a string.
func (w *Window) GetName() string {
	return w.name
}

// GetHandle retrieves the handle identifier of the current Window instance.
func (w *Window) GetHandle() string {
	return w.handle
}

// Windows manages a collection of browser windows and provides methods to interact with them.
type Windows struct {
	last      *Window
	container map[string]*Window
}

// NewWindows initializes and returns a new Windows instance with an empty container map and nil last window.
func NewWindows() *Windows {
	return &Windows{
		last:      nil,
		container: map[string]*Window{},
	}
}

// Add creates a new window instance using IWebDriver and ConfigCommand, updates the last window, and returns the timeout.
func (w *Windows) Add(wd base.IWebDriver, command ConfigCommand) int {
	win := NewWindow(wd, command)
	w.last = win
	return win.timeout
}

// Store saves the current `Window` instance into the `container` map using its variable name as the key and clears `last`.
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

// GetHandle retrieves the window handle associated with the provided target string and returns it or an error if unsupported.
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
