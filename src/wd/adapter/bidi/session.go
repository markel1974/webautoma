package bidi

import (
	"context"
	"encoding/json"
	"net/http"
)

type Session struct {
	ID     string
	client *Client `json:"-"`
}

func New(wsURL string, header http.Header) (*Session, error) {
	client := NewBiDiClient()

	if err := client.Start(wsURL, header); err != nil {
		return nil, err
	}

	return &Session{
		client: client,
	}, nil
}

type Status struct {
	Ready   bool   `json:"ready"`
	Message string `json:"message"`
}

func (s *Session) Status() (*Status, error) {
	data, err := s.client.Call(context.Background(), "session.status", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	status := &Status{}
	err = json.Unmarshal(data, &status)

	return status, err
}

func (s *Session) Subscribe(events []string) error {
	_, err := s.client.Call(context.Background(), "session.subscribe", map[string]interface{}{
		"events": events,
	})

	return err
}

func (s *Session) UnSubscribe(events []string) error {
	_, err := s.client.Call(context.Background(), "session.unsubscribe", map[string]interface{}{
		"events": events,
	})

	return err
}

func (s *Session) NewBrowsingContext(contextType BrowsingContextType, refContext *BrowsingContext) (*BrowsingContext, error) {
	params := map[string]interface{}{
		"type": contextType,
	}

	if refContext != nil {
		params["referenceContext"] = refContext.ID
	}

	data, err := s.client.Call(context.Background(), "browsingContext.create", params)
	if err != nil {
		return nil, err
	}

	context := &BrowsingContext{Type: contextType}
	if err := json.Unmarshal(data, context); err != nil {
		return nil, err
	}

	context.client = s.client

	return context, nil
}
