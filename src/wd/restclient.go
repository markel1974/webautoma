package wd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RestClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewRestClient(baseURL string) *RestClient {
	return &RestClient{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
	}
}

type APIResponse struct {
	Value json.RawMessage `json:"value"`
}

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (rc *RestClient) Do(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-charset", "utf-8")

	res, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := &APIResponse{}
	if err := json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	apiError := &APIError{}
	if err := json.Unmarshal(response.Value, apiError); err == nil && apiError.Error != "" {
		return nil, errors.New(apiError.Error)
	}

	return response.Value, nil
}

func (rc *RestClient) Get(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", rc.baseURL, path), nil)
	if err != nil {
		return nil, err
	}

	return rc.Do(req)
}

type Params map[string]interface{}

func (rc *RestClient) Post(path string, data *Params) ([]byte, error) {
	if data == nil {
		data = &Params{}
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", rc.baseURL, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	return rc.Do(req)
}

func (rc *RestClient) Delete(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", rc.baseURL, path), nil)
	if err != nil {
		return nil, err
	}

	return rc.Do(req)
}
