package wd

import (
	"crypto/tls"
	"net/http"
)

func NewDefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//IMPORTANT chromedriver support only 6 connections
			MaxConnsPerHost:     4,
			MaxIdleConnsPerHost: 2,
			MaxIdleConns:        2,
			//MaxIdleConns:          10,
			//IdleConnTimeout:       3600 * time.Second,
			//TLSHandshakeTimeout:   3600 * time.Second,
			//ResponseHeaderTimeout: 3600 * time.Second,
			//ExpectContinueTimeout: 3600 * time.Second,
			//MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}
