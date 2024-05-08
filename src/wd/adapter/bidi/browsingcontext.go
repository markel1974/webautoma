package bidi

import (
	"context"
	"encoding/json"
)

type BrowsingContextType string

const (
	BrowsingContextTypeTab    BrowsingContextType = "tab"
	BrowsingContextTypeWindow BrowsingContextType = "window"
)

type BrowsingContext struct {
	ID     string              `json:"context"`
	Type   BrowsingContextType `json:"-"`
	client *Client             `json:"-"`
}

// func (b *BrowsingContext) CaptureScreenshot() ([]byte, error) {
// 	data, err := b.client.Call(context.Background(), "browsingContext.captureScreenshot", map[string]interface{}{
// 		"context": b.ID,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	var screenshot string
// 	if err := json.Unmarshal(data, &screenshot); err != nil {
// 		return nil, err
// 	}

// 	return base64.StdEncoding.DecodeString(screenshot)
// }

func (b *BrowsingContext) Close() error {
	_, err := b.client.Call(context.Background(), "browsingContext.close", map[string]interface{}{
		"context": b.ID,
	})

	return err
}

func (b *BrowsingContext) HandleUserPrompt(accept bool, userText string) error {
	_, err := b.client.Call(context.Background(), "browsingContext.handleUserPrompt", map[string]interface{}{
		"context":  b.ID,
		"accept":   accept,
		"userText": userText,
	})

	return err
}

type BrowsingContextReadinessState string

const (
	BrowsingContextReadinessStateNone        BrowsingContextReadinessState = "none"
	BrowsingContextReadinessStateInteractive BrowsingContextReadinessState = "interactive"
	BrowsingContextReadinessStateComplete    BrowsingContextReadinessState = "complete"
)

func (b *BrowsingContext) Navigate(url string, wait BrowsingContextReadinessState) (*Navigation, error) {
	data, err := b.client.Call(context.Background(), "browsingContext.navigate", map[string]interface{}{
		"context": b.ID,
		"url":     url,
		"wait":    wait,
	})
	if err != nil {
		return nil, err
	}

	navigation := &Navigation{}
	err = json.Unmarshal(data, &navigation)

	return navigation, err
}

func (b *BrowsingContext) Reload(ignoreCache bool, wait BrowsingContextReadinessState) error {
	_, err := b.client.Call(context.Background(), "browsingContext.reload", map[string]interface{}{
		"context":     b.ID,
		"ignoreCache": ignoreCache,
		"wait":        wait,
	})

	return err
}
