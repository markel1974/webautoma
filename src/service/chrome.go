package service

import (
	"strconv"
)

type ChromeService struct {
	*Service
}

func NewChromeService(path string, port int, base string, verbose bool) (*ChromeService, error) {
	var args []string
	args = append(args, "--port="+strconv.Itoa(port))
	if len(base) > 0 {
		args = append(args, "--url-base="+base)
	}
	if verbose {
		args = append(args, "--verbose")
	}
	svc, err := NewService(path, args, base, port, "/shutdown")
	if err != nil {
		return nil, err
	}
	cs := &ChromeService{
		Service: svc,
	}
	return cs, nil
}
