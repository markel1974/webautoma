package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ServiceOption func(*Service) error

func NewSeleniumService(jarPath string, port int, opts ...ServiceOption) (*Service, error) {
	s, err := newService(exec.Command("java"), "/wd/hub", port, opts...)
	if err != nil {
		return nil, err
	}
	if s.javaPath != "" {
		s.cmd.Path = s.javaPath
	}
	if s.geckoDriverPath != "" {
		s.cmd.Args = append([]string{"java", "-Dwebdriver.gecko.driver=" + s.geckoDriverPath}, s.cmd.Args[1:]...)
	}
	if s.chromeDriverPath != "" {
		s.cmd.Args = append([]string{"java", "-Dwebdriver.chrome.driver=" + s.chromeDriverPath}, s.cmd.Args[1:]...)
	}

	var classpath []string
	if s.htmlUnitPath != "" {
		classpath = append(classpath, s.htmlUnitPath)
	}
	classpath = append(classpath, jarPath)
	s.cmd.Args = append(s.cmd.Args, "-cp", strings.Join(classpath, ":"))
	s.cmd.Args = append(s.cmd.Args, "org.openqa.grid.selenium.GridLauncherV3", "-port", strconv.Itoa(port), "-debug")

	if err := s.start(port); err != nil {
		return nil, err
	}
	return s, nil
}

func NewChromeDriverService(path string, port int, opts ...ServiceOption) (*Service, error) {
	cmd := exec.Command(path, "--port="+strconv.Itoa(port), "--url-base=wd/hub", "--verbose")
	s, err := newService(cmd, "/wd/hub", port, opts...)
	if err != nil {
		return nil, err
	}
	s.shutdownURLPath = "/shutdown"
	if err := s.start(port); err != nil {
		return nil, err
	}
	return s, nil
}

func NewGeckoDriverService(path string, port int, opts ...ServiceOption) (*Service, error) {
	cmd := exec.Command(path, "--port", strconv.Itoa(port))
	s, err := newService(cmd, "", port, opts...)
	if err != nil {
		return nil, err
	}
	if err := s.start(port); err != nil {
		return nil, err
	}
	return s, nil
}

// Service controls a locally-running Selenium subprocess.
type Service struct {
	port                      int
	addr                      string
	cmd                       *exec.Cmd
	shutdownURLPath           string
	display                   string
	xAuthPath                 string
	xvfb                      *FrameBuffer
	geckoDriverPath, javaPath string
	chromeDriverPath          string
	htmlUnitPath              string
	output                    io.Writer
}

func newService(cmd *exec.Cmd, urlPrefix string, port int, opts ...ServiceOption) (*Service, error) {
	s := &Service{
		port: port,
		addr: fmt.Sprintf("http://localhost:%d%s", port, urlPrefix),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	cmd.Stderr = s.output
	cmd.Stdout = s.output
	cmd.Env = os.Environ()
	// TODO(minusnine): Pdeathsig is only supported on Linux. Somehow, make sure
	// process cleanup happens as gracefully as possible.
	if s.display != "" {
		cmd.Env = append(cmd.Env, "DISPLAY=:"+s.display)
	}
	if s.xAuthPath != "" {
		cmd.Env = append(cmd.Env, "XAUTHORITY="+s.xAuthPath)
	}
	s.cmd = cmd
	return s, nil
}

func (s *Service) start(port int) error {
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
	return fmt.Errorf("server did not respond on port %d", port)
}

func (s *Service) SetOutput(w io.Writer) {
	s.output = w
}

func (s *Service) SetGeckoDriver(path string) {
	s.geckoDriverPath = path
}

func (s *Service) SetChromeDriver(path string) {
	s.chromeDriverPath = path
}

func (s *Service) SetJavaPath(path string) {
	s.javaPath = path
}

func (s *Service) SetHTMLUnit(path string) {
	s.htmlUnitPath = path
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
	if s.xvfb != nil {
		return s.xvfb.Stop()
	}
	return nil
}

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
