package service

import (
	"net/url"
)

type ChromeService struct {
	*Service
}

func NewChromeService(nohup bool, wdPath string, wdUrl *url.URL, verbose bool) (*ChromeService, error) {
	var wdArgs []string
	wdArgs = append(wdArgs, "--port="+wdUrl.Port())
	if len(wdUrl.Path) > 0 {
		wdArgs = append(wdArgs, "--url-base="+wdUrl.Path)
	}
	if verbose {
		wdArgs = append(wdArgs, "--verbose")
	}
	svc, err := NewService(nohup, wdPath, wdArgs, wdUrl, "/shutdown")
	if err != nil {
		return nil, err
	}
	cs := &ChromeService{
		Service: svc,
	}
	return cs, nil
}
