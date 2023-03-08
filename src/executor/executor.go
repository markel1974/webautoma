package executor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/markel1974/webautoma/src/wd/base"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Executor struct {
	driver              base.IWebDriver
	fileName            string
	fileData            []byte
	fileId              string
	start               time.Time
	wait                int
	retryInterval       int
	humanWaitBase       int
	quit                bool
	uuid                string
	useHtmlEncoding     bool
	bid                 string
	maxScreenshotLength int
	errorDisabled       bool
	debug               bool
	maxRetry            int
	timers              map[string]*Timer
	probeId             string
	networkFilter       *regexp.Regexp
	drag                base.IWebElement
	mainWindow          string
	cfg                 Config
}

func New() *Executor {
	e := &Executor{}
	//e.driver = vars.getObject('driver');
	//e.actions = vars.getObject('actions');
	//e.filename = vars.getObject('filename');
	//e.filedata = vars.getObject('filedata');
	//e.imageAdapter = vars.getObject('imageAdapter');
	//e.sidecarWriter = vars.getObject('sidecarWriter');
	//e.keys = JavaImporter(org.openqa.selenium.Keys);
	//e.logPkg = JavaImporter(org.openqa.selenium.logging);
	e.retryInterval = 1000
	e.humanWaitBase = 150
	e.uuid = uuid.New().String()
	e.useHtmlEncoding = true
	e.bid = ""
	e.maxScreenshotLength = 0
	e.errorDisabled = false
	e.debug = false
	e.maxRetry = 3
	e.timers = make(map[string]*Timer)
	e.probeId = ""
	e.networkFilter = regexp.MustCompile("(http|file|ftp|png|jpg|gif|js|css|mp4|ico|bmp)")
	return e
}

func (e *Executor) RequiredLogs() (base.LogType, base.LogLevel) {
	return base.LogPerformance, base.LogAll
}

func (e *Executor) Setup(driver base.IWebDriver, fileName string) error {
	e.driver = driver
	e.fileName = fileName
	e.fileId = e.computeFileId(e.fileName)
	e.fileData = []byte(stub)
	if err := json.Unmarshal(e.fileData, &e.cfg); err != nil {
		return err
	}
	if e.cfg.Timeout != nil {
		e.wait = *e.cfg.Timeout
	} else {
		e.wait = 60000
	}
	if e.cfg.Quit != nil {
		e.quit = *e.cfg.Quit
	} else {
		e.quit = false
	}
	return nil
}

func (e *Executor) createEvent(id string, kind string, err error, start time.Time, dur int64, shot bool) map[string]interface{} {
	stop := time.Now()
	networkLog, errorCount := e.getNetworkLogs()
	errorDesc := ""
	if err != nil {
		errorDesc = err.Error()
	}
	event := map[string]interface{}{
		"probeId":           e.probeId,
		"thread_name":       id,
		"transaction_type":  kind,
		"timestamp":         time.Now().Format(RFC3339Milli),
		"start":             start.Format(RFC3339Milli),
		"stop":              stop.Format(RFC3339Milli),
		"execution_time":    dur,
		"uuid":              e.uuid,
		"livello":           "INFO",
		"error":             errorDesc,
		"passed":            err == nil,
		"business_id":       e.bid,
		"network":           networkLog,
		"networkErrorCount": errorCount,
		"screenshotId":      nil,
		"label":             nil,
		"target":            nil,
		"command":           nil,
		"message":           nil,
	}

	if shot {
		//480p = 858 x 480 - 720p = 1280 x 720 - 1080p = 1920 x 1080 — FullHD
		screenshotId := "probes-" + uuid.New().String()
		if screenshot, err := e.driver.Screenshot(); err == nil {
			if scaled, err := Scale(screenshot, 1920, 1080); err == nil {
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
					"id":          screenshotId,
					"length":      len(screenshotData),
					"screenshot":  screenshotData,
					"business_id": e.bid,
				}
				imageData, _ := json.Marshal(imageEvent)

				//TODO WRITE PNG!!!!
				fmt.Println("PNG DATA: ", string(imageData))
				//e.sidecarWriter.Write("images.log", imageData)
			}
		}
		event["screenshotId"] = screenshotId
	}
	return event
}

func (e *Executor) doLog(message string) {
	//TODO LOG
	fmt.Println(message)
	//WDS.log.info(message)
}

func (e *Executor) logDebug(message string) {
	if e.debug {
		e.doLog(message)
	}
}

func (e *Executor) logInfo(message string) {
	e.doLog(message)
}

func (e *Executor) logEvent(event map[string]interface{}) {
	message, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	e.doLog(string(message))
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

func (e *Executor) computeFileId(filename string) string {
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

func (e *Executor) getNetworkLogs() (map[string]interface{}, int) {
	logEntries, err := e.driver.Log(base.LogPerformance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, 0
	}

	errorCount := 0
	var headersData []map[string]interface{}

	for _, entry := range logEntries {
		var message NetworkMessage
		if err := json.Unmarshal([]byte(entry.Message), &message); err != nil {
			fmt.Println(err.Error())
			continue
		}
		if message.Method != "Network.responseReceived" {
			continue
		}

		if e.networkFilter.MatchString(message.Params.Response.Url) {
			headers := message.Params.Response.Headers
			headers["Url"] = message.Params.Response.Url
			headers["Status"] = message.Params.Response.Status
			headers["Timing"] = message.Params.Response.Timing
			contentType, _ := MapToString(headers, "Content-Type")

			if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "json") {
				if len(e.bid) == 0 {
					if bid, ok := MapToString(headers, "businessID"); ok {
						e.bid = bid
					}
				}
				if status, _ := MapToFloat64(headers, "status"); status >= 400 {
					errorCount++
				}
				for key := range headers {
					if key != "url" && key != "status" && key != "businessID" {
						delete(headers, key)
					}
				}
				headersData = append(headersData, map[string]interface{}{"headers": headers})
			}
		}
	}
	network := map[string]interface{}{
		"errorCount": errorCount,
		"data":       headersData,
	}
	return network, errorCount
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

	start := time.Now().UnixMilli()
	for time.Now().UnixMilli()-start < int64(e.wait) {
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
		e.logDebug("findElement: " + err.Error())
	}
	return elm
}

func (e *Executor) isElementReady(elm base.IWebElement) bool {
	ok, err := elm.IsEnabled()
	if err != nil {
		e.logDebug("isElementReady (IsEnabled): " + err.Error())
		return false
	}
	if !ok {
		e.logDebug("isElementReady: element isn't enabled")
		return false
	}
	ok, err = elm.IsDisplayed()
	if err != nil {
		e.logDebug("isElementReady (IsDisplayed): " + err.Error())
		return false
	}
	if !ok {
		e.logDebug("isElementReady: element isn't displayed")
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
			event["message"] = "Esito sonda ok"
			e.logEvent(event)
			return nil
		}
		return fmt.Errorf("unknown timer id (timerFinalize): %s", target)
	}
	return nil
}

func (e *Executor) waitForMouse(target string, until string, value string) error {
	var err error
	var start = time.Now().UnixMilli()
	var elm = e.getElementByMode(target, 2)
	if elm == nil {
		return fmt.Errorf("element isn't ready")
	}
	for time.Now().UnixMilli()-start < int64(e.wait) {
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
	start := time.Now().UnixMilli()
	pos := 0

	for time.Now().UnixMilli()-start < int64(e.wait) {
		err = nil
		if err = elm.SendKeys(string(data[pos])); err != nil {
			if elm = e.getElementByMode(target, 2); elm == nil {
				e.logDebug("waitForType: element isn't ready, returning...")
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
	start := time.Now().UnixMilli()
	for time.Now().UnixMilli()-start < int64(e.wait) {
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
		e.logInfo("unimplemented command: " + command)
	}
	if e.errorDisabled {
		err = nil
	}
	var dur = time.Now().UnixMilli() - start.UnixMilli()
	var event = e.createEvent(e.fileId, "intermediate", err, start, dur, err != nil)
	event["label"] = id
	event["target"] = target
	event["command"] = command
	e.logEvent(event)

	return err
}

func (e *Executor) finalize(err error) {
	for _, t := range e.timers {
		if !t.Finalized {
			start := t.Start
			dur := t.Finalize()
			event := e.createEvent(t.Id, "full", nil, start, dur, false)
			event["message"] = "timer not finalized"
			e.logEvent(event)
		}
	}

	var dur = time.Now().UnixMilli() - e.start.UnixMilli()
	var event = e.createEvent(e.fileId, "full", err, e.start, dur, false)
	if err != nil {
		event["message"] = "Errore nella sonda"
	} else {
		event["message"] = "Esito sonda ok"
	}

	e.logEvent(event)

	//WDS.sampleResult.sampleEnd()
	if e.quit {
		_ = e.driver.Quit()
	}
}

func (e *Executor) Run() error {
	var err error

	e.start = time.Now()

	//WDS.sampleResult.sampleStart()

	e.logInfo("Sample started")

	//var logTypes = e.driver.manage().logs().getAvailableLogTypes()
	//e.logInfo(logTypes)

	if err := e.loadUrl(e.cfg.Url); err != nil {
		e.finalize(err)
		return err
	}

	e.mainWindow, err = e.driver.CurrentWindowHandle()
	if err != nil {
		e.finalize(err)
		return err
	}

	e.logInfo("mainWindow: " + e.mainWindow)

	for _, test := range e.cfg.Tests {
		for _, cmd := range test.Commands {
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
