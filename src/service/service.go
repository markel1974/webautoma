package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Service struct {
	port        int
	url         string
	cmd         *exec.Cmd
	shutdownURL string
	screenSize  string
	fb          *FrameBuffer
	output      io.Writer
	display     bool
}

func NewService(nohup bool, wdPath string, wdArgs []string, wdUrl *url.URL, shutdownURL string) (*Service, error) {
	if nohup {
		wdArgs = append([]string{wdPath}, wdArgs...)
		wdPath = "nohup"
	}

	port, _ := strconv.Atoi(wdUrl.Port())
	s := &Service{
		port:        port,
		url:         wdUrl.String(),
		shutdownURL: shutdownURL,
		cmd:         exec.Command(wdPath, wdArgs...),
		display:     false,
	}
	s.cmd.Env = os.Environ()
	return s, nil
}

func (s *Service) SetOutput(w io.Writer, e io.Writer) {
	s.cmd.Stdout = w
	s.cmd.Stderr = e
}

func (s *Service) EnableDisplay(screenSize string) {
	s.display = true
	s.screenSize = screenSize
}

func (s *Service) Display() (string, string) {
	if s.fb == nil {
		return "", ""
	}
	return s.fb.Display()
}

func (s *Service) Start() error {
	if s.display {
		s.fb = NewFrameBuffer(s.screenSize)
		if err := s.fb.Start(); err != nil {
			return err
		}
	}
	if err := s.cmd.Start(); err != nil {
		return err
	}
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		resp, err := http.Get(s.url + "/status")
		if err == nil {
			_ = resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusForbidden, http.StatusBadRequest, http.StatusOK:
				return nil
			}
		}
	}
	return fmt.Errorf("server did not respond on port %d", s.port)
}

func (s *Service) Stop() error {
	if len(s.shutdownURL) == 0 {
		if err := s.cmd.Process.Kill(); err != nil {
			return err
		}
	} else {
		resp, err := http.Get(s.url + s.shutdownURL)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
	}

	if err := s.cmd.Wait(); err != nil && err.Error() != "signal: killed" {
		return err
	}
	if s.fb != nil {
		return s.fb.Stop()
	}
	return nil
}

/*
func (s * Service) SetDisplay(display string, xAuthPath string, screenSize string) {
	if len(display) > 0 {
		s.cmd.Env = append(s.cmd.Env, "DISPLAY=:"+ display)
	}
	if len(xAuthPath) > 0 {
		s.cmd.Env = append(s.cmd.Env, "XAUTHORITY="+ xAuthPath)
	}
	if len(screenSize) > 0 {
		s.screenSize = screenSize
	}
}
*/
