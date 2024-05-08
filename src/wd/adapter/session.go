package adapter

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/markel1974/webautoma/src/wd/adapter/bidi"
)

type Session struct {
	ID           string       `json:"sessionId"`
	Capabilities Capabilities `json:"capabilities"`
	client       *RestClient
	biDiSession  *bidi.Session
}

func (s *Session) IsBiDiSession() bool {
	return s.Capabilities.WebSocketURL() != ""
}

func (s *Session) BiDiSession() (*bidi.Session, error) {
	if !s.IsBiDiSession() {
		return nil, errors.New("BiDi not supported")
	}

	return s.biDiSession, nil
}

type Timeouts struct {
	// Session script timeout in milliseconds.
	// Determines when to interrupt a script that is being evaluated.
	Script int `json:"script"`

	// Session page load timeout in milliseconds.
	// Provides the timeout limit used to interrupt navigation of the browsing context.
	PageLoad int `json:"pageLoad"`

	// Session implicit wait timeout in milliseconds.
	// Gives the timeout of when to abort locating an element.
	Implicit int `json:"implicit"`
}

// GetTimeouts gets timeout durations associated with the current session.
func (s *Session) GetTimeouts() (*Timeouts, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/timeouts", s.ID))
	if err != nil {
		return nil, err
	}

	timeouts := &Timeouts{}
	err = json.Unmarshal(data, &timeouts)

	return timeouts, err
}

// SetTimeouts configures the amount of time that a particular type of operation can execute for before
// they are aborted and a Timeout error is returned to the client.
func (s *Session) SetTimeouts(timeouts *Timeouts) error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/timeouts", s.ID), &Params{
		"script":   timeouts.Script,
		"pageLoad": timeouts.PageLoad,
		"implicit": timeouts.Implicit,
	})

	return err
}

// Close closes the session.
func (s *Session) Close() error {
	_, err := s.client.Delete(fmt.Sprintf("/session/%s", s.ID))
	return err
}

// NavigateTo navigates to a new URL.
func (s *Session) NavigateTo(url string) error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/url", s.ID), &Params{"url": url})
	return err
}

// GetCurrentURL gets current page URL.
func (s *Session) GetCurrentURL() (string, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/url", s.ID))
	if err != nil {
		return "", err
	}

	var url string
	err = json.Unmarshal(data, &url)

	return url, err
}

// Back navigates to previous url from history.
func (s *Session) Back() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/back", s.ID), nil)
	return err
}

// Forward navigates forward to next url from history.
func (s *Session) Forward() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/forward", s.ID), nil)
	return err
}

// Refresh refreshes the current page.
func (s *Session) Refresh() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/refresh", s.ID), nil)
	return err
}

// GetTitle gets the current page title.
func (s *Session) GetTitle() (string, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/title", s.ID))
	if err != nil {
		return "", err
	}

	var title string
	err = json.Unmarshal(data, &title)

	return title, err
}

// GetWindowHandle gets handle of current window.
func (s *Session) GetWindowHandle() (string, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/window", s.ID))
	if err != nil {
		return "", err
	}

	var handle string
	err = json.Unmarshal(data, &handle)

	return handle, err
}

// CloseWindow closes the current window.
func (s *Session) CloseWindow() error {
	_, err := s.client.Delete(fmt.Sprintf("/session/%s/window", s.ID))
	return err
}

// SwitchToWindow changes focus to another window. The window to change focus to may be specified
// by it's server assigned window handle.
func (s *Session) SwitchToWindow(handle string) error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/window", s.ID), &Params{"handle": handle})
	return err
}

// GetWindowHandles gets all window handles.
func (s *Session) GetWindowHandles() ([]string, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/window", s.ID))
	if err != nil {
		return nil, err
	}

	var handles []string
	err = json.Unmarshal(data, &handles)

	return handles, err
}

// SwitchToFrame changes focus to another frame on the page.
func (s *Session) SwitchToFrame(target Element) error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/frame", s.ID), &Params{"id": target.ID})
	return err
}

// SwitchToParentFrame changes focus to parent frame on the page.
func (s *Session) SwitchToParentFrame() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/frame/parent", s.ID), nil)
	return err
}

// WindowRect defines the Window Rect.
type WindowRect struct {
	// The screenX and screenLeft attributes must return the x-coordinate, relative to the origin of the
	// Web-exposed screen area, of the left of the client window as number of CSS pixels
	X int `json:"x"`

	// The screenY and screenTop attributes must return the y-coordinate, relative to the origin of the
	// screen of the Web-exposed screen area, of the top of the client window as number of CSS pixels
	Y int `json:"y"`

	// The outerWidth attribute must return the width of the client window.
	// If there is no client window this attribute must return zero
	Width int `json:"width"`

	// The outerWidth attribute must return the height of the client window.
	// If there is no client window this attribute must return zero
	Height int `json:"height"`
}

// GetWindowRect gets the size and position on the screen of the operating system window.
func (s *Session) GetWindowRect() (*WindowRect, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/window/rect", s.ID))
	if err != nil {
		return nil, err
	}

	windowRect := &WindowRect{}
	err = json.Unmarshal(data, &windowRect)

	return windowRect, err
}

// SetWindowRect sets the size and position on the screen of the operating system window.
func (s *Session) SetWindowRect(windowRect *WindowRect) error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/window/rect", s.ID), &Params{
		"x":      windowRect.X,
		"y":      windowRect.Y,
		"width":  windowRect.Width,
		"height": windowRect.Height,
	})

	return err
}

// MaximizeWindow maximizes the current window.
func (s *Session) MaximizeWindow() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/window/maximize", s.ID), nil)
	return err
}

// MinimizeWindow minimizes the current window.
func (s *Session) MinimizeWindow() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/window/minimize", s.ID), nil)
	return err
}

// FullscreenWindow increases current window to Full-Screen.
func (s *Session) FullscreenWindow() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/window/fullscreen", s.ID), nil)
	return err
}

// FindElement searches for an element on the page, starting from the document root.
func (s *Session) FindElement(strategy LocatorStrategy, selector string) (*Element, error) {
	data, err := s.client.Post(fmt.Sprintf("/session/%s/element", s.ID), &Params{
		"using": strategy,
		"value": selector,
	})
	if err != nil {
		return nil, err
	}

	element := Element{}
	if err := json.Unmarshal(data, &element); err != nil {
		return nil, err
	}

	element.SessionID = s.ID
	element.client = s.client

	return &element, nil
}

// FindElements searches for multiple elements on the page, starting from the document root. The
// located elements will be returned as a WebElement JSON objects. The table below lists the locator
// strategies that each server should support. Elements should be returned in the order located
// in the DOM.
func (s *Session) FindElements(strategy LocatorStrategy, selector string) ([]Element, error) {
	data, err := s.client.Post(fmt.Sprintf("/session/%s/elements", s.ID), &Params{
		"using": strategy,
		"value": selector,
	})
	if err != nil {
		return nil, err
	}
	var elements []Element
	if err := json.Unmarshal(data, &elements); err != nil {
		return nil, err
	}
	for _, element := range elements {
		element.SessionID = s.ID
		element.client = s.client
	}
	return elements, nil
}

// GetActiveElement gets the element on the page that currently has focus.
func (s *Session) GetActiveElement() (*Element, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/element/active", s.ID))
	if err != nil {
		return nil, err
	}

	element := Element{}
	if err := json.Unmarshal(data, &element); err != nil {
		return nil, err
	}

	element.SessionID = s.ID
	element.client = s.client

	return &element, nil
}

// GetPageSource returns a string serialization of the DOM of the current browsing context active document.
func (s *Session) GetPageSource() (string, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/source", s.ID))
	if err != nil {
		return "", err
	}
	var source string
	err = json.Unmarshal(data, &source)
	return source, err
}

// ExecuteScript injects a snippet of JavaScript into the page for execution in the context
// of the currently selected frame. The executed script is assumed to be synchronous and
// the result of evaluating the script is returned to the client.
func (s *Session) ExecuteScript(script string, args []interface{}) ([]byte, error) {
	data, err := s.client.Post(fmt.Sprintf("/session/%s/execute/sync", s.ID), &Params{
		"script": script,
		"args":   args,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Cookie struct {
	// The name of the cookie.
	Name string `json:"name"`

	// The cookie value.
	Value string `json:"value"`

	// The cookie path. Defaults to "/" if omitted when adding a cookie.
	Path string `json:"path"`

	// The domain the cookie is visible to.
	// Defaults to the current browsing context's document's URL domain if omitted when adding a cookie.
	Domain string `json:"domain"`

	// Whether the cookie is a secure cookie. Defaults to false if omitted when adding a cookie.
	Secure bool `json:"secure"`

	// Whether the cookie is an HTTP only cookie. Defaults to false if omitted when adding a cookie.
	HTTPOnly bool `json:"httpOnly"`

	// When the cookie expires, specified in seconds since Unix Epoch.
	// Defaults to 20 years into the future if omitted when adding a cookie.
	Expiry int `json:"expiry"`
}

// GetCookies returns all cookies associated with the address of the current browsing context's
// active document.
func (s *Session) GetCookies() ([]*Cookie, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/cookie", s.ID))
	if err != nil {
		return nil, err
	}

	var cookies []*Cookie
	err = json.Unmarshal(data, &cookies)

	return cookies, err
}

// GetCookie returns cookie based on the cookie name
func (s *Session) GetCookie(name string) (*Cookie, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/cookie/%s", s.ID, name))
	if err != nil {
		return nil, err
	}

	var cookie *Cookie
	err = json.Unmarshal(data, cookie)

	return cookie, err
}

// AddCookie adds a single cookie to the cookie store associated with the active document's address.
func (s *Session) AddCookie(cookie Cookie) error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/execute/sync", s.ID), &Params{
		"name":     cookie.Name,
		"value":    cookie.Value,
		"path":     cookie.Path,
		"domain":   cookie.Domain,
		"secure":   cookie.Secure,
		"httpOnly": cookie.HTTPOnly,
		"expiry":   cookie.Expiry,
	})

	return err
}

// DeleteCookie deletes a cookie based on its name
func (s *Session) DeleteCookie(name string) error {
	_, err := s.client.Delete(fmt.Sprintf("/session/%s/cookie/%s", s.ID, name))
	return err
}

// DeleteCookies deletes all cookies associated with the address of the current browsing context's
// active document.
func (s *Session) DeleteCookies() error {
	_, err := s.client.Delete(fmt.Sprintf("/session/%s/cookie", s.ID))
	return err
}

type ActionType string

const (
	ActionTypePause       ActionType = "pause"
	ActionTypeKeyDown     ActionType = "keyDown"
	ActionTypeKeyUp       ActionType = "keyUp"
	ActionTypePointerMove ActionType = "pointerMove"
	ActionTypePointerUp   ActionType = "pointerUp"
	ActionTypePointerDown ActionType = "pointerDown"
)

// DismissAlert dismisses the alert in current page.
func (s *Session) DismissAlert() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/alert/dismiss", s.ID), nil)
	return err
}

// AcceptAlert accepts the alert in current page.
func (s *Session) AcceptAlert() error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/alert/accept", s.ID), nil)
	return err
}

// GetAlertText returns the text from an alert.
func (s *Session) GetAlertText() (string, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/alert/text", s.ID))
	if err != nil {
		return "", err
	}
	var text string
	err = json.Unmarshal(data, &text)
	return text, err
}

// SendAlertText sets the text field of a prompt to the given value.
func (s *Session) SendAlertText(text string) error {
	_, err := s.client.Post(fmt.Sprintf("/session/%s/alert/text", s.ID), &Params{
		"text": text,
	})
	return err
}

// TakeScreenshot takes a screenshot of the top-level browsing context's viewport.
func (s *Session) TakeScreenshot() ([]byte, error) {
	data, err := s.client.Get(fmt.Sprintf("/session/%s/screenshot", s.ID))
	if err != nil {
		return nil, err
	}
	var screenshot string
	if err := json.Unmarshal(data, &screenshot); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(screenshot)
}
