package executor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type Executor struct {
	cfgFile             string
	logFile             string
	imgFile             string
	driver              base.IWebDriver
	fileId              string
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
}

func New(fileName string, logFile string, imgFile string, capture []string, variables map[string]interface{}) (*Executor, error) {
	e := &Executor{
		logFile:             "",
		imgFile:             "",
		network:             NewNetwork(capture),
		templates:           NewTemplates(variables),
		driver:              nil,
		fileId:              computeFileId(fileName),
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
	}
	var err error
	e.logFile, err = e.templates.Apply(logFile)
	if err != nil {
		return nil, err
	}
	e.imgFile, err = e.templates.Apply(imgFile)
	if err != nil {
		return nil, err
	}
	e.cfgFile, err = e.templates.Apply(fileName)
	if err != nil {
		return nil, err
	}
	fileData, err := ioutil.ReadFile(e.cfgFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fileData, &e.cfg); err != nil {
		return nil, err
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
	return e, nil
}

func (e *Executor) RequiredLogs() (base.LogType, base.LogLevel) {
	//TODO FROM CONFIG
	return base.LogPerformance, base.LogAll
}

func (e *Executor) createEvent(id string, kind string, err error, start time.Time, dur int64, shot bool) *Event {
	event := NewEvent(id, kind, err, start, dur)
	event.ProbeId = e.probeId
	event.UUID = e.uuid
	logEntries, err := e.driver.Log(base.LogPerformance)
	if err != nil {
		e.log(LogLevelCritical, err.Error())
	}
	event.Network, event.NetworkErrorCount = e.network.Compute(logEntries)
	event.Acquired = e.network.Acquired()
	if shot {
		event.ScreenShoot = "probes-" + uuid.New().String()
		if err := e.doScreenshot(event.ScreenShoot); err != nil {
			e.log(LogLevelCritical, "error generating screenshot: "+err.Error())
		}
	}
	return event
}

func (e *Executor) doScreenshot(screenshotId string) error {
	//480p = 858 x 480 - 720p = 1280 x 720 - 1080p = 1920 x 1080 — FullHD
	screenshot, err := e.driver.Screenshot()
	if err != nil {
		return err
	}
	scaled, err := Scale(screenshot, 1920, 1080)
	if err != nil {
		return err
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
		//"business_id": e.bid,
	}
	imageEvent["acquired"] = e.network.Acquired()
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

func (e *Executor) doSleep(interval int) {
	time.Sleep(time.Millisecond * time.Duration(interval))
}

func (e *Executor) humanWait() {
	rnd := rand.Float64()
	val := math.Round(rnd * 100)
	interval := e.humanWaitBase + int(val)
	e.doSleep(interval)
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
				if !e.isElementReady(elem) {
					found = true
					elem = nil
				}
			} else {
				found = true
			}
		} else {
			if elem != nil {
				if e.isElementReady(elem) {
					found = true
				}
			}
		}
		if found {
			break
		} else {
			e.doSleep(e.retryInterval)
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

func (e *Executor) isElementReady(elm base.IWebElement) bool {
	ok, err := elm.IsEnabled()
	if err != nil {
		e.log(LogLevelDebug, "isElementReady (IsEnabled): "+err.Error())
		return false
	}
	if !ok {
		e.log(LogLevelDebug, "isElementReady: element isn't enabled")
		return false
	}
	ok, err = elm.IsDisplayed()
	if err != nil {
		e.log(LogLevelDebug, "isElementReady (IsDisplayed): "+err.Error())
		return false
	}
	if !ok {
		e.log(LogLevelDebug, "isElementReady: element isn't displayed")
		return false
	}
	return true
}

func (e *Executor) timerHandler(target string, until string, value string) error {
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

func (e *Executor) waitForMouse(target string, until string, value string) error {
	var err error
	var start = UnixMilli(time.Now())
	var elm = e.getElementByMode(target, 2)
	if elm == nil {
		return fmt.Errorf("element isn't ready")
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
				//TODO TEST
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
			e.doSleep(e.retryInterval)
		}
	}
	return err
}

func (e *Executor) waitForType(target string, data string) error {
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
				e.log(LogLevelDebug, "waitForType: element isn't ready, returning...")
				return fmt.Errorf("element isn't ready")
			}
		}

		e.humanWait()

		if err == nil {
			pos++
			if pos >= len(data) {
				break
			}
		} else {
			e.doSleep(e.retryInterval)
		}
	}
	return err
}

func (e *Executor) waitForScript(script string, args []interface{}) error {
	var err error
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(e.wait) {
		_, err = e.driver.ExecuteScript(script, args)
		if err == nil {
			break
		} else {
			e.doSleep(e.retryInterval)
		}
	}
	return err
}

func (e *Executor) loadUrl(url string) error {
	if err := e.driver.Navigate(url); err != nil {
		return err
	}
	return nil
}

func (e *Executor) selectWindowMain() error {
	if v := e.driver.SwitchWindow(e.mainWindow); v == nil {
		return fmt.Errorf("invalid mainWindow")
	}
	return nil
}

func (e *Executor) selectWindowTitle(title string) error {
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

func (e *Executor) selectAlert(value int) error {
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

func (e *Executor) closeWindow() error {
	var currentWindow, err = e.driver.CurrentWindowHandle()
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

func (e *Executor) exec(id string, command string, target string, until string, value string) error {
	var err error
	start := time.Now()

	e.humanWait()

	switch command {
	case "probe":
		e.probeId = target
	case "until":
		v := 0
		if len(until) > 0 {
			v = 1
		}
		if untilElm := e.getElementByMode(target, v); untilElm != nil {
			err = fmt.Errorf("element exists")
		}
	case "open":
		//nothing to do
	case "setWindowSize":
		//nothing to do
	case "timerCreate", "timerStart", "timerStop", "timerFinalize":
		err = e.timerHandler(target, command, value)
	case "exists":
		v := 0
		if len(until) > 0 {
			v = 1
		}
		if existsElm := e.getElementByMode(target, v); existsElm == nil {
			err = fmt.Errorf("element not exists")
		}
	case "disableError":
		e.errorDisabled = true
	case "enableError":
		e.errorDisabled = false
	case "disableDebug":
		e.debug = false
	case "enableDebug":
		e.debug = true
	case "click", "rightClick", "doubleClick", "mouseUpAt", "mouseDownAt", "mouseMultipleMoveAt", "mouseMoveAt", "mouseOver":
		err = e.waitForMouse(target, command, value)
	case "mouseOut":
		//nothing to do
	case "type":
		err = e.waitForType(target, value)
	case "runScript":
		err = e.waitForScript(target, nil)
	case "humanWait":
		e.doSleep(parseInt(value))
	case "selectWindowMain":
		err = e.selectWindowMain()
	case "selectWindowTitle":
		err = e.selectWindowTitle(value)
	case "selectAlert":
		err = e.selectAlert(parseInt(value))
	case "closeWindow":
		err = e.closeWindow()
	case "actionsSendKeys":
		err = e.driver.KeyDown(value)
		//e.actions.sendKeys(value).perform()
	case "setTimeout":
		e.wait = parseInt(value)
	default:
		e.log(LogLevelWarning, "unimplemented command: "+command)
	}
	if e.errorDisabled {
		err = nil
	}
	var dur = UnixMilli(time.Now()) - UnixMilli(start)
	var event = e.createEvent(e.fileId, "intermediate", err, start, dur, err != nil)
	event.Label = id
	event.Target = target
	event.Target = command
	e.logEvent(event)

	return err
}

func (e *Executor) finalize(err error) {
	for _, t := range e.timers {
		if !t.Finalized {
			start := t.Start
			dur := t.Finalize()
			event := e.createEvent(t.Id, "full", nil, start, dur, false)
			event.Message = "timer not finalized"
			e.logEvent(event)
		}
	}

	var dur = UnixMilli(time.Now()) - UnixMilli(e.start)
	var event = e.createEvent(e.fileId, "full", err, e.start, dur, false)
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

func (e *Executor) Start(driver base.IWebDriver) error {
	var err error
	e.driver = driver
	e.start = time.Now()

	e.log(LogLevelInfo, "Sample started")

	url, err := e.templates.Apply(e.cfg.Url)
	if err != nil {
		return err
	}
	if err := e.loadUrl(url); err != nil {
		e.finalize(err)
		return err
	}
	e.mainWindow, err = e.driver.CurrentWindowHandle()
	if err != nil {
		e.finalize(err)
		return err
	}

	e.log(LogLevelInfo, "mainWindow: "+e.mainWindow)

	for _, test := range e.cfg.Tests {
		for _, cmd := range test.Commands {
			var err error
			if cmd.Id, err = e.templates.Apply(cmd.Id); err != nil {
				return err
			}
			if cmd.Command, err = e.templates.Apply(cmd.Command); err != nil {
				return err
			}
			if cmd.Target, err = e.templates.Apply(cmd.Target); err != nil {
				return err
			}
			if cmd.Until, err = e.templates.Apply(cmd.Until); err != nil {
				return err
			}
			if cmd.Value, err = e.templates.Apply(cmd.Value); err != nil {
				return err
			}
			if err = e.exec(cmd.Id, cmd.Command, cmd.Target, cmd.Until, cmd.Value); err != nil {
				e.finalize(err)
				return err
			}
		}
	}
	e.finalize(nil)
	return nil
}

func parseInt(in string) int {
	out, _ := strconv.Atoi(in)
	return out
}

func computeFileId(filename string) string {
	var pos = strings.LastIndex(filename, "\\")
	if pos < 0 {
		pos = strings.LastIndex(filename, "/")
	}
	if pos < 0 {
		pos = 0
	} else {
		pos++
	}
	name := filename[pos:]
	if pos := strings.LastIndex(name, "."); pos >= 0 {
		name = name[0:pos]
	}
	return name
}
