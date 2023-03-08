package executor

import "time"

type Timer struct {
	Id        string
	Desc      string
	Start     time.Time
	Epoch     int64
	Finalized bool
	Data      []int64
}

func NewTimer(id string, desc string) *Timer {
	return &Timer{
		Id:        id,
		Desc:      desc,
		Epoch:     0,
		Finalized: false,
		Data:      nil,
	}
}

func (t *Timer) Run() {
	t.Start = time.Now()
	t.Epoch = t.Start.UnixMilli()
}

func (t *Timer) Stop() {
	var epoch = time.Now().UnixMilli()
	var diff = epoch - t.Epoch
	t.Epoch = epoch
	t.Data = append(t.Data, diff)
}

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
