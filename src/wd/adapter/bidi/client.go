package bidi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

// Command to send to browser
type Command struct {
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Channel string      `json:"channel,omitempty"`
}

// APIResponse from browser
type APIResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
}

// APIError from browser
type APIError struct {
	ErrorCode string `json:"error"`
	Message   string `json:"message"`
}

// Error stdlib interface
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

// Event from browser
type Event struct {
	SessionID string          `json:"sessionId,omitempty"`
	Method    string          `json:"method"`
	Params    json.RawMessage `json:"params,omitempty"`
}

// EventCallback represents a callback event, associated with a method.
type EventCallback func(params json.RawMessage) error

type Client struct {
	count     uint64
	pending   sync.Map    // pending requests
	event     chan *Event // events from browser
	ws        *WebSocket
	callbacks map[string]EventCallback
}

func NewBiDiClient() *Client {
	return &Client{
		event:     make(chan *Event),
		ws:        &WebSocket{},
		callbacks: map[string]EventCallback{},
	}
}

func (c *Client) Start(wsURL string, header http.Header) error {
	if err := c.ws.Connect(context.Background(), wsURL, header); err != nil {
		return err
	}

	go c.processEvents()
	go c.readMessages()

	return nil
}

type result struct {
	msg json.RawMessage
	err error
}

func (c *Client) Call(ctx context.Context, method string, params interface{}) ([]byte, error) {
	command := &Command{
		ID:     int(c.newID()),
		Method: method,
		Params: params,
	}

	done := make(chan result)

	once := sync.Once{}

	c.pending.Store(command.ID, func(res result) {
		once.Do(func() {
			select {
			case <-ctx.Done():
			case done <- res:
			}
		})
	})

	defer c.pending.Delete(command.ID)

	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	if err := c.ws.Write(context.Background(), data); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-done:
		return res.msg, res.err
	}
}

func (c *Client) Close() error {
	return c.ws.Close()
}

func (c *Client) CallbackEvent(method string, cb EventCallback) {
	c.callbacks[method] = cb
}

// Read messages coming from the browser via the websocket.
func (c *Client) readMessages() {
	defer close(c.event)

	for {
		data, err := c.ws.Read(context.Background())
		if err != nil {
			c.pending.Range(func(_, val interface{}) bool {
				val.(func(result))(result{err: err})
				return true
			})

			return
		}

		var id struct {
			ID int `json:"id"`
		}

		if err = json.Unmarshal(data, &id); err != nil {
			panic(err)
		}

		if id.ID == 0 {
			var evt Event
			if err = json.Unmarshal(data, &evt); err != nil {
				panic(err)
			}

			c.event <- &evt

			continue
		}

		val, ok := c.pending.Load(id.ID)
		if !ok {
			continue
		}

		var (
			apiRes APIResponse
			apiErr APIError
		)

		if err = json.Unmarshal(data, &apiErr); err == nil && apiErr.ErrorCode != "" {
			val.(func(result))(result{nil, apiErr})
			continue
		} else if err = json.Unmarshal(data, &apiRes); err == nil {
			val.(func(result))(result{apiRes.Result, nil})
			continue
		} else {
			panic(err)
		}
	}
}

// Process events coming from the browser via the websocket.
func (c *Client) processEvents() {
	for event := range c.event {
		if cb, ok := c.callbacks[event.Method]; ok {
			if err := cb(event.Params); err != nil {
				panic(err)
			}
		}
	}
}

func (c *Client) newID() uint64 {
	return atomic.AddUint64(&c.count, 1)
}
