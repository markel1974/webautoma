package executor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/markel1974/webautoma/src/wd/base"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Loglevel int

const (
	LogLevelDebug    Loglevel = -1
	LogLevelInfo     Loglevel = 0
	LogLevelWarning  Loglevel = 1
	LogLevelCritical Loglevel = 2
)

const (
	DefaultResultFile = "log.json"
	DefaultImagesFile = "images.json"
)

type Executor struct {
	cfgFile             string
	logFile             string
	imgFile             string
	driver              base.IWebDriver
	fileId              string
	execId              string
	start               time.Time
	wait                int
	retryInterval       int
	humanWaitBase       int
	quit                bool
	uuid                string
	useHtmlEncoding     bool
	maxScreenshotLength int
	errorDisabled       bool
	debug               bool
	maxRetry            int
	timers              map[string]*Timer
	probeId             string
	drag                base.IWebElement
	mainWindow          string
	network             *Network
	templates           *Templates
	cfg                 Config
	stack               map[string]interface{}
	imgDump             int
	imgWidth            int
	imgHeight           int
	windowHandles       *Windows
}

func New() *Executor {
	e := &Executor{
		logFile:             DefaultResultFile,
		imgFile:             DefaultImagesFile,
		network:             nil,
		templates:           nil,
		driver:              nil,
		fileId:              "",
		quit:                false,
		retryInterval:       1000,
		humanWaitBase:       150,
		wait:                60000,
		uuid:                uuid.New().String(),
		useHtmlEncoding:     true,
		maxScreenshotLength: 0,
		errorDisabled:       false,
		debug:               false,
		maxRetry:            3,
		timers:              make(map[string]*Timer),
		probeId:             "",
		stack:               nil,
		imgDump:             0,
		imgWidth:            1920,
		imgHeight:           1080,
		windowHandles:       NewWindows(),
	}
	return e
}

func (e *Executor) SetImageDump() {
	e.imgDump = 1
}

func (e *Executor) SetImageSize(w int, h int) {
	if w > 0 && h > 0 {
		e.imgWidth = w
		e.imgHeight = h
	}
}

func (e *Executor) SetLogFile(logFile string) {
	if len(logFile) > 0 {
		e.logFile = logFile
	}
}

func (e *Executor) SetImageFile(imgFile string) {
	if len(imgFile) > 0 {
		e.imgFile = imgFile
	}
}

func (e *Executor) Setup(fileName string, capture []string, variables map[string]interface{}) error {
	var err error

	e.network = NewNetwork(capture)
	e.templates = NewTemplates(variables)
	e.fileId = computeFileId(fileName)
	e.execId = e.fileId

	if e.logFile, err = e.templates.Apply(e.logFile, nil); err != nil {
		return err
	}
	if e.imgFile, err = e.templates.Apply(e.imgFile, nil); err != nil {
		return err
	}
	if e.cfgFile, err = e.templates.Apply(fileName, nil); err != nil {
		return err
	}
	fileData, err := os.ReadFile(e.cfgFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(fileData, &e.cfg); err != nil {
		return err
	}
	if e.cfg.Timeout != nil {
		e.wait = *e.cfg.Timeout
	}
	if e.cfg.Quit != nil {
		e.quit = *e.cfg.Quit
	}
	if e.cfg.HumanWait != nil {
		e.humanWaitBase = *e.cfg.HumanWait
	}
	if e.cfg.RetryInterval != nil {
		e.retryInterval = *e.cfg.RetryInterval
	}
	if e.cfg.MaxScreenshotLength != nil {
		e.maxScreenshotLength = *e.cfg.MaxScreenshotLength
	}
	return nil
}

func (e *Executor) RequiredLogs() (base.LogType, base.LogLevel) {
	return base.LogPerformance, base.LogAll
}

func (e *Executor) createEvent(id string, kind string, err error, start time.Time, dur int64, shot bool) *Event {
	if e.errorDisabled {
		err = nil
	}
	event := NewEvent(id, kind, err, start, dur)
	event.ProbeId = e.probeId
	event.UUID = e.uuid
	logEntries, logErr := e.driver.Log(base.LogPerformance)
	if logErr != nil {
		e.log(LogLevelCritical, logErr.Error())
	}
	event.Network, event.NetworkErrorCount = e.network.Compute(logEntries)
	if acquired := e.network.Acquired(); acquired != nil {
		event.Acquired = acquired
	}
	if shot && err != nil {
		event.ScreenShoot = "probes-" + uuid.New().String()
		if err := e.doScreenshot(event.ScreenShoot); err != nil {
			e.log(LogLevelCritical, "error generating screenshot: "+err.Error())
		}
	}
	return event
}

func (e *Executor) loadUrl(url string) error {
	if err := e.driver.Navigate(url); err != nil {
		return err
	}
	return nil
}

func (e *Executor) doScreenshot(screenshotId string) error {
	//480p = 858 x 480 - 720p = 1280 x 720 - 1080p = 1920 x 1080 — FullHD
	screenshot, err := e.driver.Screenshot()
	if err != nil {
		return err
	}
	scaled, err := Scale(screenshot, e.imgWidth, e.imgHeight, true)
	if err != nil {
		return err
	}

	if e.imgDump > 0 {
		imgFile := time.Now().Format(RFC3339Milli) + "." + strconv.Itoa(e.imgDump) + ".png"
		imgFile = strings.Replace(imgFile, ":", "-", -1)
		if fi, fiErr := os.OpenFile(imgFile, os.O_CREATE|os.O_WRONLY, 0644); fiErr == nil {
			_, _ = fi.Write(scaled)
			fi.Close()
		}
		e.imgDump++
	}

	screenshotData := base64.StdEncoding.EncodeToString(scaled)
	if e.useHtmlEncoding {
		screenshotData = "data:image/png;base64," + screenshotData
	}
	if e.maxScreenshotLength > 0 {
		if len(screenshotData) > e.maxScreenshotLength {
			screenshotData = screenshotData[0:e.maxScreenshotLength]
		}
	}
	imageEvent := map[string]interface{}{
		"id":         screenshotId,
		"length":     len(screenshotData),
		"screenshot": screenshotData,
	}
	if acquired := e.network.Acquired(); acquired != nil {
		imageEvent["acquired"] = acquired
	}
	imageData, err := json.Marshal(imageEvent)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(e.imgFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, _ = f.Write(imageData)
	_, _ = f.Write([]byte("\n"))
	_ = f.Sync()
	f.Close()

	return nil
}

func (e *Executor) doWriteLog(message string) {
	f, err := os.OpenFile(e.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		e.log(LogLevelCritical, "error writing log file: "+err.Error())
		return
	}
	_, _ = f.Write([]byte(message))
	_, _ = f.Write([]byte("\n"))
	_ = f.Sync()
	f.Close()
}

func (e *Executor) log(kind Loglevel, src string) {
	if kind == -1 && !e.debug {
		return
	}
	var severity string
	switch kind {
	case LogLevelDebug:
		severity = "debug"
	case LogLevelInfo:
		severity = "info"
	case LogLevelWarning:
		severity = "warning"
	case LogLevelCritical:
		severity = "critical"
	default:
		severity = "info"
	}
	event := map[string]interface{}{"severity": severity, "data": src}
	message, _ := json.Marshal(event)
	e.doWriteLog(string(message))
}

func (e *Executor) logEvent(event *Event) {
	message, err := json.Marshal(event)
	if err != nil {
		e.log(LogLevelCritical, err.Error())
		return
	}
	e.doWriteLog(string(message))
}

func (e *Executor) sleep(interval int) {
	time.Sleep(time.Millisecond * time.Duration(interval))
}

func (e *Executor) humanWait() {
	rnd := rand.Float64()
	val := math.Round(rnd * 100)
	interval := e.humanWaitBase + int(val)
	e.sleep(interval)
}

func (e *Executor) getElementByMode(target string, until int) base.IWebElement {
	var container = strings.Split(target, "=")
	if len(container) <= 1 {
		return nil
	}
	var elem base.IWebElement = nil
	var search string
	var data = strings.TrimSpace(container[1])
	var command = strings.TrimSpace(container[0])
	switch command {
	case "xpath":
		search = base.ByXPATH // selector.xpath(data)
	case "id":
		search = base.ByID //selector.id(data)
	case "linkText":
		search = base.ByLinkText // selector.linkText(data)
	case "name":
		search = base.ByName //selector.name(data)
	case "css":
		search = base.ByCSSSelector // selector.cssSelector(data)
	case "class":
		search = base.ByClassName // selector.cssSelector(data)
	default:
		return nil
	}

	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(e.wait) {
		found := false
		elem = e.findElement(search, data)
		if until == 0 {
			found = elem != nil
		} else if until == 1 {
			if elem != nil {
				if !e.isElementReady(target, elem) {
					found = true
					elem = nil
				}
			} else {
				found = true
			}
		} else {
			if elem != nil {
				if e.isElementReady(target, elem) {
					found = true
				}
			}
		}
		if found {
			break
		} else {
			e.sleep(e.retryInterval)
		}
	}

	return elem
}

func (e *Executor) getCoords(data string) base.Point {
	var coords base.Point
	if container := strings.Split(data, ","); len(container) > 1 {
		coords.X, _ = strconv.ParseFloat(container[0], 64)
		coords.Y, _ = strconv.ParseFloat(container[1], 64)
	}
	return coords
}

func (e *Executor) findElement(by string, data string) base.IWebElement {
	elm, err := e.driver.FindElement(by, data)
	if err != nil {
		e.log(LogLevelDebug, "findElement: "+err.Error())
	}
	return elm
}

func (e *Executor) isElementReady(target string, elm base.IWebElement) bool {
	ok, err := elm.IsEnabled()
	if err != nil {
		e.log(LogLevelDebug, fmt.Sprintf("(%s) isElementReady : %s", target, err.Error()))
		return false
	}
	if !ok {
		e.log(LogLevelDebug, fmt.Sprintf("(%s) isElementReady: element isn't enabled", target))
		return false
	}
	ok, err = elm.IsDisplayed()
	if err != nil {
		e.log(LogLevelDebug, fmt.Sprintf("(%s) isElementReady (IsDisplayed): %s", target, err.Error()))
		return false
	}
	if !ok {
		e.log(LogLevelDebug, fmt.Sprintf("(%s) isElementReady: element isn't displayed", target))
		return false
	}
	return true
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
			event := e.createEvent(t.Id, "full", nil, start, dur, false)
			event.Message = "Esito sonda ok"
			e.logEvent(event)
			return nil
		}
		return fmt.Errorf("unknown timer id (timerFinalize): %s", target)
	}
	return nil
}

func (e *Executor) doMouse(target string, until string, value string) error {
	var err error
	var start = UnixMilli(time.Now())
	var elm = e.getElementByMode(target, 2)
	if elm == nil {
		return fmt.Errorf("element isn't ready (%s)", target)
	}
	for UnixMilli(time.Now())-start < int64(e.wait) {
		err = nil
		switch until {
		case "click":
			err = elm.Click()

		case "doubleClick":
			//click - sleep - click
			//err = elm.Click()
			//e.doSleep(100)
			//err = elm.Click()
			if err = elm.MoveTo(0, 0); err == nil {
				err = e.driver.DoubleClick()
			}
			//e.actions.doubleClick(elm).perform()

		case "rightClick":
			if err = elm.MoveTo(0, 0); err == nil {
				err = e.driver.Click(2)
			}
			//e.actions.contextClick(elm).perform()

		case "mouseOver":
			err = elm.MoveTo(0, 0)
			//elm.MoveTo(0, 0)
			//e.actions.moveToElement(elm).perform()

		case "mouseUpAt":
			if e.drag != nil {
				var coordsUpAt = e.getCoords(value)
				if err = e.drag.MoveTo(coordsUpAt.X, coordsUpAt.Y); err != nil {
					err = e.driver.ButtonUp()
				}
				//e.drag.moveByOffset(coordsUpAt.X, coordsUpAt.Y).release().perform()
				e.drag = nil
			}

		case "mouseDownAt":
			var coordsDownAt = e.getCoords(value)
			//This is equivalent to: Actions.moveToElement(onElement).clickAndHold()
			//e.drag = e.actions.clickAndHold(elm)
			//e.drag.moveByOffset(coordsDownAt.X, coordsDownAt.Y).perform()

			if err = elm.MoveTo(coordsDownAt.X, coordsDownAt.Y); err != nil {
				err = e.driver.ButtonDown()
				e.drag = elm
			}

		case "mouseMultipleMoveAt":
			if e.drag != nil {
				for _, m := range strings.Split(value, "|") {
					var coordsMultipleMoveAt = e.getCoords(m)
					if err = e.drag.MoveTo(coordsMultipleMoveAt.X, coordsMultipleMoveAt.Y); err != nil {
						break
					}
					//e.drag.moveByOffset(coordsMultipleMoveAt.X, coordsMultipleMoveAt.Y).perform()
					e.humanWait()
				}
			}
		case "mouseMoveAt":
			if e.drag != nil {
				var coordsMoveAt = e.getCoords(value)
				err = e.drag.MoveTo(coordsMoveAt.X, coordsMoveAt.Y)
				//e.drag.moveByOffset(coordsMoveAt.x, coordsMoveAt.y).perform()
			}
			/*
				} catch (e) {
				var exception = e.toString().toLowerCase();
					if (!exception.contains('staleelementreferenceexception')) {
					this.logDebug('waitForMouse: ' + e);
					elm = this.getElementByMode(target, 2);
					if (!elm) {
					this.logDebug("waitForMouse: element isn't ready, returning...");
					return "element isn't ready";
					}
					err = e
					}
				}

			*/
		}
		if err == nil {
			break
		} else {
			e.sleep(e.retryInterval)
		}
	}
	return err
}

func (e *Executor) doType(target string, data string) error {
	if len(data) <= 0 {
		return nil
	}

	elm := e.getElementByMode(target, 2)
	if elm == nil {
		return fmt.Errorf("element isn't ready")
	}
	var err error
	start := UnixMilli(time.Now())
	pos := 0

	for UnixMilli(time.Now())-start < int64(e.wait) {
		err = nil
		if err = elm.SendKeys(string(data[pos])); err != nil {
			if elm = e.getElementByMode(target, 2); elm == nil {
				e.log(LogLevelDebug, fmt.Sprintf("(%s) waitForType: element isn't ready", target))
				return errors.New("element isn't ready")
			}
		}

		e.humanWait()

		if err == nil {
			pos++
			if pos >= len(data) {
				break
			}
		} else {
			e.sleep(e.retryInterval)
		}
	}
	return err
}

func (e *Executor) doRunScript(script string, args []interface{}) error {
	var err error
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(e.wait) {
		_, err = e.driver.ExecuteScript(script, args)
		if err == nil {
			break
		} else {
			e.sleep(e.retryInterval)
		}
	}
	return err
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
		case "relative":
			frame = v[1]
		default:
			frame = v[1]
		}
	} else {
		frame = ""
	}

	//var frame interface{}
	//if elm := e.getElementByMode(target, 2); elm != nil {
	//	frame = elm
	//} else {
	//	frame = target
	//}

	if err := e.driver.SwitchFrame(frame); err != nil {
		return err
	}
	return nil
}

func (e *Executor) doSelectParentFrame() error {
	if v := e.driver.SwitchParentFrame(); v == nil {
		return fmt.Errorf("invalid frame")
	}
	return nil
}

func (e *Executor) doSelectWindowMain() error {
	if v := e.driver.SwitchWindow(e.mainWindow); v == nil {
		return fmt.Errorf("invalid mainWindow")
	}
	return nil
}

func (e *Executor) doSelectWindowTitle(title string) error {
	currentWindow, err := e.driver.CurrentWindowHandle()
	if err != nil {
		return err
	}
	lowerTitle := strings.ToLower(title)
	var found = false

	windows, err := e.driver.WindowHandles()
	if err != nil {
		return err
	}

	for _, window := range windows {
		if err := e.driver.SwitchWindow(window); err != nil {
			continue
		}
		windowTitle, err := e.driver.Title()
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
		_ = e.driver.SwitchWindow(currentWindow)
		return fmt.Errorf("can't find window with title " + title)
	}
	return nil
}

func (e *Executor) doSelectAlert(val string) error {
	value := parseInt(val)
	_, err := e.driver.AlertText()
	if err != nil {
		return err
	}

	switch value {
	case 1:
		if err := e.driver.AcceptAlert(); err != nil {
			return err
		}
	case 2:
		//alert.read()
	default:
		if err := e.driver.DismissAlert(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) doCloseWindow() error {
	currentWindow, err := e.driver.CurrentWindowHandle()
	if err != nil {
		return err
	}
	if currentWindow == e.mainWindow {
		return fmt.Errorf("can't close mainWindow")
	}
	_ = e.driver.Close()
	if v := e.driver.SwitchWindow(e.mainWindow); v == nil {
		return fmt.Errorf("invalid mainWindow")
	}
	return nil
}

func (e *Executor) doClose(target string) error {
	handle, err := e.windowHandles.GetHandle(target)
	if err != nil {
		return err
	}
	if p := e.driver.CloseWindow(handle); p == nil {
		return fmt.Errorf("invalid handle")
	}
	return nil
}

func (e *Executor) doStackAdd(id string, target string, until string) error {
	v := 0
	if len(until) > 0 {
		v = 1
	}
	if analyzedElm := e.getElementByMode(target, v); analyzedElm != nil {
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
	existsElm := e.getElementByMode(target, v)
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
	if existsElm := e.getElementByMode(target, v); existsElm == nil {
		return fmt.Errorf("element not exists")
	}
	return nil
}

func (e *Executor) doUntil(target string, until string) error {
	v := 0
	if len(until) > 0 {
		v = 1
	}
	if untilElm := e.getElementByMode(target, v); untilElm != nil {
		return fmt.Errorf("element exists")
	}
	return nil
}

func (e *Executor) doSelectWindow(target string) error {
	handle, err := e.windowHandles.GetHandle(target)
	if err != nil {
		return err
	}
	if err := e.driver.SwitchWindow(handle); err != nil {
		return err
	}
	return nil
}

func (e *Executor) doStoreWindowHandle(target string) error {
	return e.windowHandles.Store(target)
}

func (e *Executor) doSelect(target string, value string) error {
	v := strings.Split(value, "=")
	if len(v) < 2 {
		return fmt.Errorf("invalid value")
	}
	by := v[0]
	if by != "label" {
		return fmt.Errorf("unsupported by")
	}
	label := v[1]
	elm := e.getElementByMode(target, 2)
	if elm == nil {
		return fmt.Errorf("element isn't ready")
	}
	z, err := elm.FindElements(base.ByTagName, "option")
	if err != nil {
		return err
	}
	for _, k := range z {
		computedLabel, err := k.ComputedLabel()
		if err != nil {
			continue
		}
		if computedLabel == label {
			return k.Click()
		}
	}
	return fmt.Errorf("not found")
}

func (e *Executor) doActiveElement(id string) error {
	v, err := e.driver.ActiveElement()
	if err != nil {
		return err
	}
	fmt.Printf("%s ActiveElement =>\n", id)
	v.Print()
	return nil
}

func (e *Executor) doWindowHandles(id string) error {
	fmt.Printf("%s WindowHandles =>\n", id)
	v, _ := e.driver.WindowHandles()
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
	e.wait = parseInt(value)
	return nil
}

func (e *Executor) doPageSource(id string) error {
	v, _ := e.driver.PageSource()
	fmt.Printf("%s PageSource =>\n", id)
	fmt.Println(v)
	return nil
}

func (e *Executor) doStatus(id string) error {
	fmt.Printf("%s Status =>\n", id)
	k, _ := e.driver.Status()
	fmt.Println(k)
	return nil
}

func (e *Executor) doGetCookie(target string, id string) error {
	v, _ := e.driver.GetCookie(target)
	fmt.Printf("%s Cookie =>\n", id)
	fmt.Println(v)
	return nil
}

func (e *Executor) doGetAllCookies(id string) error {
	v, _ := e.driver.GetCookies()
	fmt.Printf("%s Cookies =>\n", id)
	fmt.Println(v)
	return nil
}

func (e *Executor) doDeleteAllCookies() error {
	return e.driver.DeleteAllCookies()
}

func (e *Executor) doDeleteCookie(target string) error {
	return e.driver.DeleteCookie(target)
}

func (e *Executor) doActionsSendKeys(value string) error {
	return e.driver.KeyDown(value)
}

func (e *Executor) doSetWindowMain() error {
	var err error
	e.mainWindow, err = e.driver.CurrentWindowHandle()
	return err
}

func (e *Executor) doHumanWait(value string) error {
	e.sleep(parseInt(value))
	return nil
}

func (e *Executor) doEnableDebug() error {
	e.debug = true
	return nil
}

func (e *Executor) doDisableError() error {
	e.errorDisabled = true
	return nil
}

func (e *Executor) doEnableError() error {
	e.errorDisabled = false
	return nil
}

func (e *Executor) doDisableDebug() error {
	e.debug = false
	return nil
}

func (e *Executor) doExecId(target string) error {
	e.execId = target
	return nil
}

func (e *Executor) doProbeId(target string) error {
	e.probeId = target
	return nil
}

func (e *Executor) commandExec(id string, command string, target string, until string, value string) error {
	var err error
	start := time.Now()

	e.humanWait()

	switch command {
	case "execId":
		err = e.doExecId(target)
	case "probe":
		err = e.doProbeId(target)
	case "until":
		err = e.doUntil(target, until)
	case "open":
		//nothing to do
	case "setWindowSize":
		//nothing to do
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
	case "noop":
		//nothing to do
	default:
		e.log(LogLevelWarning, "unimplemented command: "+command)
	}
	event := e.createEvent(e.execId, "intermediate", err, start, UnixMilli(time.Now())-UnixMilli(start), true)
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
			event := e.createEvent(t.Id, "full", err, start, dur, false)
			event.Message = "timer not finalized"
			e.logEvent(event)
		}
	}
	var dur = UnixMilli(time.Now()) - UnixMilli(e.start)
	var event = e.createEvent(e.execId, "full", err, e.start, dur, false)
	if err != nil {
		event.Message = "Errore nella sonda"
	} else {
		event.Message = "Esito sonda ok"
	}
	e.logEvent(event)
	if e.quit {
		_ = e.driver.Quit()
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
				if len(cmd.WindowHandleName) > 0 {
					e.windowHandles.Add(e.driver, cmd)
				}
				err = e.commandExec(cmd.Id, cmd.Command, cmd.Target, cmd.Until, cmd.Value)
			}
			if e.errorDisabled {
				err = nil
			}
			if err != nil {
				return cmd.Id, err
			}
		}
	}
	return "", nil
}

func (e *Executor) Start(driver base.IWebDriver) error {
	e.driver = driver
	e.start = time.Now()

	e.log(LogLevelInfo, "Sample started")

	url, err := e.templates.Apply(e.cfg.Url, nil)
	if err != nil {
		return err
	}
	if err = e.loadUrl(url); err != nil {
		e.finalize(err)
		return err
	}
	if e.mainWindow, err = e.driver.CurrentWindowHandle(); err != nil {
		e.finalize(err)
		return err
	}

	e.log(LogLevelInfo, "mainWindow: "+e.mainWindow)

	var id string
	if id, err = e.commandsLoop(); err != nil {
		e.finalize(fmt.Errorf("%s (%s)", err.Error(), id))
		return err
	}

	e.finalize(nil)
	return nil
}
