package wd

import (
	"github.com/markel1974/webautoma/src/wd/base"
	"github.com/markel1974/webautoma/src/wd/caps/firefox"
	"strings"
)

// The list of valid, top-level capability names, according to the W3C specification.
//
// This must be kept in sync with the specification:
// https://www.w3.org/TR/webdriver/#capabilities
var w3cCapabilityNames = []string{
	"acceptInsecureCerts",
	"browserName",
	"browserVersion",
	"platformName",
	"pageLoadStrategy",
	"proxy",
	"setWindowRect",
	"timeouts",
	"unhandledPromptBehavior",
}

var chromeCapabilityNames = []string{
	// This is not a standardized top-level capability name, but Chromedriver
	// expects this capability here.
	// https://cs.chromium.org/chromium/src/chrome/test/chromedriver/capabilities.cc?rcl=0754b5d0aad903439a628618f0e41845f1988f0c&l=759
	"loggingPrefs",
}

// Create a W3C-compatible capabilities instance.
func newW3CCapabilities(cs base.Capabilities) base.Capabilities {
	isValidW3CCapability := map[string]bool{}
	for _, name := range w3cCapabilityNames {
		isValidW3CCapability[name] = true
	}
	if b, ok := cs["browserName"]; ok && b == "chrome" {
		for _, name := range chromeCapabilityNames {
			isValidW3CCapability[name] = true
		}
	}

	alwaysMatch := make(base.Capabilities)
	for name, value := range cs {
		if isValidW3CCapability[name] || strings.Contains(name, ":") {
			alwaysMatch[name] = value
		}
	}

	// Move the Firefox profile setting from the old location to the new
	// location.
	if prof, ok := cs["firefox_profile"]; ok {
		if c, ok := alwaysMatch[firefox.CapabilitiesKey]; ok {
			firefoxCaps := c.(firefox.Caps)
			if firefoxCaps.Profile == "" {
				firefoxCaps.Profile = prof.(string)
			}
		} else {
			alwaysMatch[firefox.CapabilitiesKey] = firefox.Caps{
				Profile: prof.(string),
			}
		}
	}

	return base.Capabilities{
		"alwaysMatch": alwaysMatch,
	}
}
