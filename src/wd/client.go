package wd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/markel1974/webautoma/src/wd/base"
	"io"
	"mime"
	"net/http"
	"time"
)

func newDefaultClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//IMPORTANT chromedriver support only 6 connections
			MaxConnsPerHost:       4,
			MaxIdleConnsPerHost:   2,
			MaxIdleConns:          2,
			ResponseHeaderTimeout: 60 * time.Second,
			//MaxIdleConns:          10,
			//IdleConnTimeout:       3600 * time.Second,
			//TLSHandshakeTimeout:   3600 * time.Second,
			//ResponseHeaderTimeout: 3600 * time.Second,
			//ExpectContinueTimeout: 3600 * time.Second,
			//MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}

type Client struct {
	id        string
	client    *http.Client
	urlPrefix string
	debug     bool
}

func NewClient(client *http.Client, urlPrefix string, debug bool) *Client {
	if client == nil {
		client = newDefaultClient()
	}
	return &Client{
		id:        "",
		client:    client,
		urlPrefix: urlPrefix,
		debug:     debug,
	}
}

func (c *Client) SetId(id string) {
	c.id = id
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) RequestURL(template string, args ...interface{}) string {
	u := c.urlPrefix + fmt.Sprintf(template, args...)
	return u
}

func (c *Client) VoidCommand(urlTemplate string, params interface{}) error {
	rUrl := c.RequestURL(urlTemplate, c.id)
	err := c.VoidParamsCommand("POST", rUrl, params)
	return err
}

func (c *Client) VoidParamsCommand(method, url string, params interface{}) error {
	if params == nil {
		params = make(map[string]interface{})
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, err = c.Execute(method, url, data)
	return err
}

func (c *Client) StringCommand(urlTemplate string) (string, error) {
	rUrl := c.RequestURL(urlTemplate, c.id)
	response, err := c.Execute("GET", rUrl, nil)
	if err != nil {
		return "", err
	}

	reply := new(struct{ Value *string })
	if err := json.Unmarshal(response, reply); err != nil {
		return "", err
	}

	if reply.Value == nil {
		return "", fmt.Errorf("nil return value")
	}

	return *reply.Value, nil
}

func (c *Client) StringsCommand(urlTemplate string) ([]string, error) {
	rUrl := c.RequestURL(urlTemplate, c.id)
	response, err := c.Execute("GET", rUrl, nil)
	if err != nil {
		return nil, err
	}
	reply := new(struct{ Value []string })
	if err := json.Unmarshal(response, reply); err != nil {
		return nil, err
	}
	return reply.Value, nil
}

func (c *Client) BoolCommand(urlTemplate string) (bool, error) {
	rUrl := c.RequestURL(urlTemplate, c.id)
	response, err := c.Execute("GET", rUrl, nil)
	if err != nil {
		return false, err
	}
	reply := new(struct{ Value bool })
	if err := json.Unmarshal(response, reply); err != nil {
		return false, err
	}
	return reply.Value, nil
}

func (c *Client) newRequest(method string, url string, data []byte) (*http.Request, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", DefaultContentType)
	return request, nil
}

func (c *Client) Execute(method string, url string, data []byte) (json.RawMessage, error) {
	buildError := func(err error) error {
		return fmt.Errorf("[method: %s] [url: %s] [error: %s]", method, url, err.Error())
	}
	request, err := c.newRequest(method, url, data)
	if err != nil {
		return nil, buildError(err)
	}

	var response *http.Response
	var buf []byte

	if response, err = c.client.Do(request); err == nil {
		if response.Body != nil {
			buf, err = io.ReadAll(response.Body)
			_ = response.Body.Close()
		} else {
			err = buildError(errors.New("empty body"))
		}
	}

	if err != nil {
		return nil, buildError(err)
	}

	if c.debug {
		fmt.Printf("%s %s %s --> %s [%s]\n", method, url, data, response.Status, response.Header["Content-Type"])
		if err != nil {
			fmt.Println(buildError(err))
		} else {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, buf, "", "	"); err == nil {
				fmt.Println(pretty.String())
			} else {
				fmt.Println(string(buf))
			}
		}
	}

	if err != nil {
		return nil, buildError(errors.New(response.Status))
	}

	fullCType := response.Header.Get("Content-Type")
	cType, _, err := mime.ParseMediaType(fullCType)
	if err != nil {
		return nil, buildError(fmt.Errorf("got content type header %q, expected %q", fullCType, DefaultContentType))
	}
	if cType != DefaultContentType {
		return nil, buildError(fmt.Errorf("got content type %q, expected %q", cType, DefaultContentType))
	}

	reply := new(ServerReply)
	if err = json.Unmarshal(buf, reply); err != nil {
		if response.StatusCode != http.StatusOK {
			return nil, buildError(fmt.Errorf("bad server reply status: %s", response.Status))
		}
		return nil, buildError(err)
	}
	if reply.Err != "" {
		return nil, buildError(&reply.Error)
	}

	// Handle the W3C-compliant error format. In the W3C spec, the error is embedded in the 'value' field.
	if len(reply.Value) > 0 {
		respErr := new(base.Error)
		if err = json.Unmarshal(reply.Value, respErr); err == nil && respErr.Err != "" {
			respErr.HTTPCode = response.StatusCode
			return nil, buildError(respErr)
		}
	}

	// Handle the legacy error format.
	const success = 0
	if reply.Status != success {
		shortMsg, ok := base.RemoteErrors[reply.Status]
		if !ok {
			shortMsg = fmt.Sprintf("unknown error - %d", reply.Status)
		}
		longMsg := new(struct {
			Message string
		})
		if err = json.Unmarshal(reply.Value, longMsg); err != nil {
			return nil, buildError(errors.New(shortMsg))
		}
		return nil, &base.Error{
			Err:        shortMsg,
			Message:    longMsg.Message,
			HTTPCode:   response.StatusCode,
			LegacyCode: reply.Status,
		}
	}
	return buf, nil
}
