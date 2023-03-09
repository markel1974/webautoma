package base

import (
	"github.com/markel1974/webautoma/src/wd/caps/chrome"
	"github.com/markel1974/webautoma/src/wd/caps/firefox"
)

type Capabilities map[string]interface{}

func (c Capabilities) SetChrome(f chrome.Caps) {
	c[chrome.CapabilitiesKey] = f
	c[chrome.DeprecatedCapabilitiesKey] = f
}

func (c Capabilities) SetFirefox(f firefox.Caps) {
	c[firefox.CapabilitiesKey] = f
}

func (c Capabilities) SetProxy(p Proxy) {
	c["proxy"] = p
}

func (c Capabilities) SetLogging(l LogActions) {
	c[LogCapabilitiesKey] = l
}

func (c Capabilities) SetLogLevel(kind LogType, level LogLevel) {
	if _, ok := c[LogCapabilitiesKey]; !ok {
		c[LogCapabilitiesKey] = make(LogActions)
	}
	m := c[LogCapabilitiesKey].(LogActions)
	m[kind] = level
}
