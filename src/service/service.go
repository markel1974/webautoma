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
	port            int
	addr            string
	cmd             *exec.Cmd
	shutdownURLPath string
	//display                   string
	//xAuthPath                 string
	//xvfb                      *FrameBuffer
	output io.Writer
}

func NewService(cmd *exec.Cmd, urlPrefix string, port int, shutdownURLPath string) (*Service, error) {
	s := &Service{
		port:            port,
		addr:            fmt.Sprintf("http://localhost:%d/%s", port, urlPrefix),
		shutdownURLPath: shutdownURLPath,
	}

	cmd.Env = os.Environ()

	/*
		// TODO: Pdeathsig is only supported on Linux. Somehow, make sure
		// process cleanup happens as gracefully as possible.
		if s.display != "" {
			cmd.Env = append(cmd.Env, "DISPLAY=:"+s.display)
		}
		if s.xAuthPath != "" {
			cmd.Env = append(cmd.Env, "XAUTHORITY="+s.xAuthPath)
		}*/
	s.cmd = cmd
	return s, nil
}

func (s *Service) SetOutput(w io.Writer) {
	s.cmd.Stderr = w
	s.cmd.Stdout = w
}

func (s *Service) Start() error {
	if err := s.cmd.Start(); err != nil {
		return err
	}

	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		resp, err := http.Get(s.addr + "/status")
		if err == nil {
			_ = resp.Body.Close()
			switch resp.StatusCode {
			// Selenium <3 returned Forbidden and BadRequest. ChromeDriver and
			// Selenium 3 return OK.
			case http.StatusForbidden, http.StatusBadRequest, http.StatusOK:
				return nil
			}
		}
	}
	return fmt.Errorf("server did not respond on port %d", s.port)
}

func (s *Service) Stop() error {
	// Selenium 3 stopped supporting the shutdown URL by default.
	// https://github.com/SeleniumHQ/selenium/issues/2852

	if s.shutdownURLPath == "" {
		if err := s.cmd.Process.Kill(); err != nil {
			return err
		}
	} else {
		resp, err := http.Get(s.addr + s.shutdownURLPath)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
	}

	if err := s.cmd.Wait(); err != nil && err.Error() != "signal: killed" {
		return err
	}
	/*
		if s.xvfb != nil {
			return s.xvfb.Stop()
		}
	*/
	return nil
}

/*
func (s *Service) FrameBuffer() *FrameBuffer {
	return s.xvfb
}

func (s *Service) StartFrameBuffer() error {
	return s.StartFrameBufferWithOptions(FrameBufferOptions{})
}

func (s *Service) StartFrameBufferWithOptions(options FrameBufferOptions) error {
	if s.display != "" {
		return fmt.Errorf("service display already set: %v", s.display)
	}
	if s.xAuthPath != "" {
		return fmt.Errorf("service xauth path already set: %v", s.xAuthPath)
	}
	if s.xvfb != nil {
		return fmt.Errorf("service Xvfb instance already running")
	}
	fb, err := NewFrameBufferWithOptions(options)
	if err != nil {
		return fmt.Errorf("error starting frame buffer: %v", err)
	}
	s.xvfb = fb
	return s.Display(fb.Display, fb.AuthPath)
}

func (s *Service) Display(target string, xAuthPath string) error {
	if s.display != "" {
		return fmt.Errorf("service display already set: %v", s.display)
	}
	if s.xAuthPath != "" {
		return fmt.Errorf("service xauth path already set: %v", s.xAuthPath)
	}
	if !s.isDisplay(target) {
		return fmt.Errorf("supplied display %q must be of the format 'x' or 'x.y' where x and y are integers", target)
	}
	s.display = target
	s.xAuthPath = xAuthPath
	return nil
}

func (s *Service) isDisplay(target string) bool {
	ds := strings.Split(target, ".")
	if len(ds) > 2 {
		return false
	}
	for _, d := range ds {
		if _, err := strconv.Atoi(d); err != nil {
			return false
		}
	}
	return true
}
*/
