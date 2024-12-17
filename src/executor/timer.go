package executor

import "time"

// Timer represents a structure for tracking time intervals and recording elapsed durations across multiple runs.
// It includes properties for identification, description, start time, epoch timestamp, finalization status, and data storage.
// The Id and Desc fields identify and describe the Timer, respectively.
// The Start field tracks the start time of a timing operation.
// The Epoch stores the timestamp in milliseconds for the latest start time.
// The Finalized field indicates whether the Timer has been concluded, preventing further modification.
// The Data field captures the durations of each timing segment recorded by the Timer.
type Timer struct {
	Id        string
	Desc      string
	Start     time.Time
	Epoch     int64
	Finalized bool
	Data      []int64
}

// UnixMilli converts a time.Time object to a Unix timestamp in milliseconds.
func UnixMilli(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// NewTimer creates and returns a new Timer instance with the specified id and description.
func NewTimer(id string, desc string) *Timer {
	return &Timer{
		Id:        id,
		Desc:      desc,
		Epoch:     0,
		Finalized: false,
		Data:      nil,
	}
}

// Run initializes the Timer by setting the Start time to the current time and updating the Epoch value in milliseconds.
func (t *Timer) Run() {
	t.Start = time.Now()
	t.Epoch = UnixMilli(t.Start)
}

// Stop records the elapsed time since the last epoch, updates the epoch, and appends the duration to the Timer's data slice.
func (t *Timer) Stop() {
	var epoch = UnixMilli(time.Now())
	var diff = epoch - t.Epoch
	t.Epoch = epoch
	t.Data = append(t.Data, diff)
}

// Finalize aggregates the durations in the timer, clears its state, marks it as finalized, and returns the total duration.
func (t *Timer) Finalize() int64 {
	if t.Finalized {
		return 0
	}
	var dur = int64(0)
	for _, entry := range t.Data {
		dur += entry
	}
	t.Data = nil
	t.Epoch = 0
	t.Finalized = true
	return dur
}
