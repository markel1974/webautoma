package service

import (
	"io"
	"os/exec"
	"strconv"
)

type ChromeService struct {
	svc *Service
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
	cmd := exec.Command(path, args...)
	svc, err := NewService(cmd, base, port, "/shutdown")
	if err != nil {
		return nil, err
	}
	cs := &ChromeService{
		svc: svc,
	}
	return cs, nil
}

func (cs *ChromeService) SetOutput(w io.Writer) {
	cs.svc.SetOutput(w)
}

func (cs *ChromeService) Start() error {
	if err := cs.svc.Start(); err != nil {
		return err
	}
	return nil
}

func (cs *ChromeService) Stop() error {
	if err := cs.svc.Stop(); err != nil {
		return err
	}
	return nil
}
