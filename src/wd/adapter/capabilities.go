package adapter

type Capabilities map[string]interface{}

// SetBrowserName sets the desired browser name.
func (c Capabilities) SetBrowserName(name string) Capabilities {
	c["browserName"] = name
	return c
}

// A ProxyConfig instance defines the desired proxy configuration the WebDriver
// should use to proxy a Page.
//
// See: https://www.w3.org/TR/webdriver/#proxy
type ProxyConfig struct {
	ProxyType          string `json:"proxyType"`
	ProxyAutoconfigURL string `json:"proxyAutoconfigUrl,omitempty"`
	FTPProxy           string `json:"ftpProxy,omitempty"`
	HTTPProxy          string `json:"httpProxy,omitempty"`
	SSLProxy           string `json:"sslProxy,omitempty"`
	SOCKSProxy         string `json:"socksProxy,omitempty"`
	SOCKSUsername      string `json:"socksUsername,omitempty"`
	SOCKSPassword      string `json:"socksPassword,omitempty"`
	NoProxy            string `json:"noProxy,omitempty"`
}

// SetProxy Proxy sets the desired proxy configuration.
func (c Capabilities) SetProxy(p ProxyConfig) Capabilities {
	c["proxy"] = p
	return c
}

// SetBrowserVersion sets the desired browser version.
func (c Capabilities) SetBrowserVersion(version string) Capabilities {
	c["browserVersion"] = version
	return c
}

// SetPlatformName sets the desired browser platform.
func (c Capabilities) SetPlatformName(platform string) Capabilities {
	c["platformName"] = platform
	return c
}

func (c Capabilities) SetAcceptInsecureCerts(acceptInsecureCerts bool) Capabilities {
	c["acceptInsecureCerts"] = acceptInsecureCerts
	return c
}

func (c Capabilities) SetWebSocketURL(webSocketURL bool) Capabilities {
	c["webSocketUrl"] = webSocketURL
	return c
}

func (c Capabilities) WebSocketURL() string {
	if val, ok := c["webSocketUrl"]; ok {
		return val.(string)
	}

	return ""
}

// Set Sets an arbitrary key-value pair
func (c Capabilities) Set(key string, value string) Capabilities {
	c[key] = value
	return c
}

type ChromeOptions map[string]interface{}

func (c Capabilities) SetChromeOptions(co ChromeOptions) Capabilities {
	c["goog:chromeOptions"] = co
	return c
}

func (c Capabilities) ChromeOptions() ChromeOptions {
	if opts, ok := c["goog:chromeOptions"]; ok {
		return opts.(map[string]interface{})
	}
	return nil
}

func (co ChromeOptions) AddArg(arg string) ChromeOptions {
	if _, ok := co["args"]; ok {
		co["args"] = append(co["args"].([]string), arg)
	} else {
		co["args"] = []string{arg}
	}
	return co
}

func (co ChromeOptions) SetBinary(binary string) ChromeOptions {
	co["binary"] = binary
	return co
}

func (co ChromeOptions) Binary() string {
	if val, ok := co["binary"]; ok {
		return val.(string)
	}
	return ""
}

func (co ChromeOptions) DebuggerAddress() string {
	if val, ok := co["debuggerAddress"]; ok {
		return val.(string)
	}
	return ""
}
