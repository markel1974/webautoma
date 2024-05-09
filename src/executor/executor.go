package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/markel1974/webautoma/src/wd/base"
)

func RequiredLogs() (base.LogType, base.LogLevel) {
	return base.LogPerformance, base.LogAll
}

type Executor struct {
	cfgFile   string
	start     time.Time
	quit      bool
	maxRetry  int
	timers    map[string]*Timer
	templates *Templates
	cfg       Config
	stack     map[string]interface{}
	adapter   *Adapter
	execId    string
	lastFrame interface{}
	baseUrl   string
}

func New(driver base.IWebDriver) *Executor {
	e := &Executor{
		templates: nil,
		quit:      false,
		maxRetry:  3,
		timers:    make(map[string]*Timer),
		stack:     nil,
		adapter:   NewAdapter(driver),
		execId:    "",
		lastFrame: nil,
		baseUrl:   "",
	}
	return e
}

func (e *Executor) Setup(sideFile string, logFile string, imgFile string, imgDump bool, capture []string, variables map[string]interface{}) error {
	var err error
	e.execId = computeFileId(sideFile)
	e.templates = NewTemplates(variables)
	if logFile, err = e.templates.Apply(logFile, nil); err != nil {
		return err
	}
	if imgFile, err = e.templates.Apply(imgFile, nil); err != nil {
		return err
	}
	if e.cfgFile, err = e.templates.Apply(sideFile, nil); err != nil {
		return err
	}
	if err = e.adapter.Setup(logFile, imgFile, imgDump, capture); err != nil {
		return err
	}
	var fileData []byte
	if fileData, err = os.ReadFile(e.cfgFile); err != nil {
		return err
	}
	if err = json.Unmarshal(fileData, &e.cfg); err != nil {
		return err
	}
	if e.cfg.Quit != nil {
		e.quit = *e.cfg.Quit
	}
	if e.cfg.Timeout != nil {
		e.adapter.SetWait(*e.cfg.Timeout)
	}
	if e.cfg.HumanWait != nil {
		e.adapter.SetHumanWait(*e.cfg.HumanWait)
	}
	if e.cfg.RetryInterval != nil {
		e.adapter.SetRetryInterval(*e.cfg.RetryInterval)
	}
	if e.cfg.MaxScreenshotLength != nil {
		e.adapter.SetMaxScreenshotLength(*e.cfg.MaxScreenshotLength)
	}
	return nil
}

func (e *Executor) SetImageSize(w int, h int) {
	e.adapter.SetImageSize(w, h)
}

func (e *Executor) loadUrl(url string) error {
	if err := e.adapter.Navigate(url); err != nil {
		return err
	}
	e.baseUrl = url
	return nil
}

func (e *Executor) doScreenshot(screenshotId string) error {
	return e.adapter.Screenshot(screenshotId)
}

func (e *Executor) doWriteLog(message string) {
	e.adapter.WriteLog(message)
}

func (e *Executor) logEvent(event *Event) {
	message, err := json.Marshal(event)
	if err != nil {
		e.adapter.log(LogLevelCritical, err.Error())
		return
	}
	e.doWriteLog(string(message))
}

func (e *Executor) computeSelector(target string) (string, string, error) {
	var container = strings.Split(target, "=")
	if len(container) <= 1 {
		return "", "", fmt.Errorf("invalid target")
	}
	var by string
	data := strings.TrimSpace(container[1])
	command := strings.TrimSpace(container[0])
	switch command {
	case "xpath":
		by = base.ByXPATH // selector.xpath(data)
	case "id":
		by = base.ByID //selector.id(data)
	case "linkText":
		by = base.ByLinkText // selector.linkText(data)
	case "name":
		by = base.ByName //selector.name(data)
	case "css":
		by = base.ByCSSSelector // selector.cssSelector(data)
	case "class":
		by = base.ByClassName // selector.cssSelector(data)
	default:
		return "", "", fmt.Errorf("unsupported by")
	}
	return by, data, nil
}

func (e *Executor) doTimer(target string, until string, value string) error {
	switch until {
	case "timerCreate":
		e.timers[target] = NewTimer(target, value)
		return nil
	case "timerStart":
		if t, ok := e.timers[target]; ok {
			t.Run()
			return nil
		}
		return fmt.Errorf("unknown timer id (timerStart): %s", target)
	case "timerStop":
		if t, ok := e.timers[target]; ok {
			t.Stop()
			return nil
		}
		return fmt.Errorf("unknown timer id (timerStop): %s", target)
	case "timerFinalize":
		if t, ok := e.timers[target]; ok {
			start := t.Start
			dur := t.Finalize()
			event := e.adapter.CreateEvent(t.Id, "full", nil, start, dur, false)
			event.Message = "Esito sonda ok"
			e.logEvent(event)
			return nil
		}
		return fmt.Errorf("unknown timer id (timerFinalize): %s", target)
	}
	return nil
}

func (e *Executor) doMouse(target string, command string, value string) error {
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	elm := e.adapter.FindElementReady(by, data, 2)
	if elm == nil {
		return fmt.Errorf("element isn't ready (%s)", target)
	}
	return e.adapter.MouseActions(elm, command, value)
}

func (e *Executor) doType(target string, value string) error {
	if len(value) <= 0 {
		return nil
	}
	if len(target) > 0 {
		by, data, err := e.computeSelector(target)
		if err != nil {
			return err
		}
		elm := e.adapter.FindElementReady(by, data, 2)
		if elm == nil {
			return fmt.Errorf("element isn't ready")
		}
		return e.adapter.SendKeys(elm, by, value)
	}
	for _, v := range []rune(value) {
		time.Sleep(time.Millisecond * 100)
		if err := e.adapter.KeyDown(string(v)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) doRunScript(script string, args []interface{}) error {
	return e.adapter.ExecuteScript(script, args)
}

func (e *Executor) doSelectFrame(target string) error {
	var frame interface{}
	if v := strings.Split(target, "="); len(v) >= 2 {
		switch strings.ToLower(strings.TrimSpace(v[0])) {
		case "index":
			var err error
			if frame, err = strconv.Atoi(v[1]); err != nil {
				return err
			}
			e.lastFrame = frame
		case "relative":
			frame = v[1]
			e.lastFrame = frame
		default:
			frame = v[1]
			e.lastFrame = frame
		}
	} else {
		frame = ""
		e.lastFrame = nil
	}
	return e.adapter.SwitchFrame(frame)
}

func (e *Executor) doSelectParentFrame() error {
	e.lastFrame = nil
	return e.adapter.SwitchParentFrame()
}

func (e *Executor) doSelectWindowMain() error {
	return e.adapter.SelectWindowMain()
}

func (e *Executor) doSelectWindowTitle(title string) error {
	currentWindow, err := e.adapter.CurrentWindowHandle()
	if err != nil {
		return err
	}
	lowerTitle := strings.ToLower(title)
	var found = false
	windows, err := e.adapter.WindowHandles()
	if err != nil {
		return err
	}
	for _, window := range windows {
		if err := e.adapter.SwitchWindow(window); err != nil {
			continue
		}
		windowTitle, err := e.adapter.Title()
		if err != nil {
			continue
		}
		if len(lowerTitle) == 0 && len(windowTitle) == 0 {
			found = true
			break
		}
		if strings.Contains(strings.ToLower(windowTitle), lowerTitle) {
			found = true
			break
		}
	}
	if !found {
		_ = e.adapter.SwitchWindow(currentWindow)
		return fmt.Errorf("can't find window with title " + title)
	}
	return nil
}

func (e *Executor) doSelectAlert(val string) error {
	value := parseInt(val)
	_, err := e.adapter.AlertText()
	if err != nil {
		return err
	}

	switch value {
	case 1:
		if err := e.adapter.AcceptAlert(); err != nil {
			return err
		}
	case 2:
		//alert.read()
	default:
		if err := e.adapter.DismissAlert(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) doCloseWindow() error {
	currentWindow, err := e.adapter.CurrentWindowHandle()
	if err != nil {
		return err
	}
	if currentWindow == e.adapter.GetRootWindow() {
		return fmt.Errorf("can't close mainWindow")
	}
	_ = e.adapter.Close()
	if v := e.adapter.SwitchWindow(e.adapter.GetRootWindow()); v == nil {
		return fmt.Errorf("invalid mainWindow")
	}
	return nil
}

func (e *Executor) doClose(target string) error {
	handle, err := e.adapter.GetWindowHandle(target)
	if err != nil {
		return err
	}
	if p := e.adapter.CloseWindow(handle); p == nil {
		return fmt.Errorf("invalid handle")
	}
	return nil
}

func (e *Executor) doStackAdd(id string, target string, until string) error {
	v := 0
	if len(until) > 0 {
		v = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	if analyzedElm := e.adapter.FindElementReady(by, data, v); analyzedElm != nil {
		values := make(map[string]interface{})
		values["target"] = target
		values["displayed"], _ = analyzedElm.IsDisplayed()
		values["enabled"], _ = analyzedElm.IsEnabled()
		values["text"], _ = analyzedElm.Text()
		values["tagName"], _ = analyzedElm.TagName()
		if e.stack == nil {
			e.stack = make(map[string]interface{})
		}
		e.stack[id] = values
	}
	return nil
}

func (e *Executor) doStackPrint() error {
	k, _ := json.MarshalIndent(e.stack, "", "	")
	fmt.Println(string(k))
	return nil
}

func (e *Executor) doStackReset() error {
	e.stack = nil
	return nil
}

func (e *Executor) doAssert(target string, until string, caption string) error {
	v := 0
	if len(until) > 0 {
		v = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	existsElm := e.adapter.FindElementReady(by, data, v)
	if existsElm == nil {
		return fmt.Errorf("element not exists")
	}
	text, err := existsElm.Text()
	if err != nil {
		return err
	}
	if text != caption {
		return fmt.Errorf("different text %s -> %s", text, caption)
	}
	return nil
}

func (e *Executor) doExists(target string, until string) error {
	v := 0
	if len(until) > 0 {
		v = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	if existsElm := e.adapter.FindElementReady(by, data, v); existsElm == nil {
		return fmt.Errorf("element not exists")
	}
	return nil
}

func (e *Executor) doUntil(target string, until string) error {
	v := 0
	if len(until) > 0 {
		v = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	if untilElm := e.adapter.FindElementReady(by, data, v); untilElm != nil {
		return fmt.Errorf("element exists")
	}
	return nil
}

func (e *Executor) doSelectWindow(target string) error {
	handle, err := e.adapter.GetWindowHandle(target)
	if err != nil {
		return err
	}
	if err := e.adapter.SwitchWindow(handle); err != nil {
		return err
	}
	return nil
}

func (e *Executor) doStoreWindowHandle(target string) error {
	return e.adapter.StoreWindowHandle(target)
}

func (e *Executor) doSelect(target string, value string) error {
	v := strings.Split(value, "=")
	if len(v) < 2 {
		return fmt.Errorf("invalid value")
	}
	if byValue := v[0]; byValue != "label" {
		return fmt.Errorf("unsupported by")
	}
	label := v[1]
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	elm := e.adapter.FindElementReady(by, data, 2)
	if elm == nil {
		return fmt.Errorf("element isn't ready")
	}
	options, err := elm.FindElements(base.ByTagName, "option")
	if err != nil {
		return err
	}
	for _, opt := range options {
		//name, _ := k.Text()
		//fmt.Println(name)
		var computedLabel string
		computedLabel, err = opt.ComputedLabel()
		if err != nil {
			continue
		}
		if computedLabel == label {
			return opt.Click()
		}
	}
	return fmt.Errorf("not found")
}

func (e *Executor) doActiveElement(id string) error {
	elm, err := e.adapter.ActiveElement()
	if err != nil {
		return err
	}
	fmt.Printf("%s ActiveElement =>\n", id)
	elm.Print()
	return nil
}

func (e *Executor) doWindowHandles(id string) error {
	fmt.Printf("%s WindowHandles =>\n", id)
	v, _ := e.adapter.WindowHandles()
	fmt.Println(v)
	return nil
}

func (e *Executor) doPause(target string) error {
	v, err := strconv.Atoi(target)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * time.Duration(v))
	return nil
}

func (e *Executor) doSetTimeout(value string) error {
	e.adapter.SetWait(parseInt(value))
	return nil
}

func (e *Executor) doPageSource(id string) error {
	v, _ := e.adapter.PageSource()
	fmt.Printf("%s PageSource =>\n", id)
	fmt.Println(v)
	return nil
}

func (e *Executor) doStatus(id string) error {
	fmt.Printf("%s Status =>\n", id)
	k, _ := e.adapter.Status()
	fmt.Println(k)
	return nil
}

func (e *Executor) doGetCookie(target string, id string) error {
	v, _ := e.adapter.GetCookie(target)
	fmt.Printf("%s Cookie =>\n", id)
	fmt.Println(v)
	return nil
}

func (e *Executor) doGetAllCookies(id string) error {
	v, _ := e.adapter.GetCookies()
	fmt.Printf("%s Cookies =>\n", id)
	fmt.Println(v)
	return nil
}

func (e *Executor) doDeleteAllCookies() error {
	return e.adapter.DeleteAllCookies()
}

func (e *Executor) doDeleteCookie(target string) error {
	return e.adapter.DeleteCookie(target)
}

func (e *Executor) doActionsSendKeys(value string) error {
	return e.adapter.KeyDown(value)
}

func (e *Executor) doSetWindowMain() error {
	return e.adapter.SetRootWindow()
}

func (e *Executor) doHumanWait(value string) error {
	e.adapter.Sleep(parseInt(value))
	return nil
}

func (e *Executor) doEnableDebug() error {
	e.adapter.EnableDebug(true)
	return nil
}

func (e *Executor) doDisableDebug() error {
	e.adapter.EnableDebug(false)
	return nil
}

func (e *Executor) doDisableError() error {
	e.adapter.SetErrorDisabled(true)
	return nil
}

func (e *Executor) doEnableError() error {
	e.adapter.SetErrorDisabled(false)
	return nil
}

func (e *Executor) doExecId(target string) error {
	e.execId = target
	return nil
}

func (e *Executor) doProbeId(target string) error {
	e.adapter.SetProbeId(target)
	return nil
}

func (e *Executor) doSetWindowSize(target string) error {
	v := strings.Split(target, "x")
	if len(v) < 2 {
		return fmt.Errorf("invalid target")
	}
	width, _ := strconv.Atoi(v[0])
	height, _ := strconv.Atoi(v[1])
	return e.adapter.ResizeWindow("", width, height)
}

func (e *Executor) doScrollTo(target string, value string) error {
	v := strings.Split(value, "x")
	if len(v) < 2 {
		return fmt.Errorf("invalid value")
	}
	x, _ := strconv.Atoi(v[0])
	y, _ := strconv.Atoi(v[1])
	var elm base.IWebElement
	if len(target) > 0 {
		by, data, err := e.computeSelector(target)
		if err != nil {
			return err
		}
		if elm = e.adapter.FindElementReady(by, data, 2); elm == nil {
			return fmt.Errorf("element isn't ready")
		}
	} else {
		var err error
		if elm, err = e.adapter.ActiveElement(); err != nil {
			return err
		}
	}
	if err := elm.MoveTo(0, 0); err != nil {
		return err
	}
	return elm.ScrollTo(x, y)
}

func (e *Executor) doNavigate(target string) error {
	if strings.HasPrefix(target, "/") {
		target = e.baseUrl + target
	}
	return e.adapter.Navigate(target)
}

func (e *Executor) commandExec(id string, command string, target string, until string, value string) error {
	var err error
	start := time.Now()

	e.adapter.HumanWait()

	switch command {
	case "execId":
		err = e.doExecId(target)
	case "probe":
		err = e.doProbeId(target)
	case "until":
		err = e.doUntil(target, until)
	case "open":
		err = e.doNavigate(target)
	case "setWindowSize":
		err = e.doSetWindowSize(target)
	case "timerCreate", "timerStart", "timerStop", "timerFinalize":
		err = e.doTimer(target, command, value)
	case "assert":
		err = e.doAssert(target, until, value)
	case "exists":
		err = e.doExists(target, until)
	case "stackAdd":
		err = e.doStackAdd(id, target, until)
	case "stackReset":
		err = e.doStackReset()
	case "stackPrint":
		err = e.doStackPrint()
	case "disableError":
		err = e.doDisableError()
	case "enableError":
		err = e.doEnableError()
	case "disableDebug":
		err = e.doDisableDebug()
	case "enableDebug":
		err = e.doEnableDebug()
	case "click", "rightClick", "doubleClick", "mouseUpAt", "mouseDownAt", "mouseMultipleMoveAt", "mouseMoveAt", "mouseOver":
		err = e.doMouse(target, command, value)
	case "mouseOut":
		//nothing to do
	case "type":
		err = e.doType(target, value)
	case "runScript":
		err = e.doRunScript(target, nil)
	case "humanWait":
		err = e.doHumanWait(value)
	case "select":
		err = e.doSelect(target, value)
	case "setWindowMain":
		err = e.doSetWindowMain()
	case "selectWindow":
		err = e.doSelectWindow(target)
	case "selectWindowMain":
		err = e.doSelectWindowMain()
	case "selectWindowTitle":
		err = e.doSelectWindowTitle(value)
	case "selectFrame":
		err = e.doSelectFrame(target)
	case "selectParentFrame":
		err = e.doSelectParentFrame()
	case "selectAlert":
		err = e.doSelectAlert(value)
	case "closeWindow":
		err = e.doCloseWindow()
	case "actionsSendKeys":
		err = e.doActionsSendKeys(value)
	case "setTimeout":
		err = e.doSetTimeout(value)
	case "activeElement":
		err = e.doActiveElement(id)
	case "pageSource":
		err = e.doPageSource(id)
	case "getAllCookies":
		err = e.doGetAllCookies(id)
	case "getCookie":
		err = e.doGetCookie(target, id)
	case "deleteAllCookies":
		err = e.doDeleteAllCookies()
	case "deleteCookie":
		err = e.doDeleteCookie(target)
	case "windowHandles":
		err = e.doWindowHandles(id)
	case "storeWindowHandle":
		err = e.doStoreWindowHandle(target)
	case "close":
		err = e.doClose(target)
	case "status":
		err = e.doStatus(id)
	case "pause":
		err = e.doPause(target)
	case "scrollTo":
		err = e.doScrollTo(target, value)
	case "noop":
		//nothing to do
	default:
		e.adapter.log(LogLevelWarning, "unimplemented command: "+command)
	}
	event := e.adapter.CreateEvent(e.execId, "intermediate", err, start, UnixMilli(time.Now())-UnixMilli(start), true)
	event.Label = id
	event.Target = target
	event.Command = command
	e.logEvent(event)
	return err
}

func (e *Executor) finalize(err error) {
	for _, t := range e.timers {
		if !t.Finalized {
			start := t.Start
			dur := t.Finalize()
			event := e.adapter.CreateEvent(t.Id, "full", err, start, dur, false)
			event.Message = "timer not finalized"
			e.logEvent(event)
		}
	}
	var dur = UnixMilli(time.Now()) - UnixMilli(e.start)
	var event = e.adapter.CreateEvent(e.execId, "full", err, e.start, dur, false)
	if err != nil {
		event.Message = "Errore nella sonda"
	} else {
		event.Message = "Esito sonda ok"
	}
	e.logEvent(event)
	if e.quit {
		_ = e.adapter.Quit()
	}
}

func (e *Executor) doComputeJump(target string, commands []ConfigCommand) (int, error) {
	if len(target) == 0 {
		return -1, fmt.Errorf("empty target")
	}
	for idx, cmd := range commands {
		if cmd.Id == target {
			return idx, nil
		}
	}
	return -1, fmt.Errorf("undefined id: " + target)
}

func (e *Executor) commandsLoop() (string, error) {
	for _, test := range e.cfg.Tests {
		for x := 0; x < len(test.Commands); x++ {
			cmd := test.Commands[x]
			err := e.templates.BuildCommand(cmd, e.stack)
			if err != nil {
				return cmd.Id, err
			}
			if cmd.Command == "jump" {
				var jump int
				if jump, err = e.doComputeJump(cmd.Target, test.Commands); err == nil {
					x = jump
				}
			} else {
				windowTimeout := 0
				if err == nil && len(cmd.WindowHandleName) > 0 {
					windowTimeout = e.adapter.AddWindowHandle(cmd)
				}
				if e.lastFrame != nil {
					/*
						//if frameErr := e.adapter.SwitchParentFrame(); frameErr != nil {
						//	e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchParentFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
						//}
						if frameErr := e.adapter.SwitchFrame(""); frameErr != nil {
							e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
						}
						if frameErr := e.adapter.SwitchFrame(e.lastFrame); frameErr != nil {
							e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
						}
					*/
				}
				err = e.commandExec(cmd.Id, cmd.Command, cmd.Target, cmd.Until, cmd.Value)
				if windowTimeout > 0 {
					e.adapter.Sleep(windowTimeout)
				}
			}
			if e.adapter.IsErrorDisabled() {
				err = nil
			}
			if err != nil {
				return cmd.Id, err
			}
		}
	}
	return "", nil
}

func (e *Executor) Start() error {
	e.start = time.Now()

	e.adapter.log(LogLevelInfo, "Sample started")

	err := e.adapter.SetRootWindow()
	if err != nil {
		e.finalize(err)
		return err
	}

	url, err := e.templates.Apply(e.cfg.Url, nil)
	if err != nil {
		e.finalize(err)
		return err
	}
	if err = e.loadUrl(url); err != nil {
		e.finalize(err)
		return err
	}

	e.adapter.log(LogLevelInfo, "mainWindow: "+e.adapter.GetRootWindow())

	var id string
	if id, err = e.commandsLoop(); err != nil {
		e.finalize(fmt.Errorf("%s (%s)", err.Error(), id))
		return err
	}

	e.finalize(nil)
	return nil
}
