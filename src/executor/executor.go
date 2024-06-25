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

const (
	eventMessageOK                = "Esito sonda ok"
	eventMessageNOK               = "Errore nella sonda"
	eventMessageTimerNotFinalized = "timer not finalized"
)

const (
	eventKindFull         = "full"
	eventKindIntermediate = "intermediate"
)

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

func (e *Executor) Setup(sideFile string, logFile string, imgFile string, imgDump bool, profileCapture []string, variables map[string]interface{}) error {
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
	var fileData []byte
	if fileData, err = os.ReadFile(e.cfgFile); err != nil {
		return err
	}
	if err = json.Unmarshal(fileData, &e.cfg); err != nil {
		return err
	}

	if e.cfg.ProfileCapture != nil {
		profileCapture = e.cfg.ProfileCapture
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
	if e.cfg.Debug {
		e.adapter.EnableDebug(true)
	}
	if err = e.adapter.Setup(logFile, imgFile, imgDump, profileCapture, e.cfg.ProfileSupportedMethods); err != nil {
		return err
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

func (e *Executor) computeScroll(value string) (int, []string, int) {
	var coords []string
	step := 0
	interval := 0
	if strings.Contains(value, ",") {
		container := strings.Split(value, ",")
		if len(container) > 0 {
			coords = strings.Split(container[0], "x")
		}
		if len(container) > 1 {
			step, _ = strconv.Atoi(container[1])
		}
		if len(container) > 2 {
			interval, _ = strconv.Atoi(container[2])
		}
	} else {
		coords = strings.Split(value, "x")
	}
	return step, coords, interval
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
			if strings.TrimSpace(strings.ToLower(value)) == "finalize" {
				event := e.adapter.CreateEvent(t.Id, nil, eventKindFull, t.Start, false)
				t.Stop()
				event.Write(eventMessageOK, t.Finalize())
			} else {
				t.Stop()
			}
			return nil
		}
		return fmt.Errorf("unknown timer id (timerStop): %s", target)
	case "timerFinalize":
		if t, ok := e.timers[target]; ok {
			event := e.adapter.CreateEvent(t.Id, nil, eventKindFull, t.Start, false)
			event.Write(eventMessageOK, t.Finalize())
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
	return e.adapter.MouseActions(by, data, command, value)
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
		return e.adapter.SendKeys(by, data, value)
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

func (e *Executor) doClose(cmd ConfigCommand) error {
	target := cmd.Target
	handle, err := e.adapter.GetWindowHandle(target)
	if err != nil {
		return err
	}
	if p := e.adapter.CloseWindow(handle); p == nil {
		return fmt.Errorf("invalid handle")
	}
	return nil
}

func (e *Executor) doStackAdd(id string, target string, u string) error {
	until := 0
	if len(u) > 0 {
		until = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	if elm := e.adapter.FindElementReady(by, data, until); elm != nil {
		values := make(map[string]interface{})
		values["target"] = target
		values["displayed"], _ = elm.IsDisplayed()
		values["enabled"], _ = elm.IsEnabled()
		values["text"], _ = elm.Text()
		values["tagName"], _ = elm.TagName()
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

func (e *Executor) doAssert(target string, u string, caption string) error {
	until := 0
	if len(u) > 0 {
		until = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	elm := e.adapter.FindElementReady(by, data, until)
	if elm == nil {
		return fmt.Errorf("element not exists")
	}
	text, err := elm.Text()
	if err != nil {
		return err
	}
	if text != caption {
		return fmt.Errorf("different text %s -> %s", text, caption)
	}
	return nil
}

func (e *Executor) doExists(target string, u string) error {
	until := 0
	if len(u) > 0 {
		until = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	if elm := e.adapter.FindElementReady(by, data, until); elm == nil {
		return fmt.Errorf("element not exists")
	}
	return nil
}

func (e *Executor) doUntil(cmd ConfigCommand) error {
	target := cmd.Target
	u := cmd.Until
	until := 0
	if len(u) > 0 {
		until = 1
	}
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	if elm := e.adapter.FindElementReady(by, data, until); elm != nil {
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

func (e *Executor) doStoreWindowHandle(cmd ConfigCommand) error {
	target := cmd.Target
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

func (e *Executor) doPause(cmd ConfigCommand) error {
	target := cmd.Target
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

func (e *Executor) doStatus(cmd ConfigCommand) error {
	id := cmd.Id
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
	//TODO TEST
	//for x := 0; x < 100; x++ {
	//	fmt.Println("SENDING KEYDOWN TEST")
	//	e.adapter.KeyDown(base.DownArrowKey)
	//	time.Sleep(100 * time.Millisecond)
	//}

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

func (e *Executor) doExecId(cmd ConfigCommand) error {
	target := cmd.Target
	e.execId = target
	return nil
}

func (e *Executor) doProbeId(cmd ConfigCommand) error {
	target := cmd.Target
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

func (e *Executor) doScroll(cmd ConfigCommand) error {
	value := cmd.Value
	target := cmd.Target
	step, coords, interval := e.computeScroll(value)
	if len(coords) < 2 {
		return fmt.Errorf("invalid value")
	}
	x := cmd.X
	y := cmd.Y
	var elm base.IWebElement
	if len(target) > 0 {
		by, data, err := e.computeSelector(target)
		if err != nil {
			return err
		}
		if elm = e.adapter.FindElementReady(by, data, 2); elm == nil {
			return fmt.Errorf("element isn't ready")
		}
		p, err := elm.Location()
		if err != nil {
			return err
		}
		x = p.X
		y = p.Y
		//TODO COMPLETARE
		//y += 50
		//x = 1050
		//y = 50
	}

	x += cmd.OffsetX
	y += cmd.OffsetY

	dx, _ := strconv.Atoi(coords[0])
	dy, _ := strconv.Atoi(coords[1])
	if step == 0 {
		return e.adapter.Scroll(int(x), int(y), dx, dy)
	}
	deltaX := dx / step
	deltaY := dy / step
	currentX := 0
	currentY := 0
	for k := 0; k <= step; k++ {
		if err := e.adapter.Scroll(int(x), int(y), currentX, currentY); err != nil {
			return err
		}
		currentX += deltaX
		currentY += deltaY
		if interval > 0 {
			e.adapter.Sleep(interval)
		}
	}
	return nil

}

func (e *Executor) doScrollTo(cmd ConfigCommand) error {
	value := cmd.Value
	target := cmd.Target
	step, coords, interval := e.computeScroll(value)
	if len(coords) < 2 {
		return fmt.Errorf("invalid value")
	}
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
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])
	x += int(cmd.OffsetX)
	y += int(cmd.OffsetY)

	if step == 0 {
		return elm.ScrollTo(x, y)
	}
	deltaX := x / step
	deltaY := y / step
	currentX := 0
	currentY := 0
	for k := 0; k <= step; k++ {
		if err := elm.ScrollTo(currentX, currentY); err != nil {
			return err
		}
		currentX += deltaX
		currentY += deltaY
		if interval > 0 {
			e.adapter.Sleep(interval)
		}
	}
	return nil
}

func (e *Executor) doNavigate(target string) error {
	if strings.HasPrefix(target, "/") {
		target = e.baseUrl + target
	}
	return e.adapter.Navigate(target)
}

func (e *Executor) commandExec(cmd ConfigCommand) error {
	start := time.Now()
	var err error
	id := cmd.Id
	command := cmd.Command
	target := cmd.Target
	until := cmd.Until
	value := cmd.Value
	windowTimeout := 0

	if len(cmd.WindowHandleName) > 0 {
		windowTimeout = e.adapter.AddWindowHandle(cmd)
	}

	e.adapter.HumanWait()

	switch command {
	case "execId":
		err = e.doExecId(cmd)
	case "probe":
		err = e.doProbeId(cmd)
	case "until":
		err = e.doUntil(cmd)
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
		err = e.doStoreWindowHandle(cmd)
	case "close":
		err = e.doClose(cmd)
	case "status":
		err = e.doStatus(cmd)
	case "pause":
		err = e.doPause(cmd)
	case "scrollTo":
		err = e.doScrollTo(cmd)
	case "scroll":
		err = e.doScroll(cmd)
	case "noop":
		//nothing to do
	default:
		err = fmt.Errorf("unimplemented command: %s", command)
		//e.adapter.log(LogLevelWarning, "unimplemented command: "+command)
	}

	if err == nil {
		if windowTimeout > 0 {
			e.adapter.Sleep(windowTimeout)
		}
	}

	event := e.adapter.CreateEvent(e.execId, err, eventKindIntermediate, start, true)
	event.Label = id
	event.Target = target
	event.Command = command
	event.Write("", UnixMilli(time.Now())-UnixMilli(start))
	return err
}

func (e *Executor) finalize(err error) {
	for _, t := range e.timers {
		if !t.Finalized {
			event := e.adapter.CreateEvent(t.Id, err, eventKindFull, t.Start, false)
			event.Write(eventMessageTimerNotFinalized, t.Finalize())
		}
	}
	msg := eventMessageOK
	if err != nil {
		msg = eventMessageNOK
	}
	event := e.adapter.CreateEvent(e.execId, err, eventKindFull, e.start, false)
	event.Write(msg, UnixMilli(time.Now())-UnixMilli(e.start))

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

func (e *Executor) doCommand(commands []ConfigCommand, idx int) (string, int, error) {
	if idx >= len(commands) {
		return "", -1, fmt.Errorf("invalid index %d", idx)
	}
	cmd := commands[idx]
	jump := -1
	if e.adapter.IsDebugEnabled() {
		fmt.Printf("--------------------------------------------------------------\n")
		fmt.Printf("Next command is [%s] %s: %s\n", cmd.Id, cmd.Command, cmd.Target)
	}
	err := e.templates.BuildCommand(cmd, e.stack)
	if err != nil {
		return cmd.Id, -1, err
	}
	if cmd.Command == "jump" {
		var j int
		if j, err = e.doComputeJump(cmd.Target, commands); err == nil {
			jump = j
		}
	} else {
		if e.lastFrame != nil {
			////if frameErr := e.adapter.SwitchParentFrame(); frameErr != nil {
			////	e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchParentFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
			////}
			//if frameErr := e.adapter.SwitchFrame(""); frameErr != nil {
			//	e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
			//}
			//if frameErr := e.adapter.SwitchFrame(e.lastFrame); frameErr != nil {
			//	e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
			//}
		}
		err = e.commandExec(cmd)
	}
	if e.adapter.IsDebugEnabled() {
		fmt.Printf("Error: %v\n", err)
	}
	if e.adapter.IsErrorDisabled() {
		err = nil
	}
	if err != nil {
		return cmd.Id, -1, err
	}
	return "", jump, err
}

func (e *Executor) commandsLoop() (string, error) {
	for _, test := range e.cfg.Tests {
		for x := 0; x < len(test.Commands); x++ {
			id, jump, err := e.doCommand(test.Commands, x)
			if err != nil {
				return id, err
			}
			if jump >= 0 {
				x = jump
			}
			/*
				cmd := test.Commands[x]
				if e.adapter.IsDebugEnabled() {
					fmt.Printf("--------------------------------------------------------------\n")
					fmt.Printf("Next command is [%s] %s: %s\n", cmd.Id, cmd.Command, cmd.Target)
				}
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
					if e.lastFrame != nil {
						////if frameErr := e.adapter.SwitchParentFrame(); frameErr != nil {
						////	e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchParentFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
						////}
						//if frameErr := e.adapter.SwitchFrame(""); frameErr != nil {
						//	e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
						//}
						//if frameErr := e.adapter.SwitchFrame(e.lastFrame); frameErr != nil {
						//	e.adapter.log(LogLevelWarning, fmt.Sprintf("SwitchFrame error :%s (last frame: %v)", frameErr.Error(), e.lastFrame))
						//}
					}
					err = e.commandExec(cmd)
				}
				if e.adapter.IsDebugEnabled() {
					fmt.Printf("Error: %v\n", err)
				}
				if e.adapter.IsErrorDisabled() {
					err = nil
				}
				if err != nil {
					return cmd.Id, err
				}
			*/
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
