package executor

import "time"

type Event struct {
	ProbeId           string                 `json:"probeId"`
	ThreadName        string                 `json:"thread_name"`
	TransactionType   string                 `json:"transaction_type"`
	Timestamp         string                 `json:"timestamp"`
	Start             string                 `json:"start"`
	Stop              string                 `json:"stop"`
	ExecutionTime     int64                  `json:"execution_time"`
	UUID              string                 `json:"uuid"`
	Level             string                 `json:"livello"`
	Error             string                 `json:"error"`
	Passed            bool                   `json:"passed"`
	NetworkErrorCount int                    `json:"networkErrorCount"`
	Network           map[string]interface{} `json:"network"`
	Acquired          map[string]interface{} `json:"acquired,omitempty"`
	ScreenShoot       string                 `json:"screenshotId,omitempty"`
	Label             string                 `json:"label,omitempty"`
	Target            string                 `json:"target,omitempty"`
	Command           string                 `json:"command,omitempty"`
	Message           string                 `json:"message,omitempty"`
}

func NewEvent(id string, kind string, err error, start time.Time, dur int64) *Event {
	stop := time.Now()
	errorDesc := ""
	if err != nil {
		errorDesc = err.Error()
	}
	event := &Event{
		ProbeId:           "",
		ThreadName:        id,
		TransactionType:   kind,
		Timestamp:         time.Now().Format(RFC3339Milli),
		Start:             start.Format(RFC3339Milli),
		Stop:              stop.Format(RFC3339Milli),
		ExecutionTime:     dur,
		UUID:              "",
		Level:             "INFO",
		Error:             errorDesc,
		Passed:            err == nil,
		Network:           nil,
		NetworkErrorCount: 0,
		ScreenShoot:       "",
		Label:             "",
		Target:            "",
		Command:           "",
		Message:           "",
	}
	return event
}
