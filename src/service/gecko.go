package service

import (
	"net/url"
)

type GeckoService struct {
	*Service
}

func NewGeckoService(nohup bool, path string, wdUrl *url.URL) (*GeckoService, error) {
	var args []string
	args = append(args, "--port="+wdUrl.Port())

	s, err := NewService(nohup, path, args, wdUrl, "")
	if err != nil {
		return nil, err
	}

	gs := &GeckoService{
		Service: s,
	}

	return gs, nil
}
