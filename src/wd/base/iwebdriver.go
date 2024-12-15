package base

import (
	"time"
)

type Condition func(wd IWebDriver) (bool, error)

type IWebDriver interface {
	// Status returns various pieces of information about the server environment.
	Status() (*Status, error)

	// NewSession starts a new session and returns the session ID.
	NewSession() (string, error)

	// SessionID returns the current session ID.
	SessionID() string

	// SwitchSession switches to the given session ID.
	SwitchSession(sessionID string) error

	// Capabilities returns the current session's capabilities.
	Capabilities() (Capabilities, error)

	// SetAsyncScriptTimeout sets the amount of time that asynchronous scripts
	// are permitted to run before they are aborted. The timeout will be rounded
	// to nearest millisecond.
	SetAsyncScriptTimeout(timeout time.Duration) error

	// SetImplicitWaitTimeout sets the amount of time the driver should wait when
	// searching for elements. The timeout will be rounded to nearest millisecond.
	SetImplicitWaitTimeout(timeout time.Duration) error

	// SetPageLoadTimeout sets the amount of time the driver should wait when
	// loading a page. The timeout will be rounded to nearest millisecond.
	SetPageLoadTimeout(timeout time.Duration) error

	// Quit ends the current session. The browser instance will be closed.
	Quit() error

	// CurrentWindowHandle returns the ID of current window handle.
	CurrentWindowHandle() (string, error)

	// WindowHandles returns the IDs of current open windows.
	WindowHandles() ([]string, error)

	// CurrentURL returns the browser's current URL.
	CurrentURL() (string, error)

	// Title returns the current page's title.
	Title() (string, error)

	// PageSource returns the current page's source.
	PageSource() (string, error)

	// Close closes the current window.
	Close() error

	// SwitchFrame switches to the given frame. The frame parameter can be the
	// frame's ID as a string, its WebElement instance as returned by
	// GetElement, or nil to switch to the current top-level browsing context.
	SwitchFrame(frame interface{}) error

	// SwitchParentFrame switches to the parent frame.
	SwitchParentFrame() error

	// SwitchWindow switches the context to the specified window.
	SwitchWindow(handle string) error

	// CloseWindow closes the specified window.
	CloseWindow(handle string) error

	// MaximizeWindow maximizes a window. If the name is empty, the current
	// window will be maximized.
	MaximizeWindow(handle string) error

	// ResizeWindow changes the dimensions of a window. If the name is empty, the
	// current window will be maximized.
	ResizeWindow(handle string, width, height int) error

	// Navigate Navigates the browser to the provided URL.
	Navigate(url string) error

	// Forward moves in history.
	Forward() error

	// Back moves backward in history.
	Back() error

	// Refresh refreshes the page.
	Refresh() error

	// ComputedLabel ...
	ComputedLabel() (string, error)

	// FindElement finds exactly one element in the current page's DOM.
	FindElement(by string, value string) (IWebElement, error)

	// FindElements finds potentially many elements in the current page's DOM.
	FindElements(by string, value string) ([]IWebElement, error)

	// ActiveElement returns the currently active element on the page.
	ActiveElement() (IWebElement, error)

	// DecodeElement decodes a single element response.
	DecodeElement([]byte) (IWebElement, error)

	// DecodeElements decodes a multi-element response.
	DecodeElements([]byte) ([]IWebElement, error)

	// GetCookies returns all the cookies in the browser's jar.
	GetCookies() ([]Cookie, error)

	// GetCookie returns the named cookie in the jar, if present. This method is
	// only implemented for Firefox.
	GetCookie(name string) (Cookie, error)

	// AddCookie adds a cookie to the browser's jar.
	AddCookie(cookie *Cookie) error

	// DeleteAllCookies deletes all the cookies in the browser's jar.
	DeleteAllCookies() error

	// DeleteCookie deletes a cookie to the browser's jar.
	DeleteCookie(name string) error

	// Click clicks a mouse button. The button should be one of RightButton,
	// MiddleButton or LeftButton.
	Click(button int) error

	// DoubleClick clicks the left mouse button twice.
	DoubleClick() error

	// ButtonDown causes the left mouse button to be held down.
	ButtonDown() error

	// ButtonUp causes the left mouse button to be released.
	ButtonUp() error

	// StoreKeyActions store provided actions until they are executed
	// by PerformActions or released by ReleaseActions.
	// inputID is a string used as a unique virtual device identifier for this
	// and future actions, the value can be set to any valid string
	// and used to refer to this specific device in future calls.
	StoreKeyActions(inputID string, actions ...KeyAction)

	// StorePointerActions store provided actions until they are executed
	// by PerformActions or released by ReleaseActions.
	// inputID is a string used as a unique virtual device identifier for this
	// and future actions, the value can be set to any valid string
	// and used to refer to this specific device in future calls.
	StorePointerActions(inputID string, pointer PointerType, actions ...PointerAction)

	// StoreWheelActions store provided actions until they are executed
	// by PerformActions or released by ReleaseActions.
	// inputID is a string used as a unique virtual device identifier for this
	// and future actions, the value can be set to any valid string
	// and used to refer to this specific device in future calls.
	StoreWheelActions(inputID string, actions ...WheelAction)

	// PerformActions executes actions previously stored by calls to StorePointerActions and StoreKeyActions.
	PerformActions() error

	// ReleaseActions releases keys and pointer buttons if they are pressed,
	// triggering any events as if they were performed by a regular action.
	ReleaseActions() error

	// SendModifier sends the modifier key to the active element. The modifier
	// can be one of ShiftKey, ControlKey, AltKey, MetaKey.
	//
	// Deprecated: Use KeyDown or KeyUp instead.
	SendModifier(modifier string, isDown bool) error

	// KeyDown sends a sequence of keystrokes to the active element. This method
	// is similar to SendKeys but without the implicit termination. Modifiers are
	// not released at the end of each call.
	KeyDown(keys string) error

	// KeyUp indicates that a previous keystroke sent by KeyDown should be
	// released.
	KeyUp(keys string) error

	// Screenshot takes a screenshot of the browser window.
	Screenshot() ([]byte, error)

	// Log fetches the logs. Log types must be previously configured in the
	// capabilities.
	Log(kind LogType) ([]LogMessage, error)

	// DismissAlert dismisses current alert.
	DismissAlert() error

	// AcceptAlert accepts the current alert.
	AcceptAlert() error

	// AlertText returns the current alert text.
	AlertText() (string, error)

	// SetAlertText sets the current alert text.
	SetAlertText(text string) error

	// ExecuteScript executes a script.
	ExecuteScript(script string, args []interface{}) (interface{}, error)

	// ExecuteScriptAsync asynchronously executes a script.
	ExecuteScriptAsync(script string, args []interface{}) (interface{}, error)

	// ExecuteScriptRaw executes a script but does not perform JSON decoding.
	ExecuteScriptRaw(script string, args []interface{}) ([]byte, error)

	// ExecuteScriptAsyncRaw asynchronously executes a script but does not
	// perform JSON decoding.
	ExecuteScriptAsyncRaw(script string, args []interface{}) ([]byte, error)

	// WaitWithTimeoutAndInterval waits for the condition to evaluate to true.
	WaitWithTimeoutAndInterval(condition Condition, timeout, interval time.Duration) error

	// WaitWithTimeout works like WaitWithTimeoutAndInterval, but with default polling interval.
	WaitWithTimeout(condition Condition, timeout time.Duration) error

	//Wait works like WaitWithTimeoutAndInterval, but using the default timeout and polling interval.
	Wait(condition Condition) error
}
