package executor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/markel1974/webautoma/src/wd/base"
)

// RFC3339Milli defines a timestamp format with millisecond precision conforming to the RFC3339 specification.
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

// DefaultResultFile specifies the default filename for storing log results in JSON format.
// DefaultImagesFile specifies the default filename for storing image information in JSON format.
const (
	DefaultResultFile = "log.json"
	DefaultImagesFile = "images.json"
)

// Loglevel represents the severity level of a log message.
// It can be used to distinguish between debug, info, warning, and critical log messages.
type Loglevel int

// LogLevelDebug represents debug-level logging, used for detailed and diagnostic information.
// LogLevelInfo represents informational logging, typically for general application operations.
// LogLevelWarning represents warning-level logging, indicating potential issues requiring attention.
// LogLevelCritical represents critical-level logging, used for severe error messages likely requiring immediate action.
const (
	LogLevelDebug    Loglevel = -1
	LogLevelInfo     Loglevel = 0
	LogLevelWarning  Loglevel = 1
	LogLevelCritical Loglevel = 2
)

// Adapter represents a wrapper around a WebDriver to manage browser interactions and configurations.
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

// NewAdapter initializes and returns a new Adapter instance based on the provided IWebDriver implementation.
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

// Setup initializes the Adapter instance with configuration for logging, image handling, and network profiling.
// It sets up the log file, image file, image dump flag, and creates a new Network instance with specified profiles.
// Returns an error if network initialization fails.
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

// GetWindowHandle retrieves the window handle for the given identifier and returns it as a string along with any error encountered.
func (e *Adapter) GetWindowHandle(id string) (string, error) {
	return e.windowHandles.GetHandle(id)
}

// StoreWindowHandle stores the current window handle identified by the provided id into the windowHandles container of the Adapter.
func (e *Adapter) StoreWindowHandle(id string) error {
	return e.windowHandles.Store(id)
}

// AddWindowHandle adds a new window handle to the internal window collection using the provided ConfigCommand.
// It returns the timeout value associated with the created window handle.
func (e *Adapter) AddWindowHandle(cmd ConfigCommand) int {
	return e.windowHandles.Add(e.driver, cmd)
}

// Quit releases the underlying driver resources and terminates the adapter's operation.
func (e *Adapter) Quit() error {
	return e.driver.Quit()
}

// GetRootWindow retrieves the identifier for the root window of the adapter instance.
func (e *Adapter) GetRootWindow() string {
	return e.rootWindow
}

// SetRootWindow sets the current window handle as the root window for the adapter. Returns an error if retrieval fails.
func (e *Adapter) SetRootWindow() error {
	handle, err := e.CurrentWindowHandle()
	if err != nil {
		return err
	}
	e.rootWindow = handle
	return err
}

// SetProbeId assigns the provided probe ID to the adapter's internal probeId field.
func (e *Adapter) SetProbeId(id string) {
	e.probeId = id
}

// SetRetryInterval sets the interval in milliseconds for retrying failed operations within the adapter.
func (e *Adapter) SetRetryInterval(interval int) {
	e.retryInterval = interval
}

// SetWait sets the wait duration for the adapter in milliseconds.
func (e *Adapter) SetWait(wait int) {
	e.wait = wait
}

// SetHumanWait sets the base wait time in milliseconds for human-like delays during execution.
func (e *Adapter) SetHumanWait(wait int) {
	e.humanWaitBase = wait
}

// SetErrorDisabled sets the errorDisabled flag to the provided boolean value, enabling or disabling error handling.
func (e *Adapter) SetErrorDisabled(val bool) {
	e.errorDisabled = val
}

// IsErrorDisabled determines if error handling is currently disabled in the adapter configuration.
func (e *Adapter) IsErrorDisabled() bool {
	return e.errorDisabled
}

// EnableDebug toggles debug mode for the adapter based on the provided boolean parameter.
func (e *Adapter) EnableDebug(d bool) {
	e.debug = d
}

// IsDebugEnabled checks if the debug mode is enabled for the adapter and returns true if enabled, otherwise false.
func (e *Adapter) IsDebugEnabled() bool {
	return e.debug
}

// SetMaxScreenshotLength sets the maximum allowed length for a screenshot in the adapter configuration.
func (e *Adapter) SetMaxScreenshotLength(len int) {
	e.maxScreenshotLength = len
}

// SetImageSize sets the image dimensions for width and height if both values are greater than zero.
func (e *Adapter) SetImageSize(w int, h int) {
	if w > 0 && h > 0 {
		e.imgWidth = w
		e.imgHeight = h
	}
}

// Close terminates the underlying driver's connection and releases associated resources. Returns an error if it fails.
func (e *Adapter) Close() error {
	return e.driver.Close()
}

// Navigate loads the specified URL in the web driver and returns an error if navigation fails.
func (e *Adapter) Navigate(url string) error {
	return e.driver.Navigate(url)
}

// SwitchParentFrame switches the context to the parent frame of the current frame in the web driver session.
func (e *Adapter) SwitchParentFrame() error {
	return e.tester(e.wait, func() error { return e.driver.SwitchParentFrame() })
}

// SelectWindowMain switches the focus of the driver to the main root window of the application.
func (e *Adapter) SelectWindowMain() error {
	return e.tester(e.wait, func() error { return e.driver.SwitchWindow(e.rootWindow) })
}

// CurrentWindowHandle retrieves the handle of the current browser window and returns it as a string.
func (e *Adapter) CurrentWindowHandle() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.CurrentWindowHandle() })
}

// WindowHandles retrieves a list of handles for all open browser windows or tabs associated with the current session.
func (e *Adapter) WindowHandles() ([]string, error) {
	return e.testerStringArray(e.wait, func() ([]string, error) { return e.driver.WindowHandles() })
}

// SwitchWindow switches the active window to the specified window using its handle and waits for the operation to complete.
func (e *Adapter) SwitchWindow(window string) error {
	return e.tester(e.wait, func() error { return e.driver.SwitchWindow(window) })
}

// SwitchFrame switches the browser's focus to the specified frame.
// The frame parameter can represent the frame index, name, or a relative frame reference.
// Returns an error if the frame cannot be switched.
func (e *Adapter) SwitchFrame(frame interface{}) error {
	return e.tester(e.wait, func() error { return e.driver.SwitchFrame(frame) })
}

// Title retrieves the title of the current window in the web driver and returns it as a string along with any error encountered.
func (e *Adapter) Title() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.Title() })
}

// ExecuteScript executes a script on the underlying driver with the provided arguments and returns any resulting error.
func (e *Adapter) ExecuteScript(script string, args []interface{}) error {
	return e.tester(e.wait, func() error { _, err := e.driver.ExecuteScript(script, args); return err })
}

// ActiveElement returns the currently focused element on the page as an IWebElement or an error if the operation fails.
func (e *Adapter) ActiveElement() (base.IWebElement, error) {
	return e.testerElement(e.wait, func() (base.IWebElement, error) { return e.driver.ActiveElement() })
}

// Click sends a click event to the UI element identified by the provided buttonId and returns an error if the operation fails.
func (e *Adapter) Click(buttonId int) error {
	return e.tester(e.wait, func() error { return e.driver.Click(buttonId) })
}

// DoubleClick performs a double-click action using the underlying driver and tester mechanisms, returning an error if it fails.
func (e *Adapter) DoubleClick() error {
	return e.tester(e.wait, func() error { return e.driver.DoubleClick() })
}

// ButtonUp triggers the action for releasing a button, utilizing the tester and driver implementations.
func (e *Adapter) ButtonUp() error {
	return e.tester(e.wait, func() error { return e.driver.ButtonUp() })
}

// ButtonDown simulates pressing and holding down a mouse button via the underlying driver. Returns an error if the operation fails.
func (e *Adapter) ButtonDown() error {
	return e.tester(e.wait, func() error { return e.driver.ButtonDown() })
}

// AlertText retrieves the text from an active browser alert dialog and returns it as a string.
func (e *Adapter) AlertText() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.AlertText() })
}

// AcceptAlert confirms and accepts the currently active alert dialog, if present, within the web driver session.
func (e *Adapter) AcceptAlert() error {
	return e.tester(e.wait, func() error { return e.driver.AcceptAlert() })
}

// DismissAlert dismisses the active alert dialog and returns an error if the operation fails.
func (e *Adapter) DismissAlert() error {
	return e.tester(e.wait, func() error { return e.driver.DismissAlert() })
}

// CloseWindow closes the window identified by the given handle string and returns an error if the operation fails.
func (e *Adapter) CloseWindow(handle string) error {
	return e.tester(e.wait, func() error { return e.driver.CloseWindow(handle) })
}

// PageSource retrieves the current page's source as a string using the underlying driver, with a wait condition applied.
func (e *Adapter) PageSource() (string, error) {
	return e.testerString(e.wait, func() (string, error) { return e.driver.PageSource() })
}

// Status retrieves the current status of the driver and returns it as a JSON-encoded string, or an error if retrieval fails.
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

// DeleteCookie removes a cookie identified by the given ID from the current browser session.
func (e *Adapter) DeleteCookie(id string) error {
	return e.tester(e.wait, func() error { return e.driver.DeleteCookie(id) })
}

// DeleteAllCookies removes all cookies associated with the current session. It returns an error if the operation fails.
func (e *Adapter) DeleteAllCookies() error {
	return e.tester(e.wait, func() error { return e.driver.DeleteAllCookies() })
}

// GetCookie retrieves a cookie by its name and returns it as a string.
// The output format is determined by the mode: "curl" for curl format or default for JSON.
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

// GetCookies retrieves all browser cookies and formats them either as JSON or Curl command, based on the given mode.
// The mode parameter specifies the desired output format: "curl" for Curl syntax or any other string for JSON.
// Returns the formatted string of cookies or an error in case of failure.
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

// GetCookiesHttp retrieves all cookies from the underlying driver and converts them to http.Cookie format.
// Returns a slice of http.Cookie pointers and an error if the operation fails.
func (e *Adapter) GetCookiesHttp() ([]*http.Cookie, error) {
	cookies, err := e.driver.GetCookies()
	if err != nil {
		return nil, err
	}
	var out []*http.Cookie
	for _, c := range cookies {
		out = append(out, c.Http())
	}
	return out, nil
}

// KeyDown simulates pressing and holding down the specified keys in the context of the associated driver.
func (e *Adapter) KeyDown(keys string) error {
	return e.tester(e.wait, func() error { return e.driver.KeyDown(keys) })
}

// ResizeWindow resizes the specified window to the given width and height using the provided handle.
// Returns an error if the operation fails.
func (e *Adapter) ResizeWindow(handle string, w int, h int) error {
	return e.tester(e.wait, func() error { return e.driver.ResizeWindow(handle, w, h) })
}

// ElementSendKeys sends the specified string of keys to the given web element, using a defined tester function for execution.
func (e *Adapter) ElementSendKeys(elm base.IWebElement, keys string) error {
	return e.tester(e.wait, func() error { return elm.SendKeys(keys) })
}

// ElementClick triggers a click action on the provided web element, waiting for it to be ready before execution.
func (e *Adapter) ElementClick(elm base.IWebElement) error {
	return e.tester(e.wait, func() error { return elm.Click() })
}

// ElementMoveTo moves the specified web element by the given x and y offsets.
// elm represents the web element, xOffset and yOffset specify the movement distances in the respective directions.
// Returns an error if the operation fails.
func (e *Adapter) ElementMoveTo(elm base.IWebElement, xOffset float64, yOffset float64) error {
	return e.tester(e.wait, func() error { return elm.MoveTo(xOffset, yOffset) })
}

// _findElementReady attempts to locate a web element and checks its readiness based on specified conditions.
// The `by` parameter specifies the strategy to find the element (e.g., ID, class, XPath).
// The `data` parameter indicates the search criteria for the element.
// The `until` parameter specifies the readiness condition (0 for existence, 1 for non-readiness, 2 for readiness).
// Returns the found element and a boolean indicating if the expected condition was met.
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

// FindElementReady searches for a web element matching the given criteria and waits until it is ready or timeout is reached.
// The 'by' parameter specifies the locating strategy (e.g., ID, XPath).
// The 'data' parameter contains the locator value used for finding the element.
// The 'until' parameter defines the condition to be met for the element to be considered ready.
// Returns the located web element if found and ready; otherwise, nil when the maximum wait time is exceeded.
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

// FindAttribute retrieves the value of a specified attribute from a web element identified by selector and data.
// It retries until the element is ready or the timeout period expires.
// Returns the attribute value or an error if the element or attribute is unavailable.
func (e *Adapter) FindAttribute(by string, data string, attribute string) (string, error) {
	var err error = nil
	var val string
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
		if elm != nil {
			val, err = elm.GetAttribute(attribute)
		}
		if err == nil {
			break
		}
		e.RetryWait()
	}
	return val, err
}

func (e *Adapter) SendKeysActions(values []string) error {
	var keys []base.KeyAction
	for _, v := range values {
		m := strings.Split(v, ":")
		if len(m) >= 2 {
			key := strings.TrimSpace(strings.ToLower(m[0]))
			val := m[1]
			switch key {
			case "press":
				keys = append(keys, base.KeyDownAction(base.KeyFromMapping(val)))
			case "release":
				keys = append(keys, base.KeyUpAction(base.KeyFromMapping(val)))
			case "pause":
				w := uint(100)
				if a, err := strconv.Atoi(val); err == nil {
					w = uint(a)
				}
				keys = append(keys, base.KeyPauseAction(w))
			}
		} else {
			keys = append(keys, base.KeyDownAction(v))
		}
	}
	if len(keys) == 0 {
		return nil
	}
	e.driver.StoreKeyActions("keyboard "+uuid.New().String(), keys...)
	if err := e.driver.PerformActions(); err != nil {
		return err
	}
	if err := e.driver.ReleaseActions(); err != nil {
		return err
	}
	return nil
}

// SendKeys sends a sequence of characters to the targeted web element located by the specified selector.
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

// MouseActions performs various mouse interactions (e.g., click, drag, mouse move) on a web element based on the passed command.
// Parameters:
// - by: The strategy to locate the element (e.g., CSS selector, XPath).
// - data: The selector or value to locate the target element.
// - command: The mouse action to execute (e.g., "click", "doubleClick").
// - value: Additional data for certain actions (e.g., coordinates for drag or move actions).
// Returns an error if the operation fails or the element is not ready.
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

// NetworkHeaders retrieves network headers from performance logs and returns them as a formatted JSON string.
// It collects log entries, processes them through the network's Headers method, and logs critical issues if encountered.
func (e *Adapter) NetworkHeaders() (string, error) {
	logEntries, logErr := e.driver.Log(base.LogPerformance)
	if logErr != nil {
		e.log(LogLevelCritical, logErr.Error())
	}
	n := e.network.Headers(logEntries)
	out, _ := json.MarshalIndent(n, "", "  ")
	return string(out), nil
}

// Screenshot captures a screenshot using the provided driver and processes it for scaling, encoding, and storing as needed.
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

// HumanWait simulates a human-like delay by introducing a randomized sleep interval based on a base value.
func (e *Adapter) HumanWait() {
	rnd := rand.Float64()
	val := math.Round(rnd * 100)
	interval := e.humanWaitBase + int(val)
	e.Sleep(interval)
}

// RetryWait pauses execution for the duration of the retry interval defined in the adapter.
func (e *Adapter) RetryWait() {
	e.Sleep(e.retryInterval)
}

// Sleep pauses the current execution for the specified interval in milliseconds.
func (e *Adapter) Sleep(interval int) {
	time.Sleep(time.Millisecond * time.Duration(interval))
}

// CreateEvent initializes and returns a new Event instance with the provided parameters.
// It computes network data, optionally triggers a screenshot, and logs errors if encountered.
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

// WriteLog appends a given message to the log file specified in the adapter, creating the file if it doesn't exist.
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

// log writes a log entry with a specified severity level and message.
// It discards debug logs if debug mode is disabled.
// Severity levels include debug, info, warning, and critical.
// Logs are serialized to JSON format and written via WriteLog.
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

// isElementReady checks if a given web element is both enabled and displayed.
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

// getCoords parses a string containing coordinates in "X,Y" format and returns a base.Point with X and Y values as floats.
func (e *Adapter) getCoords(data string) base.Point {
	var coords base.Point
	if container := strings.Split(data, ","); len(container) > 1 {
		coords.X, _ = strconv.ParseFloat(container[0], 64)
		coords.Y, _ = strconv.ParseFloat(container[1], 64)
	}
	return coords
}

// tester retries the provided function `fn` until it succeeds or the specified `wait` time (in milliseconds) elapses.
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

// testerString retries a function until successful or timeout occurs; waits `wait` ms and calls `RetryWait` between attempts.
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

// testerStringArray repeatedly executes the provided function within a specified time limit until it succeeds or time expires.
// The wait parameter specifies the timeout duration in milliseconds.
// The fn parameter is a function returning a slice of strings and an error, representing the operation to be retried.
// Returns the slice of strings from the successful execution of fn or an error if all retries fail.
// Utilizes RetryWait for any wait between retries when retries are required.
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

// testerElement attempts to retrieve a web element through retries within a specified timeout duration (wait in milliseconds).
// It uses the provided callback function fn to attempt the fetch and returns the element or an error if unsuccessful.
// RetryWait is invoked between attempts to manage retry intervals.
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

// Scroll performs a scroll action at the specified coordinates with the given deltas.
func (e *Adapter) Scroll(x int, y int, deltaX int, deltaY int) error {
	e.driver.StoreWheelActions("wheel1", base.CreateWheelAction(x, y, deltaX, deltaY))
	if err := e.driver.PerformActions(); err != nil {
		return err
	}
	return nil
}
