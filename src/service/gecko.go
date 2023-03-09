package service

import (
	"strconv"
)

type GeckoService struct {
	*Service
}

func NewGeckoService(path string, port int) (*GeckoService, error) {
	var args []string
	args = append(args, "--port="+strconv.Itoa(port))

	s, err := NewService(path, args, "", port, "")
	if err != nil {
		return nil, err
	}

	gs := &GeckoService{
		Service: s,
	}

	return gs, nil
}
