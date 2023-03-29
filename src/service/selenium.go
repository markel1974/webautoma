package service

import (
	"net/url"
	"strings"
)

type SeleniumService struct {
	*Service
}

func NewSeleniumService(nohup bool, wdUrl *url.URL, jarPath string, java string, gecko string, chrome string, htmlUnit string, debug bool) (*SeleniumService, error) {
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
	if len(classpath) > 0 {
		args = append(args, "-cp", strings.Join(classpath, ":"))
	}
	args = append(args, "org.openqa.grid.selenium.GridLauncherV3")
	args = append(args, "-port", wdUrl.Port())
	if debug {
		args = append(args, "-debug")
	}
	svc, err := NewService(nohup, path, args, wdUrl, "")
	if err != nil {
		return nil, err
	}
	ss := &SeleniumService{
		Service: svc,
	}
	return ss, nil
}
