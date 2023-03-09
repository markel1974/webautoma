package service

import (
	"io"
	"os/exec"
	"strconv"
)

type GeckoService struct {
	svc *Service
}

func NewGeckoService(path string, port int) (*GeckoService, error) {
	var args []string
	args = append(args, "--port="+strconv.Itoa(port))

	cmd := exec.Command(path, args...)
	s, err := NewService(cmd, "", port, "")
	if err != nil {
		return nil, err
	}

	gs := &GeckoService{svc: s}

	return gs, nil
}

func (cs *GeckoService) SetOutput(w io.Writer) {
	cs.svc.SetOutput(w)
}

func (gs *GeckoService) Start() error {
	if err := gs.svc.Start(); err != nil {
		return err
	}
	return nil
}

func (gs *GeckoService) Stop() error {
	if err := gs.svc.Stop(); err != nil {
		return err
	}
	return nil
}
