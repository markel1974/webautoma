package base

import "fmt"

var RemoteErrors = map[int]string{
	6:  "invalid session ID",
	7:  "no such element",
	8:  "no such frame",
	9:  "unknown command",
	10: "stale element reference",
	11: "element not visible",
	12: "invalid element state",
	13: "unknown error",
	15: "element is not selectable",
	17: "javascript error",
	19: "xpath lookup error",
	21: "timeout",
	23: "no such window",
	24: "invalid cookie domain",
	25: "unable to set cookie",
	26: "unexpected alert open",
	27: "no alert open",
	28: "script timeout",
	29: "invalid element coordinates",
	32: "invalid selector",
}

// see https://www.w3.org/TR/webdriver/#handling-errors .
type Error struct {
	Err        string `json:"error"`
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace"`
	HTTPCode   int
	LegacyCode int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Message)
}
