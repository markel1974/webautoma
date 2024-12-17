package executor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/markel1974/webautoma/src/wd/base"
)

// RequiredLogs returns the log type and log level required for performance monitoring and debugging purposes.
func RequiredLogs() (base.LogType, base.LogLevel) {
	return base.LogPerformance, base.LogAll
}

// eventMessageOK represents a message indicating the success of a probe.
// eventMessageNOK represents a message indicating an error occurred in the probe.
// eventMessageTimerNotFinalized indicates that the timer operation was not finalized.
const (
	eventMessageOK                = "Esito sonda ok"
	eventMessageNOK               = "Errore nella sonda"
	eventMessageTimerNotFinalized = "timer not finalized"
)

// Event kind constants represent different types of events: full or intermediate.
const (
	eventKindFull         = "full"
	eventKindIntermediate = "intermediate"
)

// Executor is a type that manages execution context and orchestrates components like templates, timers, and file downloads.
// It handles configuration, state management, and interactions with the SamConsole and Downloader.
type Executor struct {
	cfgFile      string
	start        time.Time
	quit         bool
	maxRetry     int
	timers       map[string]*Timer
	templates    *Templates
	cfg          Config
	stack        map[string]interface{}
	adapter      *Adapter
	execId       string
	lastFrame    interface{}
	baseUrl      string
	sam          *SamConsole
	downloadPath string
	downloader   *Downloader
}

// New initializes and returns a new Executor instance configured with the provided IWebDriver.
func New(driver base.IWebDriver) *Executor {
	e := &Executor{
		templates:  nil,
		quit:       false,
		maxRetry:   3,
		timers:     make(map[string]*Timer),
		stack:      nil,
		adapter:    NewAdapter(driver),
		execId:     "",
		lastFrame:  nil,
		baseUrl:    "",
		sam:        nil,
		downloader: NewDownloader(&http.Client{}),
	}
	downloadPath, err := os.UserHomeDir()
	if err == nil {
		downloadPath = filepath.Join(downloadPath, "Downloads") + string(os.PathSeparator)
		e.downloadPath = downloadPath
	}
	return e
}

// Setup initializes the executor with configuration, templates, and adapter settings based on provided parameters.
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
	if len(e.cfg.DownloadPath) > 0 {
		e.downloadPath = e.cfg.DownloadPath
	}
	if e.cfg.Debug {
		e.adapter.EnableDebug(true)
	}
	if err = e.adapter.Setup(logFile, imgFile, imgDump, profileCapture, e.cfg.ProfileSupportedMethods); err != nil {
		return err
	}
	return nil
}

// SetDownloadPath sets the file download path for the Executor to the specified directory.
func (e *Executor) SetDownloadPath(downloadPath string) {
	e.downloadPath = downloadPath
}

// SetImageSize sets the width and height of the image using the provided dimensions (w for width, h for height).
func (e *Executor) SetImageSize(w int, h int) {
	e.adapter.SetImageSize(w, h)
}

// loadUrl navigates to the specified URL using the adapter and updates the base URL in the Executor. It returns any error.
func (e *Executor) loadUrl(url string) error {
	if err := e.adapter.Navigate(url); err != nil {
		return err
	}
	e.baseUrl = url
	return nil
}

// computeScroll parses the scroll value, extracting scroll coordinates, step count, and interval from the input string.
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

// computeSelector parses the target, determines the selection strategy, and returns the strategy, value, or an error.
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

// doScreenshot captures a screenshot using the provided screenshot ID and delegates the operation to the adapter.
func (e *Executor) doScreenshot(screenshotId string) error {
	return e.adapter.Screenshot(screenshotId)
}

// doTimer handles the creation, starting, stopping, and finalization of timers based on the provided action parameters.
// target specifies the timer identifier.
// until determines the operation to perform: "timerCreate", "timerStart", "timerStop", or "timerFinalize".
// value provides additional context or settings for the operation, such as finalization instructions.
// Returns an error if the timer operation fails or if the provided target is unrecognized.
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

// doMouse performs mouse-related actions on a target element using the specified command and value.
// It computes the element's selector and applies the necessary mouse action through the adapter.
// Returns an error if the selector computation or mouse action execution fails.
func (e *Executor) doMouse(target string, command string, value string) error {
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	if err = e.adapter.MouseActions(by, data, command, value); err != nil {
		return err
	}
	return nil
}

// doClickDownload downloads a file by finding the "href" attribute of an element and performing an HTTP request.
// target specifies the selector of the element; value specifies the file name for saving the download.
// Returns an error if the operation fails at any stage: selector computation, attribute retrieval, or download.
func (e *Executor) doClickDownload(target string, value string) error {
	by, data, err := e.computeSelector(target)
	if err != nil {
		return err
	}
	val, err := e.adapter.FindAttribute(by, data, "href")
	if err != nil {
		return err
	}
	cookies, err := e.adapter.GetCookiesHttp()
	if err != nil {
		return err
	}
	if err = e.downloader.Do(val, cookies, value); err != nil {
		return err
	}
	return nil
}

// doType types a sequence of characters into a target element or sends individual keystrokes, with optional delays.
// Returns an error if the operation fails due to an invalid target, selector computation issue, or keystroke failure.
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

// doRunScript executes the given script with specified arguments using the adapter and returns any resulting errors.
func (e *Executor) doRunScript(script string, args []interface{}) error {
	return e.adapter.ExecuteScript(script, args)
}

// doSelectFrame switches the execution context to a specified frame based on the given target string.
// The target string can specify the frame by index, relative name, or other identifiers.
// Returns an error if the target cannot be parsed or the frame switch fails.
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

// doSelectParentFrame switches to the parent frame of the current frame and resets the lastFrame field to nil.
func (e *Executor) doSelectParentFrame() error {
	e.lastFrame = nil
	return e.adapter.SwitchParentFrame()
}

// doSelectWindowMain switches the context to the main application window and interacts with the adapter to perform the operation.
func (e *Executor) doSelectWindowMain() error {
	return e.adapter.SelectWindowMain()
}

// doListWindows retrieves the titles of all browser windows and returns them as a slice of strings.
func (e *Executor) doListWindows() ([]string, error) {
	var out []string
	currentWindow, err := e.adapter.CurrentWindowHandle()
	if err != nil {
		return nil, err
	}
	windows, err := e.adapter.WindowHandles()
	if err != nil {
		return nil, err
	}
	for _, window := range windows {
		if err = e.adapter.SwitchWindow(window); err != nil {
			continue
		}
		windowTitle, err := e.adapter.Title()
		if err != nil {
			continue
		}
		out = append(out, windowTitle)
	}
	_ = e.adapter.SwitchWindow(currentWindow)
	return out, nil
}

// doSelectWindowTitle switches to a browser window by its title. Returns an error if no such window is found.
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

// doSelectAlert processes an alert dialog based on the input value: accepts, dismisses, or performs no action.
// Returns an error if the alert interaction or parsing fails.
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

// doCloseWindow closes the current browser window unless it is the root window, switching back to the root window afterward.
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

// doClose closes a window identified by the target in the provided ConfigCommand. Returns error if the handle is invalid.
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

// doStackAdd adds an element's properties to the execution stack under the specified id if the element is found.
// id specifies the key to store the element data in the stack.
// target is the selector to locate the element.
// u is a conditional string that determines wait behavior. If non-empty, readiness is awaited before detection.
// Returns an error if the element cannot be located or any other issue arises during processing.
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

// doStackPrint serializes the stack data into an indented JSON format and prints it to the standard output. Returns an error.
func (e *Executor) doStackPrint() error {
	k, _ := json.MarshalIndent(e.stack, "", "	")
	fmt.Println(string(k))
	return nil
}

// doStackReset clears the Executor's stack by setting it to nil and returns any error that occurs.
func (e *Executor) doStackReset() error {
	e.stack = nil
	return nil
}

// doAssert verifies if the text of the found element matches the given caption and returns an error if they differ or fail.
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

// doExists verifies if the specified element exists on the page based on the provided target and condition parameters.
// Returns an error if the element is not found or if there is an issue during the operation.
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

// doUntil executes a command until a specific condition is evaluated, using the provided ConfigCommand parameters.
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

// doSelectWindow switches to the window identified by the given target string, using the adapter's window handling API.
// Returns an error if the window cannot be found or if switching to the window fails.
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

// doStoreWindowHandle stores the current window handle using the provided target string from the command configuration.
func (e *Executor) doStoreWindowHandle(cmd ConfigCommand) error {
	target := cmd.Target
	return e.adapter.StoreWindowHandle(target)
}

// doSelect attempts to select an option from a dropdown element based on its label and provided target and value inputs.
// Returns an error if the input is invalid, no matching option is found, or interaction with the element fails.
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

// doActiveElement retrieves the currently active element using the adapter and prints its details with a given identifier.
func (e *Executor) doActiveElement(id string) error {
	elm, err := e.adapter.ActiveElement()
	if err != nil {
		return err
	}
	fmt.Printf("%s ActiveElement =>\n", id)
	elm.Print()
	return nil
}

// doWindowHandles retrieves and prints the current window handles, associating them with the provided id.
func (e *Executor) doWindowHandles(id string) error {
	fmt.Printf("%s WindowHandles =>\n", id)
	v, _ := e.adapter.WindowHandles()
	fmt.Println(v)
	return nil
}

// doPause pauses execution for a duration specified in the Target field of the given ConfigCommand.
// Returns an error if Target cannot be converted to an integer or if any other issue occurs during execution.
func (e *Executor) doPause(cmd ConfigCommand) error {
	target := cmd.Target
	v, err := strconv.Atoi(target)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * time.Duration(v))
	return nil
}

// doSetTimeout sets a timeout value by parsing the input string and updating the internal adapter wait configuration.
func (e *Executor) doSetTimeout(value string) error {
	e.adapter.SetWait(parseInt(value))
	return nil
}

// doRetrieveNetworkHeaders retrieves network headers via the adapter and returns them as a string along with any error.
func (e *Executor) doRetrieveNetworkHeaders() (string, error) {
	return e.adapter.NetworkHeaders()
}

// doRetrievePageSource retrieves the page source using the adapter and returns it along with any potential error.
func (e *Executor) doRetrievePageSource() (string, error) {
	v, err := e.adapter.PageSource()
	return v, err
}

// doPageSource retrieves the current page source via the adapter and prints it along with the provided identifier.
func (e *Executor) doPageSource(id string) error {
	v, _ := e.adapter.PageSource()
	fmt.Printf("%s PageSource =>\n", id)
	fmt.Println(v)
	return nil
}

// doStatus retrieves and prints the current status from the adapter, associating it with the provided command ID.
func (e *Executor) doStatus(cmd ConfigCommand) error {
	id := cmd.Id
	fmt.Printf("%s Status =>\n", id)
	k, _ := e.adapter.Status()
	fmt.Println(k)
	return nil
}

// doGetCookie retrieves a cookie value based on the target specifications and optionally saves it to a file.
// It uses the adapter's GetCookie method and can output the result to the console or a file based on the value parameter.
// Returns an error if there are issues in retrieving or saving the cookie.
func (e *Executor) doGetCookie(target string, id string, value string) error {
	mode := ""
	name := target
	k := strings.Split(name, ":")
	if len(k) > 1 {
		mode = k[0]
		name = k[1]
	}
	v, err := e.adapter.GetCookie(mode, name)
	if err != nil {
		return err
	}
	if len(value) > 0 {
		if err = os.WriteFile(value, []byte(v), 0644); err != nil {
			return err
		}
		return nil
	}
	fmt.Printf("%s Cookie =>\n", id)
	fmt.Println(v)
	return nil
}

// doGetAllCookies retrieves all cookies for a specified target and writes them to a file or outputs them to the console.
func (e *Executor) doGetAllCookies(target string, id string, value string) error {
	mode := target
	v, err := e.adapter.GetCookies(mode)
	if err != nil {
		return err
	}
	if len(value) > 0 {
		if err = os.WriteFile(value, []byte(v), 0644); err != nil {
			return err
		}
		return nil
	}
	fmt.Printf("%s Cookies =>\n", id)
	fmt.Println(v)
	return nil
}

// doDeleteAllCookies removes all cookies stored in the current browser session and returns an error if it fails.
func (e *Executor) doDeleteAllCookies() error {
	return e.adapter.DeleteAllCookies()
}

// doDeleteCookie removes a specific browser cookie identified by the target string.
// It returns an error if the operation fails.
func (e *Executor) doDeleteCookie(target string) error {
	return e.adapter.DeleteCookie(target)
}

// doActionsSendKeys simulates sending a sequence of key presses defined by the specified value using the adapter interface.
func (e *Executor) doActionsSendKeys(value string) error {
	//TODO TEST
	//for x := 0; x < 100; x++ {
	//	fmt.Println("SENDING KEYDOWN TEST")
	//	e.adapter.KeyDown(base.DownArrowKey)
	//	time.Sleep(100 * time.Millisecond)
	//}

	return e.adapter.KeyDown(value)
}

// doSetWindowMain sets the root window as the primary window for the current session using the adapter.
func (e *Executor) doSetWindowMain() error {
	return e.adapter.SetRootWindow()
}

// doHumanWait causes a delay by parsing the given string value as an integer and invoking the adapter's Sleep method.
func (e *Executor) doHumanWait(value string) error {
	e.adapter.Sleep(parseInt(value))
	return nil
}

// doEnableDebug enables the debug mode by configuring the adapter to allow detailed debugging information.
func (e *Executor) doEnableDebug() error {
	e.adapter.EnableDebug(true)
	return nil
}

// doDisableDebug disables the debug mode in the adapter by setting the debug flag to false.
func (e *Executor) doDisableDebug() error {
	e.adapter.EnableDebug(false)
	return nil
}

// doDisableError disables error handling for the associated adapter and returns any error encountered during the process.
func (e *Executor) doDisableError() error {
	e.adapter.SetErrorDisabled(true)
	return nil
}

// doEnableError enables error handling by setting the adapter's error state to not disabled. Returns an error if any issue occurs.
func (e *Executor) doEnableError() error {
	e.adapter.SetErrorDisabled(false)
	return nil
}

// doExecId sets the Executor's execId field to the target value from the provided ConfigCommand. Returns an error if any.
func (e *Executor) doExecId(cmd ConfigCommand) error {
	target := cmd.Target
	e.execId = target
	return nil
}

// doProbeId sets the probe ID for the provided target using the adapter in the Executor. Returns an error if the operation fails.
func (e *Executor) doProbeId(cmd ConfigCommand) error {
	target := cmd.Target
	e.adapter.SetProbeId(target)
	return nil
}

// doSetWindowSize adjusts the window size based on the target string input formatted as "widthxheight".
// Returns an error if the input format is invalid or resizing fails.
func (e *Executor) doSetWindowSize(target string) error {
	v := strings.Split(target, "x")
	if len(v) < 2 {
		return fmt.Errorf("invalid target")
	}
	width, _ := strconv.Atoi(v[0])
	height, _ := strconv.Atoi(v[1])
	return e.adapter.ResizeWindow("", width, height)
}

// doScroll performs a scroll action based on the provided command configuration, including offsets and element location.
// It calculates the scroll steps, coordinates, and interval, executing the scroll in steps if required.
// Returns an error if the configuration is invalid, the target element is not ready, or the adapter scroll operation fails.
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

// doScrollTo performs a smooth scrolling operation to a specified target element or coordinates with adjustable steps and interval.
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

// doNavigate constructs a full URL if the target is relative and navigates to the specified target using the adapter.
// Returns an error if the navigation operation fails.
func (e *Executor) doNavigate(target string) error {
	if strings.HasPrefix(target, "/") {
		target = e.baseUrl + target
	}
	return e.adapter.Navigate(target)
}

// doDownload performs a file download by constructing the full target URL, retrieving cookies, and invoking the downloader.
// target specifies the file's URL path or endpoint.
// value is the filename or destination for the downloaded file.
// Returns an error if URL construction, cookie retrieval, or download process fails.
func (e *Executor) doDownload(target string, value string) error {
	if strings.HasPrefix(target, "/") {
		target = e.baseUrl + target
	}
	cookies, err := e.adapter.GetCookiesHttp()
	if err != nil {
		return err
	}
	err = e.downloader.Do(target, cookies, value)
	if err != nil {
		return err
	}
	return nil
}

// commandExec executes a given ConfigCommand using Executor's adapter and specified logic for diverse command types.
// It performs various browser or system-related actions based on the command type, target, and additional parameters.
// Returns an error if the command execution fails or is unimplemented.
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
	case "download":
		err = e.doDownload(target, value)
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
	case "clickDownload":
		err = e.doClickDownload(target, value)
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
		err = e.doGetAllCookies(target, id, value)
	case "getCookie":
		err = e.doGetCookie(target, id, value)
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

// finalize ensures all timers associated with the executor are finalized and logs the operation's result and duration.
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

// doComputeJump finds the index of the target command in the commands list based on the given target ID.
// Returns the index if found, or -1 with an error if the target is not found or is empty.
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

// doCommand executes a specified command from the given list at the provided index and returns its result or error.
func (e *Executor) doCommand(commands []ConfigCommand, idx int) (string, int, error) {
	if idx >= len(commands) {
		return "", -1, fmt.Errorf("invalid index %d", idx)
	}
	cmd := commands[idx]
	jump := -1
	if e.adapter.IsDebugEnabled() {
		l1 := fmt.Sprintf("--------------------------------------------------------------")
		l2 := fmt.Sprintf("Next command is [%s] %s: %s", cmd.Id, cmd.Command, cmd.Target)
		if e.sam != nil {
			e.sam.Print(l1)
			e.sam.Print(l2)
		} else {
			fmt.Printf("%s\n", l1)
			fmt.Printf("%s\n", l2)
		}
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
		var l1 string
		if err != nil {
			l1 = fmt.Sprintf("Error: %v", err)
		} else {
			l1 = fmt.Sprintf("Ok")
		}
		if e.sam != nil {
			e.sam.Print(l1)
		} else {
			fmt.Printf("%s\n", l1)
		}
	}
	if e.adapter.IsErrorDisabled() {
		err = nil
	}
	if err != nil {
		return cmd.Id, -1, err
	}
	return "", jump, err
}

// commandsLoop iterates through test commands and executes them, handling jumps and errors during execution.
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
		}
	}
	return "", nil
}

// Start initializes and begins execution of the Executor with optional SAM mode based on the input parameter.
// It performs setup actions such as logging, applying templates, loading URLs, and starting the command loop or SAM.
// Returns an error if initialization or execution encounters failures.
func (e *Executor) Start(sam bool) error {
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
	if !sam {
		var id string
		if id, err = e.commandsLoop(); err != nil {
			e.finalize(fmt.Errorf("%s (%s)", err.Error(), id))
			return err
		}
		e.finalize(nil)
	} else {
		e.sam = NewSam(e)
		if err = e.sam.Start(); err != nil {
			return err
		}
	}
	return nil
}

// SamSendMessage sends an ISamMessage instance to the configured SamConsole if it is not nil.
func (e *Executor) SamSendMessage(s ISamMessage) {
	if e.sam == nil {
		return
	}
	e.sam.SendMessage(s)
}
