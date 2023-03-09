package service

import (
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
	classpath = append(classpath, jarPath)
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
	sln := &SeleniumService{svc: svc}
	return sln, nil
}
