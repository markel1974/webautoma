package service

import (
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type SeleniumService struct {
	svc *Service
}

func NewSeleniumService(port int, base string, jarPath string, java string, gecko string, chrome string, htmlUnit string, debug bool) (*SeleniumService, error) {
	path := "java"
	if len(java) > 0 {
		path = java
	}
	var args []string
	if len(gecko) > 0 {
		args = append(args, "-Dwebdriver.gecko.driver="+gecko)
	}
	if chrome != "" {
		args = append(args, "-Dwebdriver.chrome.driver="+chrome)
	}
	var classpath []string
	if htmlUnit != "" {
		classpath = append(classpath, htmlUnit)
	}
	if len(jarPath) > 0 {
		classpath = append(classpath, jarPath)
	}
	args = append(args, "-cp", strings.Join(classpath, ":"))
	args = append(args, "org.openqa.grid.selenium.GridLauncherV3")
	args = append(args, "-port", strconv.Itoa(port))
	if debug {
		args = append(args, "-debug")
	}
	cmd := exec.Command(path, args...)
	svc, err := NewService(cmd, base, port, "")
	if err != nil {
		return nil, err
	}
	ss := &SeleniumService{svc: svc}
	return ss, nil
}

func (ss *SeleniumService) SetOutput(w io.Writer) {
	ss.svc.SetOutput(w)
}

func (ss *SeleniumService) SetDisplay(screenSize string) {
	ss.svc.SetDisplay(screenSize)
}

func (ss *SeleniumService) Display() (string, string) {
	return ss.svc.Display()
}

func (ss *SeleniumService) Start() error {
	if err := ss.svc.Start(); err != nil {
		return err
	}
	return nil
}

func (ss *SeleniumService) Stop() error {
	if err := ss.svc.Stop(); err != nil {
		return err
	}
	return nil
}
