package base

import "time"

// LogType represents a component capable of logging.
type LogType string

// The valid log types.
const (
	LogServer      LogType = "server"
	LogBrowser     LogType = "browser"
	LogClient      LogType = "client"
	LogDriver      LogType = "driver"
	LogPerformance LogType = "performance"
	LogProfiler    LogType = "profiler"
)

// LogLevel represents a logging level of different components in the browser,
// the driver, or any intermediary WebDriver servers.
//
// See the documentation of each driver for what browser specific logging
// components are available.
type LogLevel string

// The valid log levels.
const (
	LogOff     LogLevel = "OFF"
	LogSevere  LogLevel = "SEVERE"
	LogWarning LogLevel = "WARNING"
	LogInfo    LogLevel = "INFO"
	LogDebug   LogLevel = "DEBUG"
	LogAll     LogLevel = "ALL"
)

// LogCapabilitiesKey is the key for the logging preferences entry in the JSON
// structure representing WebDriver capabilities.
//
// Note that the W3C spec does not include logging right now, and starting with
// Chrome 75, "loggingPrefs" has been changed to "goog:loggingPrefs"
const LogCapabilitiesKey = "goog:loggingPrefs"

// LogActions is the map to include in the WebDriver capabilities structure
// to configure logging.
type LogActions map[LogType]LogLevel

// LogMessage is a log message returned from the Log method.
type LogMessage struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}
