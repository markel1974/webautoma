package executor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/markel1974/webautoma/src/wd/base"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

const (
	DefaultResultFile = "log.json"
	DefaultImagesFile = "images.json"
)

type Loglevel int

const (
	LogLevelDebug    Loglevel = -1
	LogLevelInfo     Loglevel = 0
	LogLevelWarning  Loglevel = 1
	LogLevelCritical Loglevel = 2
)

type Adapter struct {
	driver              base.IWebDriver
	retryInterval       int
	logFile             string
	imgFile             string
	imgDump             int
	imgWidth            int
	imgHeight           int
	maxScreenshotLength int
	useHtmlEncoding     bool
	drag                base.IWebElement
	rootWindow          string
	wait                int
	humanWaitBase       int
	network             *Network
	errorDisabled       bool
	probeId             string
	uuid                string
	windowHandles       *Windows
	debug               bool
}

func NewAdapter(driver base.IWebDriver) *Adapter {
	return &Adapter{
		driver:              driver,
		retryInterval:       1000,
		logFile:             DefaultResultFile,
		imgFile:             DefaultImagesFile,
		imgDump:             0,
		imgWidth:            1920,
		imgHeight:           1080,
		useHtmlEncoding:     true,
		drag:                nil,
		rootWindow:          "",
		wait:                60000,
		maxScreenshotLength: 0,
		humanWaitBase:       150,
		network:             nil,
		errorDisabled:       false,
		probeId:             "",
		uuid:                uuid.New().String(),
		windowHandles:       NewWindows(),
		debug:               false,
	}
}

func (e *Adapter) Setup(logFile string, imgFile string, imgDump bool, profileCapture []string, profileSupportedMethod []string) error {
	if len(logFile) > 0 {
		e.logFile = logFile
	}
	if len(imgFile) > 0 {
		e.imgFile = imgFile
	}
	if imgDump {
		e.imgDump = 1
	}
	var err error
	e.network, err = NewNetwork(profileCapture, profileSupportedMethod)
	if err != nil {
		return err
	}
	return nil
}

func (e *Adapter) GetWindowHandle(id string) (string, error) {
	return e.windowHandles.GetHandle(id)
}

func (e *Adapter) StoreWindowHandle(id string) error {
	return e.windowHandles.Store(id)
}

func (e *Adapter) AddWindowHandle(cmd ConfigCommand) int {
	return e.windowHandles.Add(e.driver, cmd)
}

func (e *Adapter) Quit() error {
	return e.driver.Quit()
}

func (e *Adapter) GetRootWindow() string {
	return e.rootWindow
}

func (e *Adapter) SetRootWindow() error {
	handle, err := e.CurrentWindowHandle()
	if err != nil {
		return err
	}
	e.rootWindow = handle
	return err
}

func (e *Adapter) SetProbeId(id string) {
	e.probeId = id
}

func (e *Adapter) SetRetryInterval(interval int) {
	e.retryInterval = interval
}

func (e *Adapter) SetWait(wait int) {
	e.wait = wait
}

func (e *Adapter) SetHumanWait(wait int) {
	e.humanWaitBase = wait
}

func (e *Adapter) SetErrorDisabled(val bool) {
	e.errorDisabled = val
}

func (e *Adapter) IsErrorDisabled() bool {
	return e.errorDisabled
}

func (e *Adapter) EnableDebug(d bool) {
	e.debug = d
}

func (e *Adapter) IsDebugEnabled() bool {
	return e.debug
}

func (e *Adapter) SetMaxScreenshotLength(len int) {
	e.maxScreenshotLength = len
}

func (e *Adapter) SetImageSize(w int, h int) {
	if w > 0 && h > 0 {
		e.imgWidth = w
		e.imgHeight = h
	}
}

func (e *Adapter) Close() error {
	return e.driver.Close()
}

func (e *Adapter) Navigate(url string) error {
	return e.driver.Navigate(url)
}

func (e *Adapter) SwitchParentFrame() error {
	return e.tester(e.wait, func() error { return e.driver.SwitchParentFrame() })
}

func (e *Adapter) SelectWindowMain() error {
	return e.tester(e.wait, func() error { return e.driver.SwitchWindow(e.rootWindow) })
}

func (e *Adapter) CurrentWindowHandle() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.CurrentWindowHandle() })
}

func (e *Adapter) WindowHandles() ([]string, error) {
	return e.testerStringArray(e.wait, func() ([]string, error) { return e.driver.WindowHandles() })
}

func (e *Adapter) SwitchWindow(window string) error {
	return e.tester(e.wait, func() error { return e.driver.SwitchWindow(window) })
}

func (e *Adapter) SwitchFrame(frame interface{}) error {
	return e.tester(e.wait, func() error { return e.driver.SwitchFrame(frame) })
}

func (e *Adapter) Title() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.Title() })
}

func (e *Adapter) ExecuteScript(script string, args []interface{}) error {
	return e.tester(e.wait, func() error { _, err := e.driver.ExecuteScript(script, args); return err })
}

func (e *Adapter) ActiveElement() (base.IWebElement, error) {
	return e.testerElement(e.wait, func() (base.IWebElement, error) { return e.driver.ActiveElement() })
}

func (e *Adapter) Click(buttonId int) error {
	return e.tester(e.wait, func() error { return e.driver.Click(buttonId) })
}

func (e *Adapter) DoubleClick() error {
	return e.tester(e.wait, func() error { return e.driver.DoubleClick() })
}

func (e *Adapter) ButtonUp() error {
	return e.tester(e.wait, func() error { return e.driver.ButtonUp() })
}

func (e *Adapter) ButtonDown() error {
	return e.tester(e.wait, func() error { return e.driver.ButtonDown() })
}

func (e *Adapter) AlertText() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.AlertText() })
}

func (e *Adapter) AcceptAlert() error {
	return e.tester(e.wait, func() error { return e.driver.AcceptAlert() })
}

func (e *Adapter) DismissAlert() error {
	return e.tester(e.wait, func() error { return e.driver.DismissAlert() })
}

func (e *Adapter) CloseWindow(handle string) error {
	return e.tester(e.wait, func() error { return e.driver.CloseWindow(handle) })
}

func (e *Adapter) PageSource() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.PageSource() })
}

func (e *Adapter) Status() (string, error) {
	return e.testerString(e.wait, func() (string, error) {
		status, err := e.driver.Status()
		if err != nil {
			return "", err
		}
		src, _ := json.Marshal(status)
		return string(src), nil
	})
}

func (e *Adapter) DeleteCookie(id string) error {
	return e.tester(e.wait, func() error { return e.driver.DeleteCookie(id) })
}

func (e *Adapter) DeleteAllCookies() error {
	return e.tester(e.wait, func() error { return e.driver.DeleteAllCookies() })
}

func (e *Adapter) GetCookie(mode string, name string) (string, error) {
	return e.testerString(e.wait, func() (string, error) {
		cookie, err := e.driver.GetCookie(name)
		if err != nil {
			return "", err
		}
		if mode == "curl" {
			return base.Cookies2Curl([]base.Cookie{cookie})
		}
		return base.Cookies2Json([]base.Cookie{cookie})
	})
}

func (e *Adapter) GetCookies(mode string) (string, error) {
	return e.testerString(e.wait, func() (string, error) {
		cookies, err := e.driver.GetCookies()
		if err != nil {
			return "", err
		}
		if mode == "curl" {
			return base.Cookies2Curl(cookies)
		}
		return base.Cookies2Json(cookies)
	})
}

func (e *Adapter) KeyDown(keys string) error {
	return e.tester(e.wait, func() error { return e.driver.KeyDown(keys) })
}

func (e *Adapter) ResizeWindow(handle string, w int, h int) error {
	return e.tester(e.wait, func() error { return e.driver.ResizeWindow(handle, w, h) })
}

func (e *Adapter) ElementSendKeys(elm base.IWebElement, keys string) error {
	return e.tester(e.wait, func() error { return elm.SendKeys(keys) })
}

func (e *Adapter) ElementClick(elm base.IWebElement) error {
	return e.tester(e.wait, func() error { return elm.Click() })
}

func (e *Adapter) ElementMoveTo(elm base.IWebElement, xOffset float64, yOffset float64) error {
	return e.tester(e.wait, func() error { return elm.MoveTo(xOffset, yOffset) })
}

func (e *Adapter) _findElementReady(by string, data string, until int) (base.IWebElement, bool) {
	found := false
	elem, _ := e.driver.FindElement(by, data)
	if until == 0 {
		found = elem != nil
	} else if until == 1 {
		if elem != nil {
			if !e.isElementReady(by, data, elem) {
				found = true
				elem = nil
			}
		} else {
			found = true
		}
	} else {
		if elem != nil {
			if e.isElementReady(by, data, elem) {
				found = true
			}
		}
	}
	return elem, found
}

func (e *Adapter) FindElementReady(by string, data string, until int) base.IWebElement {
	var elem base.IWebElement = nil
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(e.wait) {
		var found bool
		elem, found = e._findElementReady(by, data, until)
		if found {
			break
		}
		e.RetryWait()
	}
	return elem
}

func (e *Adapter) SendKeys(by string, data string, value string) error {
	if len(value) == 0 {
		return nil
	}
	var err error = nil
	var elm base.IWebElement = nil
	pos := 0
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(e.wait) {
		if elm == nil || err != nil {
			if elm, _ = e._findElementReady(by, data, 2); elm == nil {
				err = fmt.Errorf("element isn't ready")
			} else {
				err = nil
			}
		}
		if elm != nil && err == nil {
			err = elm.SendKeys(string(value[pos]))
		}
		e.HumanWait()
		if err == nil {
			pos++
			if pos >= len(value) {
				break
			}
		} else {
			e.RetryWait()
		}
	}
	return nil
}

func (e *Adapter) MouseActions(by string, data string, command string, value string) error {
	var err error = nil
	var elm base.IWebElement = nil
	var start = UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(e.wait) {
		if elm == nil || err != nil {
			if elm, _ = e._findElementReady(by, data, 2); elm == nil {
				err = fmt.Errorf("element isn't ready")
			} else {
				err = nil
			}
		}

		if elm != nil && err == nil {
			switch command {
			case "click":
				err = elm.Click()
			case "doubleClick":
				if err = elm.MoveTo(0, 0); err == nil {
					err = e.driver.DoubleClick()
				}
			case "rightClick":
				if err = elm.MoveTo(0, 0); err == nil {
					err = e.driver.Click(2)
				}
			case "mouseOver":
				err = elm.MoveTo(0, 0)
			case "mouseUpAt":
				if e.drag != nil {
					var coordsUpAt = e.getCoords(value)
					if err = e.drag.MoveTo(coordsUpAt.X, coordsUpAt.Y); err == nil {
						err = e.driver.ButtonUp()
					}
					e.drag = nil
				}
			case "mouseDownAt":
				var coordsDownAt = e.getCoords(value)
				if err = elm.MoveTo(coordsDownAt.X, coordsDownAt.Y); err == nil {
					err = e.driver.ButtonDown()
					e.drag = elm
				}
			case "mouseMultipleMoveAt":
				if e.drag != nil {
					//TODO l'iterazione deve essere nel ciclo principale
					baseCords := base.Point{X: 0, Y: 0}
					for _, m := range strings.Split(value, "|") {
						var coordsMultipleMoveAt = e.getCoords(m)
						x := baseCords.X + coordsMultipleMoveAt.X
						y := baseCords.Y + coordsMultipleMoveAt.Y
						baseCords.X = x
						baseCords.Y = y
						if err = e.drag.MoveTo(x, y); err != nil {
							break
						}
						e.HumanWait()
					}
				}
			case "mouseMoveAt":
				if e.drag != nil {
					var coordsMoveAt = e.getCoords(value)
					err = e.drag.MoveTo(coordsMoveAt.X, coordsMoveAt.Y)
				}
			}
		}
		if err == nil {
			break
		}
		e.RetryWait()
	}
	return err
}

func (e *Adapter) NetworkHeaders() (string, error) {
	logEntries, logErr := e.driver.Log(base.LogPerformance)
	if logErr != nil {
		e.log(LogLevelCritical, logErr.Error())
	}
	n := e.network.Headers(logEntries)
	out, _ := json.MarshalIndent(n, "", "  ")
	return string(out), nil
}

func (e *Adapter) Screenshot(screenshotId string) error {
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

func (e *Adapter) HumanWait() {
	rnd := rand.Float64()
	val := math.Round(rnd * 100)
	interval := e.humanWaitBase + int(val)
	e.Sleep(interval)
}

func (e *Adapter) RetryWait() {
	e.Sleep(e.retryInterval)
}

func (e *Adapter) Sleep(interval int) {
	time.Sleep(time.Millisecond * time.Duration(interval))
}

func (e *Adapter) CreateEvent(id string, err error, kind string, start time.Time, shot bool) *Event {
	if e.errorDisabled {
		err = nil
	}
	event := NewEvent(e, id, kind, err, start)
	event.ProbeId = e.probeId
	event.UUID = e.uuid
	logEntries, logErr := e.driver.Log(base.LogPerformance) //e.Log(base.LogPerformance)
	if logErr != nil {
		e.log(LogLevelCritical, logErr.Error())
	}
	event.Network, event.NetworkErrorCount = e.network.Compute(logEntries)
	if acquired := e.network.Acquired(); acquired != nil {
		event.Acquired = acquired
	}
	if shot && err != nil {
		event.ScreenShoot = "probes-" + uuid.New().String()
		if screenErr := e.Screenshot(event.ScreenShoot); screenErr != nil {
			e.log(LogLevelCritical, "error generating screenshot: "+screenErr.Error())
		}
	}
	return event
}

func (e *Adapter) WriteLog(message string) {
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

func (e *Adapter) log(kind Loglevel, src string) {
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
	e.WriteLog(string(message))
}

func (e *Adapter) isElementReady(by string, data string, elm base.IWebElement) bool {
	ok, err := elm.IsEnabled()
	if err != nil {
		e.log(LogLevelDebug, fmt.Sprintf("(%s=%s) isElementReady : %s", by, data, err.Error()))
		return false
	}
	if !ok {
		e.log(LogLevelDebug, fmt.Sprintf("(%s=%s) isElementReady: element isn't enabled", by, data))
		return false
	}
	ok, err = elm.IsDisplayed()
	if err != nil {
		e.log(LogLevelDebug, fmt.Sprintf("(%s=%s) isElementReady (IsDisplayed): %s", by, data, err.Error()))
		return false
	}
	if !ok {
		e.log(LogLevelDebug, fmt.Sprintf("(%s=%s) isElementReady: element isn't displayed", by, data))
		return false
	}
	return true
}

func (e *Adapter) getCoords(data string) base.Point {
	var coords base.Point
	if container := strings.Split(data, ","); len(container) > 1 {
		coords.X, _ = strconv.ParseFloat(container[0], 64)
		coords.Y, _ = strconv.ParseFloat(container[1], 64)
	}
	return coords
}

func (e *Adapter) tester(wait int, fn func() error) error {
	var err error
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(wait) {
		if err = fn(); err == nil {
			return nil
		}
		e.RetryWait()
	}
	return err
}

func (e *Adapter) testerString(wait int, fn func() (string, error)) (string, error) {
	var val string
	var err error
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(wait) {
		if val, err = fn(); err == nil {
			return val, nil
		}
		e.RetryWait()
	}
	return "", err
}

func (e *Adapter) testerStringArray(wait int, fn func() ([]string, error)) ([]string, error) {
	var val []string
	var err error
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(wait) {
		if val, err = fn(); err == nil {
			return val, nil
		}
		e.RetryWait()
	}
	return nil, err
}

func (e *Adapter) testerElement(wait int, fn func() (base.IWebElement, error)) (base.IWebElement, error) {
	var val base.IWebElement
	var err error
	start := UnixMilli(time.Now())
	for UnixMilli(time.Now())-start < int64(wait) {
		if val, err = fn(); err == nil {
			return val, nil
		}
		e.RetryWait()
	}
	return nil, err
}

func (e *Adapter) Scroll(x int, y int, deltaX int, deltaY int) error {
	e.driver.StoreWheelActions("wheel1", base.CreateWheelAction(x, y, deltaX, deltaY))
	if err := e.driver.PerformActions(); err != nil {
		return err
	}
	return nil
}
