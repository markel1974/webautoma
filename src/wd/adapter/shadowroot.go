package adapter

import (
	"encoding/json"
	"fmt"
)

type ShadowRoot struct {
	ID        string `json:"shadow-6066-11e4-a52e-4f735466cecf"`
	SessionID string `json:"-"`
	client    *RestClient
}

func (s *ShadowRoot) FindElement(strategy LocatorStrategy, selector string) (*Element, error) {
	data, err := s.client.Post(fmt.Sprintf("/session/%s/shadow/%s/element", s.SessionID, s.ID), &Params{
		"using": strategy,
		"value": selector,
	})
	if err != nil {
		return nil, err
	}
	element := Element{}
	if err := json.Unmarshal(data, &element); err != nil {
		return nil, err
	}
	element.SessionID = s.ID
	element.client = s.client
	return &element, nil
}

func (s *ShadowRoot) FindElements(strategy LocatorStrategy, selector string) ([]Element, error) {
	data, err := s.client.Post(fmt.Sprintf("/session/%s/shadow/%s/elements", s.SessionID, s.ID), &Params{
		"using": strategy,
		"value": selector,
	})
	if err != nil {
		return nil, err
	}
	var elements []Element
	if err := json.Unmarshal(data, &elements); err != nil {
		return nil, err
	}
	for _, element := range elements {
		element.SessionID = s.ID
		element.client = s.client
	}
	return elements, nil
}
