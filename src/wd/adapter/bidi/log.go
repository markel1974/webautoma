package bidi

import (
	"encoding/json"
	"math"
	"time"
)

type LogType string

const (
	LogTypeText       LogType = "text"
	LogTypeConsole    LogType = "console"
	LogTypeJavascript LogType = "javascript"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type baseLogEntry struct {
	Type       LogType    `json:"type"`
	Level      LogLevel   `json:"level"`
	Source     Source     `json:"source"`
	Text       string     `json:"text"`
	Timestamp  Timestamp  `json:"timestamp"`
	StackTrace StackTrace `json:"stackTrace"`
}

type GenericLogEntry struct {
	baseLogEntry
}

type ConsoleLogEntry struct {
	baseLogEntry
	Method string      `json:"method"`
	Args   interface{} `json:"args"`
}

type JavascriptLogEntry struct {
	baseLogEntry
}

type OnLogEntryHandler struct {
	LogTypeTextHandlerFunc       func(entry *GenericLogEntry) error
	LogTypeConsoleHandlerFunc    func(entry *ConsoleLogEntry) error
	LogTypeJavascriptHandlerFunc func(entry *JavascriptLogEntry) error
}

func (s *Session) OnLogEntry(handler *OnLogEntryHandler) {
	s.client.CallbackEvent("log.entryAdded", func(params json.RawMessage) error {
		type entry struct {
			Type LogType `json:"type"`
		}

		e := &entry{}
		if err := json.Unmarshal(params, &e); err != nil {
			return err
		}

		switch e.Type {
		case LogTypeText:
			if handler.LogTypeTextHandlerFunc != nil {
				e := &GenericLogEntry{}
				if err := json.Unmarshal(params, &e); err != nil {
					return err
				}
				return handler.LogTypeTextHandlerFunc(e)
			}
		case LogTypeConsole:
			if handler.LogTypeConsoleHandlerFunc != nil {
				e := &ConsoleLogEntry{}
				if err := json.Unmarshal(params, &e); err != nil {
					return err
				}
				return handler.LogTypeConsoleHandlerFunc(e)
			}
		case LogTypeJavascript:
			if handler.LogTypeJavascriptHandlerFunc != nil {
				e := &JavascriptLogEntry{}
				if err := json.Unmarshal(params, &e); err != nil {
					return err
				}
				return handler.LogTypeJavascriptHandlerFunc(e)
			}
		}

		return nil
	})
}

type Timestamp struct {
	time.Time
}

// UnmarshalJSON decodes an int64 timestamp into a time.Time object
func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	var raw float64
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	//TODO
	sec, dec := math.Modf(raw)
	p.Time = time.Unix(int64(sec), int64(dec*(1e9)))

	return nil
}
