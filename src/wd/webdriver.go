// See https://www.w3.org/TR/webdriver for the protocol.

package wd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/markel1974/webautoma/src/wd/base"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/blang/semver"
)

type ServerReply struct {
	SessionID *string
	Value     json.RawMessage
	Status    int
	State     string
	base.Error
}

// Actions stores KeyActions and PointerActions for later execution.
type Actions []map[string]interface{}

const (
	DefaultContentType = "application/json"

	DefaultWaitInterval = 100 * time.Millisecond

	DefaultWaitTimeout = 60 * time.Second
)

const (
	// legacyWebElementIdentifier is the string constant used in the old
	// WebDriver JSON protocol that is the key for the map that contains an
	// unique element identifier.
	legacyWebElementIdentifier = "ELEMENT"

	// webElementIdentifier is the string constant defined by the W3C
	// specification that is the key for the map that contains a unique element identifier.
	webElementIdentifier = "element-6066-11e4-a52e-4f735466cecf"
)

func elementIDFromValue(v map[string]string) string {
	for _, key := range []string{webElementIdentifier, legacyWebElementIdentifier} {
		v, ok := v[key]
		if !ok || v == "" {
			continue
		}
		return v
	}
	return ""
}

func parseVersion(v string) (semver.Version, error) {
	parts := strings.Split(v, ".")
	var err error
	for i := len(parts); i > 0; i-- {
		var ver semver.Version
		ver, err = semver.ParseTolerant(strings.Join(parts[:i], "."))
		if err == nil {
			return ver, nil
		}
	}
	return semver.Version{}, err
}

type WebDriver struct {
	id             string
	urlPrefix      string
	capabilities   base.Capabilities
	w3cCompatible  bool
	storedActions  Actions
	browser        string
	browserVersion semver.Version
	client         *http.Client
	debug          bool
}

func NewWebDriver(capabilities base.Capabilities, wdUrl *url.URL, debug bool) *WebDriver {
	wd := &WebDriver{
		urlPrefix:    wdUrl.String(),
		capabilities: capabilities,
		client:       http.DefaultClient,
		debug:        debug,
	}
	if b := capabilities["browserName"]; b != nil {
		wd.browser = b.(string)
	}
	return wd
}

func (wd *WebDriver) Start() error {
	if _, err := wd.NewSession(); err != nil {
		return err
	}
	return nil
}

func (wd *WebDriver) DeleteSession(urlPrefix, id string) error {
	u, err := url.Parse(urlPrefix)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "session", id)
	return wd.voidParamsCommand("DELETE", u.String(), nil)
}

func (wd *WebDriver) Status() (*base.Status, error) {
	rUrl := wd.requestURL("/status")
	reply, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return nil, err
	}
	status := new(struct{ Value base.Status })
	if err := json.Unmarshal(reply, status); err != nil {
		return nil, err
	}
	return &status.Value, nil
}

func (wd *WebDriver) NewSession() (string, error) {
	// Detect whether the remote end complies with the W3C specification:
	// non-compliant implementations use the top-level 'desiredCapabilities' JSON
	// key, whereas the specification mandates the 'capabilities' key.
	//
	// However, Selenium 3 currently does not implement this part of the specification.
	// https://github.com/SeleniumHQ/selenium/issues/2827
	//
	// TODO: audit which ones of these are still relevant. The W3C
	// standard switched to the "alwaysMatch" version in February 2017.
	attempts := []struct{ params map[string]interface{} }{
		{map[string]interface{}{
			"capabilities":        newW3CCapabilities(wd.capabilities),
			"desiredCapabilities": wd.capabilities,
		}},
		{map[string]interface{}{
			"capabilities": map[string]interface{}{
				"desiredCapabilities": wd.capabilities,
			},
		}},
		{map[string]interface{}{
			"desiredCapabilities": wd.capabilities,
		}}}

	for i, s := range attempts {
		data, err := json.Marshal(s.params)
		if err != nil {
			return "", err
		}

		response, err := wd.execute("POST", wd.requestURL("/session"), data)
		if err != nil {
			return "", err
		}

		reply := new(ServerReply)
		if err := json.Unmarshal(response, reply); err != nil {
			if i < len(attempts) {
				continue
			}
			return "", err
		}

		if reply.Status != 0 && i < len(attempts) {
			continue
		}
		if reply.SessionID != nil {
			wd.id = *reply.SessionID
		}

		if len(reply.Value) > 0 {
			type returnedCapabilities struct {
				// firefox via geckodriver: 55.0a1
				BrowserVersion string
				// chrome via chromedriver: 61.0.3116.0
				// firefox via selenium 2: 45.9.0
				// htmlunit: 9.4.3.v20170317
				Version          string
				PageLoadStrategy string
				Proxy            base.Proxy
				Timeouts         struct {
					Implicit       float32
					PageLoadLegacy float32 `json:"page load"`
					PageLoad       float32
					Script         float32
				}
			}

			value := struct {
				SessionID string

				// The W3C specification moved most of the returned data into the
				// "capabilities" field.
				Capabilities *returnedCapabilities

				// Legacy implementations returned most data directly in the "values"
				// key.
				returnedCapabilities
			}{}

			if err := json.Unmarshal(reply.Value, &value); err != nil {
				return "", fmt.Errorf("error unmarshalling value: %v", err)
			}
			if value.SessionID != "" && wd.id == "" {
				wd.id = value.SessionID
			}
			var cs returnedCapabilities
			if value.Capabilities != nil {
				cs = *value.Capabilities
				wd.w3cCompatible = true
			} else {
				cs = value.returnedCapabilities
			}

			for _, s := range []string{cs.Version, cs.BrowserVersion} {
				if s == "" {
					continue
				}
				v, err := parseVersion(s)
				if err != nil {
					//debugLog("error parsing version: %v\n", err)
					continue
				}
				wd.browserVersion = v
			}
		}
		return wd.id, nil
	}
	return "", fmt.Errorf("unreachable")
}

func (wd *WebDriver) SessionID() string {
	return wd.id
}

func (wd *WebDriver) SwitchSession(sessionID string) error {
	wd.id = sessionID
	return nil
}

func (wd *WebDriver) Capabilities() (base.Capabilities, error) {
	rUrl := wd.requestURL("/session/%s", wd.id)
	response, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return nil, err
	}

	c := new(struct{ Value base.Capabilities })
	if err := json.Unmarshal(response, c); err != nil {
		return nil, err
	}

	return c.Value, nil
}

func (wd *WebDriver) SetAsyncScriptTimeout(timeout time.Duration) error {
	if !wd.w3cCompatible {
		return wd.voidCommand("/session/%s/timeouts/async_script", map[string]uint{
			"ms": uint(timeout / time.Millisecond),
		})
	}
	return wd.voidCommand("/session/%s/timeouts", map[string]uint{
		"script": uint(timeout / time.Millisecond),
	})
}

func (wd *WebDriver) SetImplicitWaitTimeout(timeout time.Duration) error {
	if !wd.w3cCompatible {
		return wd.voidCommand("/session/%s/timeouts/implicit_wait", map[string]uint{
			"ms": uint(timeout / time.Millisecond),
		})
	}
	return wd.voidCommand("/session/%s/timeouts", map[string]uint{
		"implicit": uint(timeout / time.Millisecond),
	})
}

func (wd *WebDriver) SetPageLoadTimeout(timeout time.Duration) error {
	if !wd.w3cCompatible {
		return wd.voidCommand("/session/%s/timeouts", map[string]interface{}{
			"ms":   uint(timeout / time.Millisecond),
			"type": "page load",
		})
	}
	return wd.voidCommand("/session/%s/timeouts", map[string]uint{
		"pageLoad": uint(timeout / time.Millisecond),
	})
}

func (wd *WebDriver) Quit() error {
	if wd.id == "" {
		return nil
	}
	_, err := wd.execute("DELETE", wd.requestURL("/session/%s", wd.id), nil)
	if err == nil {
		wd.id = ""
	}
	return err
}

func (wd *WebDriver) CurrentWindowHandle() (string, error) {
	if !wd.w3cCompatible {
		return wd.stringCommand("/session/%s/window_handle")
	}
	return wd.stringCommand("/session/%s/window")
}

func (wd *WebDriver) WindowHandles() ([]string, error) {
	if !wd.w3cCompatible {
		return wd.stringsCommand("/session/%s/window_handles")
	}
	return wd.stringsCommand("/session/%s/window/handles")
}

func (wd *WebDriver) CurrentURL() (string, error) {
	rUrl := wd.requestURL("/session/%s/url", wd.id)
	response, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return "", err
	}
	reply := new(struct{ Value *string })
	if err := json.Unmarshal(response, reply); err != nil {
		return "", err
	}

	return *reply.Value, nil
}

func (wd *WebDriver) Navigate(rUrl string) error {
	if !strings.HasPrefix(strings.ToLower(rUrl), "http") {
		rUrl = "https://" + rUrl
	}
	rUrl = strings.TrimSpace(rUrl)
	requestURL := wd.requestURL("/session/%s/url", wd.id)
	params := map[string]string{"url": rUrl}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, err = wd.execute("POST", requestURL, data)
	return err
}

func (wd *WebDriver) Forward() error {
	return wd.voidCommand("/session/%s/forward", nil)
}

func (wd *WebDriver) Back() error {
	return wd.voidCommand("/session/%s/back", nil)
}

func (wd *WebDriver) Refresh() error {
	return wd.voidCommand("/session/%s/refresh", nil)
}

func (wd *WebDriver) Title() (string, error) {
	return wd.stringCommand("/session/%s/title")
}

func (wd *WebDriver) PageSource() (string, error) {
	return wd.stringCommand("/session/%s/source")
}

func (wd *WebDriver) DecodeElement(data []byte) (base.IWebElement, error) {
	reply := new(struct{ Value map[string]string })
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, err
	}

	id := elementIDFromValue(reply.Value)
	if id == "" {
		return nil, fmt.Errorf("invalid element returned: %+v", reply)
	}
	return &WebElement{
		parent: wd,
		id:     id,
	}, nil
}

func (wd *WebDriver) DecodeElements(data []byte) ([]base.IWebElement, error) {
	reply := new(struct{ Value []map[string]string })
	if err := json.Unmarshal(data, reply); err != nil {
		return nil, err
	}

	elements := make([]base.IWebElement, len(reply.Value))
	for i, elem := range reply.Value {
		id := elementIDFromValue(elem)
		if id == "" {
			return nil, fmt.Errorf("invalid element returned: %+v", reply)
		}
		elements[i] = &WebElement{
			parent: wd,
			id:     id,
		}
	}

	return elements, nil
}

func (wd *WebDriver) FindElement(by, value string) (base.IWebElement, error) {
	response, err := wd.find(by, value, "", "")
	if err != nil {
		return nil, err
	}
	return wd.DecodeElement(response)
}

func (wd *WebDriver) FindElements(by, value string) ([]base.IWebElement, error) {
	response, err := wd.find(by, value, "s", "")
	if err != nil {
		return nil, err
	}

	return wd.DecodeElements(response)
}

func (wd *WebDriver) Close() error {
	rUrl := wd.requestURL("/session/%s/window", wd.id)
	_, err := wd.execute("DELETE", rUrl, nil)
	return err
}

func (wd *WebDriver) SwitchWindow(handle string) error {
	params := make(map[string]string)
	if !wd.w3cCompatible {
		params["name"] = handle
	} else {
		params["handle"] = handle
	}
	return wd.voidCommand("/session/%s/window", params)
}

func (wd *WebDriver) CloseWindow(name string) error {
	return wd.modifyWindow(name, "DELETE", "", nil)
}

func (wd *WebDriver) MaximizeWindow(name string) error {
	if !wd.w3cCompatible {
		if name != "" {
			var err error
			name, err = wd.CurrentWindowHandle()
			if err != nil {
				return err
			}
		}
		rUrl := wd.requestURL("/session/%s/window/%s/maximize", wd.id, name)
		_, err := wd.execute("POST", rUrl, nil)
		return err
	}
	return wd.modifyWindow(name, "POST", "maximize", map[string]string{})
}

func (wd *WebDriver) MinimizeWindow(name string) error {
	return wd.modifyWindow(name, "POST", "minimize", map[string]string{})
}

func (wd *WebDriver) ResizeWindow(name string, width, height int) error {
	if !wd.w3cCompatible {
		return wd.modifyWindow(name, "POST", "size", map[string]int{
			"width":  width,
			"height": height,
		})
	}
	return wd.modifyWindow(name, "POST", "rect", map[string]float64{
		"width":  float64(width),
		"height": float64(height),
	})
}

func (wd *WebDriver) SwitchFrame(frame interface{}) error {
	params := map[string]interface{}{}
	switch f := frame.(type) {
	case WebElement, int, nil:
		params["id"] = f
	case string:
		if f == "" {
			params["id"] = nil
		} else if wd.w3cCompatible {
			e, err := wd.FindElement(base.ByID, f)
			if err != nil {
				return err
			}
			params["id"] = e
		} else { // Legacy, non W3C-spec behavior.
			params["id"] = f
		}
	default:
		return fmt.Errorf("invalid type %T", frame)
	}
	return wd.voidCommand("/session/%s/frame", params)
}

func (wd *WebDriver) ActiveElement() (base.IWebElement, error) {
	verb := "GET"
	if wd.browser == "firefox" && wd.browserVersion.Major < 47 {
		verb = "POST"
	}
	rUrl := wd.requestURL("/session/%s/element/active", wd.id)
	response, err := wd.execute(verb, rUrl, nil)
	if err != nil {
		return nil, err
	}
	return wd.DecodeElement(response)
}

func (wd *WebDriver) GetCookie(name string) (base.Cookie, error) {
	if wd.browser == "chrome" {
		cs, err := wd.GetCookies()
		if err != nil {
			return base.Cookie{}, err
		}
		for _, c := range cs {
			if c.Name == name {
				return c, nil
			}
		}
		return base.Cookie{}, errors.New("cookie not found")
	}
	rUrl := wd.requestURL("/session/%s/cookie/%s", wd.id, name)
	data, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return base.Cookie{}, err
	}

	// GeckoDriver returns a list of cookies for this method. Try both a single
	// cookie and a list.
	//
	// https://github.com/mozilla/geckodriver/issues/761
	reply := new(struct{ Value base.Cookie })
	if err := json.Unmarshal(data, reply); err == nil {
		return reply.Value.Sanitize(), nil
	}
	listReply := new(struct{ Value []base.Cookie })
	if err := json.Unmarshal(data, listReply); err != nil {
		return base.Cookie{}, err
	}
	if len(listReply.Value) == 0 {
		return base.Cookie{}, errors.New("no cookies returned")
	}
	return listReply.Value[0].Sanitize(), nil
}

func (wd *WebDriver) GetCookies() ([]base.Cookie, error) {
	rUrl := wd.requestURL("/session/%s/cookie", wd.id)
	data, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return nil, err
	}

	reply := new(struct{ Value []base.Cookie })
	if err := json.Unmarshal(data, reply); err != nil {
		return nil, err
	}

	cookies := make([]base.Cookie, len(reply.Value))
	for i, c := range reply.Value {
		cookies[i] = c.Sanitize()
	}
	return cookies, nil
}

func (wd *WebDriver) AddCookie(cookie *base.Cookie) error {
	return wd.voidCommand("/session/%s/cookie", map[string]*base.Cookie{"cookie": cookie})
}

func (wd *WebDriver) DeleteAllCookies() error {
	rUrl := wd.requestURL("/session/%s/cookie", wd.id)
	_, err := wd.execute("DELETE", rUrl, nil)
	return err
}

func (wd *WebDriver) DeleteCookie(name string) error {
	rUrl := wd.requestURL("/session/%s/cookie/%s", wd.id, name)
	_, err := wd.execute("DELETE", rUrl, nil)
	return err
}

func (wd *WebDriver) Click(button int) error {
	return wd.voidCommand("/session/%s/click", map[string]int{"button": button})
}

func (wd *WebDriver) DoubleClick() error {
	return wd.voidCommand("/session/%s/doubleclick", nil)
}

func (wd *WebDriver) ButtonDown() error {
	return wd.voidCommand("/session/%s/buttondown", nil)
}

func (wd *WebDriver) ButtonUp() error {
	return wd.voidCommand("/session/%s/buttonup", nil)
}

func (wd *WebDriver) SendModifier(modifier string, isDown bool) error {
	if isDown {
		return wd.KeyDown(modifier)
	}
	return wd.KeyUp(modifier)
}

func (wd *WebDriver) KeyDown(keys string) error {
	// Selenium implemented the actions API but has not yet updated its new session response.
	if !wd.w3cCompatible && !(wd.browser == "firefox" && wd.browserVersion.Major > 47) {
		return wd.voidCommand("/session/%s/keys", wd.processKeyString(keys))
	}
	return wd.keyAction("keyDown", keys)
}

func (wd *WebDriver) KeyUp(keys string) error {
	// Selenium implemented the actions API but has not yet updated its new session response.
	if !wd.w3cCompatible && !(wd.browser == "firefox" && wd.browserVersion.Major > 47) {
		return wd.KeyDown(keys)
	}
	return wd.keyAction("keyUp", keys)
}

func (wd *WebDriver) StoreKeyActions(inputID string, actions ...base.KeyAction) {
	var rawActions []map[string]interface{}
	for _, action := range actions {
		rawActions = append(rawActions, action)
	}
	wd.storedActions = append(wd.storedActions, map[string]interface{}{
		"type":    "key",
		"id":      inputID,
		"actions": rawActions,
	})
}

func (wd *WebDriver) StorePointerActions(inputID string, pointer base.PointerType, actions ...base.PointerAction) {
	var rawActions []map[string]interface{}
	for _, action := range actions {
		rawActions = append(rawActions, action)
	}
	wd.storedActions = append(wd.storedActions, map[string]interface{}{
		"type":       "pointer",
		"id":         inputID,
		"parameters": map[string]string{"pointerType": string(pointer)},
		"actions":    rawActions,
	})
}

func (wd *WebDriver) PerformActions() error {
	err := wd.voidCommand("/session/%s/actions", map[string]interface{}{
		"actions": wd.storedActions,
	})
	wd.storedActions = nil
	return err
}

func (wd *WebDriver) ReleaseActions() error {
	return wd.voidParamsCommand("DELETE", wd.requestURL("/session/%s/actions", wd.id), nil)
}

func (wd *WebDriver) DismissAlert() error {
	return wd.voidCommand("/session/%s/alert/dismiss", nil)
}

func (wd *WebDriver) AcceptAlert() error {
	return wd.voidCommand("/session/%s/alert/accept", nil)
}

func (wd *WebDriver) AlertText() (string, error) {
	return wd.stringCommand("/session/%s/alert/text")
}

func (wd *WebDriver) SetAlertText(text string) error {
	data := map[string]string{"text": text}
	return wd.voidCommand("/session/%s/alert/text", data)
}

func (wd *WebDriver) ExecuteScript(script string, args []interface{}) (interface{}, error) {
	if !wd.w3cCompatible {
		return wd.execScript(script, args, "")
	}
	return wd.execScript(script, args, "/sync")
}

func (wd *WebDriver) ExecuteScriptAsync(script string, args []interface{}) (interface{}, error) {
	if !wd.w3cCompatible {
		return wd.execScript(script, args, "_async")
	}
	return wd.execScript(script, args, "/async")
}

func (wd *WebDriver) ExecuteScriptRaw(script string, args []interface{}) ([]byte, error) {
	if !wd.w3cCompatible {
		return wd.execScriptRaw(script, args, "")
	}
	return wd.execScriptRaw(script, args, "/sync")
}

func (wd *WebDriver) ExecuteScriptAsyncRaw(script string, args []interface{}) ([]byte, error) {
	if !wd.w3cCompatible {
		return wd.execScriptRaw(script, args, "_async")
	}
	return wd.execScriptRaw(script, args, "/async")
}

func (wd *WebDriver) Screenshot() ([]byte, error) {
	data, err := wd.stringCommand("/session/%s/screenshot")
	if err != nil {
		return nil, err
	}
	buf := []byte(data)
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(buf))
	return ioutil.ReadAll(decoder)
}

func (wd *WebDriver) WaitWithTimeoutAndInterval(condition base.Condition, timeout, interval time.Duration) error {
	startTime := time.Now()
	for {
		done, err := condition(wd)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		if elapsed := time.Since(startTime); elapsed > timeout {
			return fmt.Errorf("timeout after %v", elapsed)
		}
		time.Sleep(interval)
	}
}

func (wd *WebDriver) WaitWithTimeout(condition base.Condition, timeout time.Duration) error {
	return wd.WaitWithTimeoutAndInterval(condition, timeout, DefaultWaitInterval)
}

func (wd *WebDriver) Wait(condition base.Condition) error {
	return wd.WaitWithTimeoutAndInterval(condition, DefaultWaitTimeout, DefaultWaitInterval)
}

func (wd *WebDriver) Log(kind base.LogType) ([]base.LogMessage, error) {
	rUrl := wd.requestURL("/session/%s/log", wd.id)
	params := map[string]base.LogType{"type": kind}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	response, err := wd.execute("POST", rUrl, data)
	if err != nil {
		return nil, err
	}

	c := new(struct {
		Value []struct {
			Timestamp int64
			Level     string
			Message   string
		}
	})
	if err = json.Unmarshal(response, c); err != nil {
		return nil, err
	}

	val := make([]base.LogMessage, len(c.Value))
	for i, v := range c.Value {
		val[i] = base.LogMessage{
			// n.b.: Chrome, which is the only browser that supports this API,
			// supplies timestamps in milliseconds since the Epoch.
			Timestamp: time.Unix(0, v.Timestamp*int64(time.Millisecond)),
			Level:     base.LogLevel(v.Level),
			Message:   v.Message,
		}
	}

	return val, nil
}

func (wd *WebDriver) requestURL(template string, args ...interface{}) string {
	return wd.urlPrefix + fmt.Sprintf(template, args...)
}

func (wd *WebDriver) newRequest(method string, url string, data []byte) (*http.Request, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", DefaultContentType)
	return request, nil
}

func (wd *WebDriver) execute(method, url string, data []byte) (json.RawMessage, error) {
	request, err := wd.newRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	response, err := wd.client.Do(request)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(response.Body)

	if wd.debug {
		fmt.Printf("%s %s %s --> %s [%s]\n", method, url, data, response.Status, response.Header["Content-Type"])
		if err != nil {
			fmt.Println(err.Error())
		} else {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, buf, "", "	"); err == nil {
				fmt.Println(pretty.String())
			} else {
				fmt.Println(string(buf))
			}
		}
	}

	if err != nil {
		return nil, errors.New(response.Status)
	}

	fullCType := response.Header.Get("Content-Type")
	cType, _, err := mime.ParseMediaType(fullCType)
	if err != nil {
		return nil, fmt.Errorf("got content type header %q, expected %q", fullCType, DefaultContentType)
	}
	if cType != DefaultContentType {
		return nil, fmt.Errorf("got content type %q, expected %q", cType, DefaultContentType)
	}

	reply := new(ServerReply)
	if err := json.Unmarshal(buf, reply); err != nil {
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("bad server reply status: %s", response.Status)
		}
		return nil, err
	}
	if reply.Err != "" {
		return nil, &reply.Error
	}

	// Handle the W3C-compliant error format. In the W3C spec, the error is embedded in the 'value' field.
	if len(reply.Value) > 0 {
		respErr := new(base.Error)
		if err := json.Unmarshal(reply.Value, respErr); err == nil && respErr.Err != "" {
			respErr.HTTPCode = response.StatusCode
			return nil, respErr
		}
	}

	// Handle the legacy error format.
	const success = 0
	if reply.Status != success {
		shortMsg, ok := base.RemoteErrors[reply.Status]
		if !ok {
			shortMsg = fmt.Sprintf("unknown error - %d", reply.Status)
		}

		longMsg := new(struct {
			Message string
		})
		if err := json.Unmarshal(reply.Value, longMsg); err != nil {
			return nil, errors.New(shortMsg)
		}
		return nil, &base.Error{
			Err:        shortMsg,
			Message:    longMsg.Message,
			HTTPCode:   response.StatusCode,
			LegacyCode: reply.Status,
		}
	}

	return buf, nil
}

func (wd *WebDriver) stringCommand(urlTemplate string) (string, error) {
	rUrl := wd.requestURL(urlTemplate, wd.id)
	response, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return "", err
	}

	reply := new(struct{ Value *string })
	if err := json.Unmarshal(response, reply); err != nil {
		return "", err
	}

	if reply.Value == nil {
		return "", fmt.Errorf("nil return value")
	}

	return *reply.Value, nil
}

func (wd *WebDriver) voidCommand(urlTemplate string, params interface{}) error {
	return wd.voidParamsCommand("POST", wd.requestURL(urlTemplate, wd.id), params)
}

func (wd *WebDriver) voidParamsCommand(method, url string, params interface{}) error {
	if params == nil {
		params = make(map[string]interface{})
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, err = wd.execute(method, url, data)
	return err
}

func (wd *WebDriver) stringsCommand(urlTemplate string) ([]string, error) {
	rUrl := wd.requestURL(urlTemplate, wd.id)
	response, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return nil, err
	}
	reply := new(struct{ Value []string })
	if err := json.Unmarshal(response, reply); err != nil {
		return nil, err
	}
	return reply.Value, nil
}

func (wd *WebDriver) boolCommand(urlTemplate string) (bool, error) {
	rUrl := wd.requestURL(urlTemplate, wd.id)
	response, err := wd.execute("GET", rUrl, nil)
	if err != nil {
		return false, err
	}
	reply := new(struct{ Value bool })
	if err := json.Unmarshal(response, reply); err != nil {
		return false, err
	}
	return reply.Value, nil
}

func (wd *WebDriver) find(by string, value string, suffix string, url string) ([]byte, error) {
	// The W3C specification removed the specific ID and Name locator strategies,
	// instead only providing a CSS-based strategy. Emulate the old behavior to
	// maintain API compatibility.
	if wd.w3cCompatible {
		switch by {
		case base.ByID:
			by = base.ByCSSSelector
			value = "#" + value
		case base.ByName:
			by = base.ByCSSSelector
			value = fmt.Sprintf("input[name=%q]", value)
		}
	}

	params := map[string]string{"using": by, "value": value}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if len(url) == 0 {
		url = "/session/%s/element"
	}

	return wd.execute("POST", wd.requestURL(url+suffix, wd.id), data)
}

func (wd *WebDriver) modifyWindow(handle string, verb string, command string, params interface{}) error {
	// The original protocol allowed for maximizing any named window. The W3C
	// specification only allows the current window be be modified. Emulate the
	// previous behavior by switching to the target window, maximizing the
	// current window, and switching back to the original window.
	var startWindow string
	if handle != "" && wd.w3cCompatible {
		var err error
		startWindow, err = wd.CurrentWindowHandle()
		if err != nil {
			return err
		}
		if handle != startWindow {
			if err := wd.SwitchWindow(handle); err != nil {
				return err
			}
		}
	}

	rUrl := wd.requestURL("/session/%s/window", wd.id)
	if command != "" {
		if wd.w3cCompatible {
			rUrl = wd.requestURL("/session/%s/window/%s", wd.id, command)
		} else {
			rUrl = wd.requestURL("/session/%s/window/%s/%s", wd.id, handle, command)
		}
	}

	var data []byte
	if params != nil {
		var err error
		if data, err = json.Marshal(params); err != nil {
			return err
		}
	}

	if _, err := wd.execute(verb, rUrl, data); err != nil {
		return err
	}

	// TODO: add a test for switching back to the original window.
	if handle != startWindow && wd.w3cCompatible {
		if err := wd.SwitchWindow(startWindow); err != nil {
			return err
		}
	}

	return nil
}

func (wd *WebDriver) keyAction(action, keys string) error {
	type keyAction struct {
		Type string `json:"type"`
		Key  string `json:"value"`
	}
	actions := make([]keyAction, 0, len(keys))
	for _, key := range keys {
		actions = append(actions, keyAction{
			Type: action,
			Key:  string(key),
		})
	}
	return wd.voidCommand("/session/%s/actions", map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"type":    "key",
				"id":      "default keyboard",
				"actions": actions,
			}},
	})
}

func (wd *WebDriver) execScriptRaw(script string, args []interface{}, suffix string) ([]byte, error) {
	if args == nil {
		args = make([]interface{}, 0)
	}
	data, err := json.Marshal(map[string]interface{}{
		"script": script,
		"args":   args,
	})
	if err != nil {
		return nil, err
	}
	return wd.execute("POST", wd.requestURL("/session/%s/execute"+suffix, wd.id), data)
}

func (wd *WebDriver) execScript(script string, args []interface{}, suffix string) (interface{}, error) {
	response, err := wd.execScriptRaw(script, args, suffix)
	if err != nil {
		return nil, err
	}
	reply := new(struct{ Value interface{} })
	if err = json.Unmarshal(response, reply); err != nil {
		return nil, err
	}
	return reply.Value, nil
}
