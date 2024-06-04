package executor

import (
	"encoding/json"
	"github.com/markel1974/webautoma/src/wd/base"
	"log"
	"regexp"
	"strings"
)

var _supportedMethod = []string{"Network.responseReceived"}

var _supportedContentType = []string{"text/html", "json", "xml"}

type NetworkMessage struct {
	Message struct {
		Method string `json:"method"`
		Params struct {
			Response struct {
				Headers map[string]interface{} `json:"headers"`
				Timing  map[string]interface{} `json:"timing"`
				Status  int                    `json:"status"`
				URL     string                 `json:"url"`
			} `json:"response"`
		} `json:"params"`
	} `json:"message"`
}

type Network struct {
	filter   *regexp.Regexp
	capture  []string
	acquired map[string]interface{}
}

func NewNetwork(capture []string) *Network {
	return &Network{
		filter:  regexp.MustCompile("(http|file|ftp|png|jpg|gif|js|css|mp4|ico|bmp)"),
		capture: capture,
	}
}

func (n *Network) Acquired() map[string]interface{} {
	return n.acquired
}

func (n *Network) Compute(logEntries []base.LogMessage) (map[string]interface{}, int) {
	errorCount := 0
	var headersData map[string]interface{}

	for _, entry := range logEntries {
		var network NetworkMessage
		if err := json.Unmarshal([]byte(entry.Message), &network); err != nil {
			log.Println(err.Error())
			continue
		}

		hasMethodSupport := false
		for _, v := range _supportedMethod {
			if v == network.Message.Method {
				hasMethodSupport = true
				break
			}
		}
		if !hasMethodSupport {
			continue
		}

		currentURL := network.Message.Params.Response.URL
		if n.filter.MatchString(currentURL) {
			currentHeaders := network.Message.Params.Response.Headers
			contentType, _ := MapToString(currentHeaders, "content-type")
			hasContentSupport := false
			for _, v := range _supportedContentType {
				if strings.Contains(contentType, v) {
					hasContentSupport = true
					break
				}
			}
			if !hasContentSupport {
				continue
			}
			currentStatus := network.Message.Params.Response.Status
			for _, capture := range n.capture {
				if res, ok := MapToString(currentHeaders, capture); ok {
					if n.acquired == nil {
						n.acquired = make(map[string]interface{})
					}
					n.acquired[capture] = res
				}
			}
			if currentStatus >= 400 {
				errorCount++
			}
			if headersData == nil {
				headersData = make(map[string]interface{})
			}
			headersData[currentURL] = map[string]interface{}{
				"url":     currentURL,
				"headers": currentHeaders,
				//"timing":       currenTiming,
				"content-type":  contentType,
				"status":        currentStatus,
				"messageMethod": network.Message.Method,
			}
		}
	}
	network := map[string]interface{}{
		"errorCount": errorCount,
		"result":     headersData,
	}
	return network, errorCount
}
