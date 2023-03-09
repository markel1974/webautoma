package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Service struct {
	port        int
	addr        string
	cmd         *exec.Cmd
	shutdownURL string
	screenSize  string
	fb          *FrameBuffer
	output      io.Writer
}

func NewService(path string, args []string, urlPrefix string, port int, shutdownURL string) (*Service, error) {
	s := &Service{
		port:        port,
		addr:        fmt.Sprintf("http://localhost:%d/%s", port, urlPrefix),
		shutdownURL: shutdownURL,
		cmd:         exec.Command(path, args...),
	}
	s.cmd.Env = os.Environ()
	return s, nil
}

func (s *Service) SetOutput(w io.Writer, e io.Writer) {
	s.cmd.Stdout = w
	s.cmd.Stderr = e
}

func (s *Service) SetDisplay(screenSize string) {
	s.screenSize = screenSize
}

func (s *Service) Display() (string, string) {
	if s.fb == nil {
		return "", ""
	}
	return s.fb.Display()
}

func (s *Service) Start() error {
	if len(s.screenSize) > 0 {
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
		resp, err := http.Get(s.addr + "/status")
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
		resp, err := http.Get(s.addr + s.shutdownURL)
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
