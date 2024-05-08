package adapter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Element struct {
	ID        string      `json:"element-6066-11e4-a52e-4f735466cecf"`
	SessionID string      `json:"-"`
	client    *RestClient `json:"-"`
}

// GetShadowRoot returns a shadow root of the element if there is one or an error.
func (e *Element) GetShadowRoot() (*ShadowRoot, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/shadow", e.SessionID, e.ID))
	if err != nil {
		return nil, err
	}

	shadowRoot := ShadowRoot{}
	if err = json.Unmarshal(data, &shadowRoot); err != nil {
		return nil, err
	}

	shadowRoot.SessionID = e.SessionID
	shadowRoot.client = e.client

	return &shadowRoot, err
}

// FindElement searches for an element on the page, starting from the referenced web element.
func (e *Element) FindElement(strategy LocatorStrategy, selector string) (*Element, error) {
	data, err := e.client.Post(fmt.Sprintf("/session/%s/element/%s", e.SessionID, e.ID), &Params{
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

	element.SessionID = e.SessionID
	element.client = e.client

	return &element, nil
}

// FindElements searches for multiple elements on the page, starting from the referenced web element. The located
// elements will be returned as a WebElement JSON objects. The table below lists the locator
// strategies that each server should support. Elements should be returned in the order located
// in the DOM.
func (e *Element) FindElements(strategy LocatorStrategy, selector string) ([]Element, error) {
	data, err := e.client.Post(fmt.Sprintf("/session/%s/elements/%s", e.SessionID, e.ID), &Params{
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
		element.SessionID = e.SessionID
		element.client = e.client
	}

	return elements, nil
}

// IsSelected determines if the referenced element is selected or not.
// This operation only makes sense on input elements of the Checkbox- and Radio Button states, or on option elements.
func (e *Element) IsSelected() (bool, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/selected", e.SessionID, e.ID))
	if err != nil {
		return false, err
	}

	var selected bool
	err = json.Unmarshal(data, &selected)

	return selected, err
}

// GetAttribute returns the attribute value of the referenced web element.
func (e *Element) GetAttribute(name string) (string, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/attribute/%s", e.SessionID, e.ID, name))
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal(data, &value)

	return value, err
}

// GetProperty returns the property of the referenced web element.
func (e *Element) GetProperty(name string) (string, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/property/%s", e.SessionID, e.ID, name))
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal(data, &value)

	return value, err
}

// GetCSSValue returns the computed value of the given CSS property for the element.
func (e *Element) GetCSSValue(name string) (string, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/css/%s", e.SessionID, e.ID, name))
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal(data, &value)

	return value, err
}

// GetText returns the visible text for the element.
func (e *Element) GetText() (string, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/text", e.SessionID, e.ID))
	if err != nil {
		return "", err
	}

	var text string
	err = json.Unmarshal(data, &text)

	return text, err
}

// GetTagName returns the tagName of an element
func (e *Element) GetTagName() (string, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/name", e.SessionID, e.ID))
	if err != nil {
		return "", err
	}

	var name string
	err = json.Unmarshal(data, &name)

	return name, err
}

// ElementRect defines the Element Rect.
type ElementRect struct {
	// X axis position of the top-left corner of the element relative to the current browsing context's document element in CSS pixels
	X int `json:"x"`

	// Y axis position of the top-left corner of the element relative to the current browsing context's document element in CSS pixels
	Y int `json:"y"`

	// Height of the element's bounding rectangle in CSS pixels
	Width int `json:"width"`

	// Width of the web element's bounding rectangle in CSS pixels
	Height int `json:"height"`
}

// GetRect Returns the dimensions and coordinates of the referenced element
func (e *Element) GetRect() (*ElementRect, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/rect", e.SessionID, e.ID))
	if err != nil {
		return nil, err
	}

	elementRect := &ElementRect{}
	err = json.Unmarshal(data, &elementRect)

	return elementRect, err
}

// IsEnabled determines if the referenced element is enabled or not.
func (e *Element) IsEnabled() (bool, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/enabled", e.SessionID, e.ID))
	if err != nil {
		return false, err
	}

	var enabled bool
	err = json.Unmarshal(data, &enabled)

	return enabled, err
}

// Click clicks on an element.
func (e *Element) Click() error {
	_, err := e.client.Post(fmt.Sprintf("/session/%s/element/%s/click", e.SessionID, e.ID), nil)
	return err
}

// Clear clears content of an element.
func (e *Element) Clear() error {
	_, err := e.client.Post(fmt.Sprintf("/session/%s/element/%s/clear", e.SessionID, e.ID), nil)
	return err
}

// SendKeys sends a sequence of key strokes to an element.
func (e *Element) SendKeys(text string) error {
	_, err := e.client.Post(fmt.Sprintf("/session/%s/element/%s/value", e.SessionID, e.ID), &Params{
		"text": text,
	})

	return err
}

// TakeScreenshot takes a screenshot of the visible region encompassed by the bounding rectangle of an element.
func (e *Element) TakeScreenshot() ([]byte, error) {
	data, err := e.client.Get(fmt.Sprintf("/session/%s/element/%s/screenshot", e.SessionID, e.ID))
	if err != nil {
		return nil, err
	}

	var screenshot string
	if err := json.Unmarshal(data, &screenshot); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(screenshot)
}
