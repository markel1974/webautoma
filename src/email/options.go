package email

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Options struct {
	subjectRgx        *regexp.Regexp
	bodyRgx           *regexp.Regexp
	verifyIntervalSec int
	validityMin       int
}

func NewOptions(value string) (*Options, error) {
	t := &Options{
		subjectRgx:        nil,
		bodyRgx:           nil,
		verifyIntervalSec: 60,
		validityMin:       5,
	}
	opt := strings.Split(value, "|||")
	if len(opt) < 2 {
		return nil, fmt.Errorf("invalid value, missing separator")
	}
	var err error
	t.subjectRgx, err = regexp.Compile(opt[0])
	if err != nil {
		return nil, fmt.Errorf("invalid subject regexp, %s", err.Error())
	}
	t.bodyRgx, err = regexp.Compile(opt[1])
	if err != nil {
		return nil, fmt.Errorf("invalid body regexp, %s", err.Error())
	}
	if len(opt) > 2 {
		t.verifyIntervalSec, err = strconv.Atoi(opt[2])
		if err != nil {
			return nil, fmt.Errorf("invalid interval, %s", err.Error())
		}
	}
	if len(opt) > 3 {
		t.validityMin, err = strconv.Atoi(opt[3])
		if err != nil {
			return nil, fmt.Errorf("invalid interval, %s", err.Error())
		}
	}
	return t, nil
}

func (o *Options) VerifyIntervalSec() int64 {
	return int64(o.verifyIntervalSec)
}

func (o *Options) ValidityMin() int64 {
	return int64(o.validityMin)
}
