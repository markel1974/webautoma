package executor

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/markel1974/webautoma/src/wd/base"
)

// _supportedMethod is a map that defines supported network methods with corresponding integer identifiers for processing.
var _supportedMethod = map[string]int{"Network.requestWillBeSent": 0, "Network.responseReceived": 1}

// _supportedContentType is a list of supported content types used for filtering based on content handling logic.
var _supportedContentType = []string{"text/html", "json", "xml"}

// NetworkMessage represents the structure of a network log message containing request and response data.
type NetworkMessage struct {
	Message struct {
		Method string `json:"method"`
		Params struct {
			Request struct {
				Headers map[string]interface{} `json:"headers"`
				URL     string                 `json:"url"`
			} `json:"request"`
			Response struct {
				Headers map[string]interface{} `json:"headers"`
				Timing  map[string]interface{} `json:"timing"`
				Status  int                    `json:"status"`
				URL     string                 `json:"url"`
			} `json:"response"`
		} `json:"params"`
	} `json:"message"`
}

// rgxDef is a constant string prefix used for identifying or marking regular expression-based definitions.
const rgxDef = "rgx:"

// Network is a structure for managing and processing network-related data and messages with filtering and capturing features.
type Network struct {
	filter           *regexp.Regexp
	capture          []string
	acquired         map[string]interface{}
	supportedMethods map[string]int
	//supportedMethods    []string
	//supportedMethodsRgx []*regexp.Regexp
}

// NewNetwork initializes a Network instance with capture filters and supported methods config. Returns the instance or an error.
func NewNetwork(capture []string, sm []string) (*Network, error) {
	/*
		var supportedMethods []string
		var supportedMethodsRgx []*regexp.Regexp

		if sm != nil {
			for _, m := range sm {
				if strings.HasPrefix(m, rgxDef) {
					m = m[len(rgxDef):]
					r, err := regexp.Compile(m)
					if err != nil {
						return nil, err
					}
					supportedMethodsRgx = append(supportedMethodsRgx, r)
				} else {
					supportedMethods = append(supportedMethods, m)
				}
			}
		}

		if len(supportedMethods) == 0 && len(supportedMethodsRgx) == 0 {
			supportedMethods = _supportedMethod
		}
	*/
	supportedMethods := _supportedMethod

	return &Network{
		filter:           regexp.MustCompile("(http|file|ftp|png|jpg|gif|js|css|mp4|ico|bmp)"),
		capture:          capture,
		supportedMethods: supportedMethods,
		//supportedMethodsRgx: supportedMethodsRgx,
	}, nil
}

// Acquired retrieves the map of acquired data from the network instance, typically containing captured headers or parameters.
func (n *Network) Acquired() map[string]interface{} {
	return n.acquired
}

// Headers processes a slice of log entries and extracts network-related data such as headers, status, and content type.
func (n *Network) Headers(logEntries []base.LogMessage) map[string]interface{} {
	headersData := make(map[string]interface{})
	for _, entry := range logEntries {
		var network NetworkMessage
		if err := json.Unmarshal([]byte(entry.Message), &network); err != nil {
			log.Println(err.Error())
			continue
		}
		currentURL, currentStatus, currentHeaders, contentType := n.getResponseParams(network)
		headersData[currentURL] = map[string]interface{}{
			"url":           currentURL,
			"headers":       currentHeaders,
			"content-type":  contentType,
			"status":        currentStatus,
			"messageMethod": network.Message.Method,
		}
	}
	return headersData
}

// Compute processes a list of log entries to extract network data and counts errors with status codes >= 400.
// It returns a map containing network details and the total number of errors encountered.
func (n *Network) Compute(logEntries []base.LogMessage) (map[string]interface{}, int) {
	errorCount := 0
	var headersData map[string]interface{}

	for _, entry := range logEntries {
		var network NetworkMessage
		if err := json.Unmarshal([]byte(entry.Message), &network); err != nil {
			log.Println(err.Error())
			continue
		}
		hasMethodSupport := -1
		if v, ok := _supportedMethod[network.Message.Method]; ok {
			hasMethodSupport = v
		}

		var currentURL string
		var currentStatus int
		var currentHeaders map[string]interface{}
		var contentType string

		if hasMethodSupport == 0 {
			currentURL, currentStatus, currentHeaders, contentType = n.getRequestParams(network)
		} else if hasMethodSupport == 1 {
			currentURL, currentStatus, currentHeaders, contentType = n.getResponseParams(network)
		} else {
			continue
		}

		if n.filter.MatchString(currentURL) {
			hasContentSupport := false
			if len(contentType) == 0 {
				hasContentSupport = true
			} else {
				for _, v := range _supportedContentType {
					if strings.Contains(contentType, v) {
						hasContentSupport = true
						break
					}
				}
			}
			if !hasContentSupport {
				continue
			}
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
		//TODO le url non possono stare tra le chiavi....
		//"result":     headersData,
	}
	//if data {
	//	network["result"] = headersData
	//}
	return network, errorCount
}

// getRequestParams extracts the request URL, status, headers, and content type from a NetworkMessage instance.
func (n *Network) getRequestParams(nm NetworkMessage) (string, int, map[string]interface{}, string) {
	currentURL := nm.Message.Params.Request.URL
	currentHeaders := nm.Message.Params.Request.Headers
	currentStatus := 0
	contentType := ""
	return currentURL, currentStatus, currentHeaders, contentType
}

// getResponseParams extracts and returns the URL, status, headers, and content type from a response in a NetworkMessage.
func (n *Network) getResponseParams(nm NetworkMessage) (string, int, map[string]interface{}, string) {
	currentURL := nm.Message.Params.Response.URL
	currentHeaders := nm.Message.Params.Response.Headers
	currentStatus := nm.Message.Params.Response.Status
	contentType, _ := MapToString(currentHeaders, "content-type")
	return currentURL, currentStatus, currentHeaders, contentType
}
