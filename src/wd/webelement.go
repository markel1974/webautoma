package wd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/markel1974/webautoma/src/wd/base"
	"io"
)

type WebElement struct {
	wd     *WebDriver
	client *Client
	id     string
}

func NewWebElement(wd *WebDriver, client *Client, id string) *WebElement {
	return &WebElement{
		wd:     wd,
		client: client,
		id:     id,
	}
}

func (elem *WebElement) Click() error {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/click", elem.id)
	return elem.client.VoidCommand(rUrl, nil)
}

func (elem *WebElement) SendKeys(keys string) error {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/value", elem.id)
	return elem.client.VoidCommand(rUrl, elem.wd.processKeyString(keys))
}

func (wd *WebDriver) processKeyString(keys string) interface{} {
	if !wd.w3cCompatible {
		//chars := make([]string, len(keys))
		//for i, c := range keys {
		//	chars[i] = string(c)
		//}
		b := []rune(keys)
		chars := make([]string, len(b))
		for i, c := range keys {
			chars[i] = string(c)
		}
		return map[string][]string{"value": chars}
	}
	return map[string]string{"text": keys}
}

func (elem *WebElement) TagName() (string, error) {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/name", elem.id)
	return elem.client.StringCommand(rUrl)
}

func (elem *WebElement) Text() (string, error) {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/text", elem.id)
	return elem.client.StringCommand(rUrl)
}

func (elem *WebElement) Submit() error {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/submit", elem.id)
	return elem.client.VoidCommand(rUrl, nil)
}

func (elem *WebElement) Clear() error {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/clear", elem.id)
	return elem.client.VoidCommand(rUrl, nil)
}

func (elem *WebElement) MoveTo(xOffset float64, yOffset float64) error {
	return elem.client.VoidCommand("/session/%s/moveto", map[string]interface{}{
		"element": elem.id,
		"xoffset": xOffset,
		"yoffset": yOffset,
	})
}

func (elem *WebElement) ComputedLabel() (string, error) {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/computedlabel", elem.id)
	response, err := elem.wd.computedLabel(rUrl)
	return response, err
}

func (elem *WebElement) FindElement(by string, value string) (base.IWebElement, error) {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/element", elem.id)
	response, err := elem.wd.find(by, value, "", rUrl)
	if err != nil {
		return nil, err
	}
	return elem.wd.DecodeElement(response)
}

func (elem *WebElement) FindElements(by string, value string) ([]base.IWebElement, error) {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/element", elem.id)
	response, err := elem.wd.find(by, value, "s", rUrl)
	if err != nil {
		return nil, err
	}
	return elem.wd.DecodeElements(response)
}

func (elem *WebElement) boolQuery(urlTemplate string) (bool, error) {
	return elem.client.BoolCommand(fmt.Sprintf(urlTemplate, elem.id))
}

func (elem *WebElement) IsSelected() (bool, error) {
	return elem.boolQuery("/session/%%s/element/%s/selected")
}

func (elem *WebElement) IsEnabled() (bool, error) {
	return elem.boolQuery("/session/%%s/element/%s/enabled")
}

func (elem *WebElement) IsDisplayed() (bool, error) {
	return elem.boolQuery("/session/%%s/element/%s/displayed")
}

func (elem *WebElement) GetProperty(name string) (string, error) {
	template := "/session/%%s/element/%s/property/%s"
	rUrl := fmt.Sprintf(template, elem.id, name)
	return elem.client.StringCommand(rUrl)
}

func (elem *WebElement) GetAttribute(name string) (string, error) {
	template := "/session/%%s/element/%s/attribute/%s"
	rUrl := fmt.Sprintf(template, elem.id, name)
	return elem.client.StringCommand(rUrl)
}

// rect implements the "Get Element Rect" method of the W3C standard.
func (elem *WebElement) rect() (*base.Rect, error) {
	rUrl := elem.client.RequestURL("/session/%s/element/%s/rect", elem.client.GetId(), elem.id)
	response, err := elem.client.Execute("GET", rUrl, nil)
	if err != nil {
		return nil, err
	}
	r := new(struct{ Value base.Rect })
	if err := json.Unmarshal(response, r); err != nil {
		return nil, err
	}
	return &r.Value, nil
}

func (elem *WebElement) CSSProperty(name string) (string, error) {
	return elem.client.StringCommand(fmt.Sprintf("/session/%%s/element/%s/css/%s", elem.id, name))
}

func (elem *WebElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"ELEMENT": elem.id, webElementIdentifier: elem.id})
}

func (elem *WebElement) Screenshot( /* scroll */ _ bool) ([]byte, error) {
	data, err := elem.client.StringCommand(fmt.Sprintf("/session/%%s/element/%s/screenshot", elem.id))
	if err != nil {
		return nil, err
	}
	buf := []byte(data)
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(buf))
	return io.ReadAll(decoder)
}

func (elem *WebElement) location(suffix string) (*base.Point, error) {
	if !elem.wd.w3cCompatible {
		rPath := "/session/%s/element/%s/location" + suffix
		rUrl := elem.client.RequestURL(rPath, elem.client.GetId(), elem.id)
		response, err := elem.client.Execute("GET", rUrl, nil)
		if err != nil {
			return nil, err
		}
		reply := new(struct{ Value base.Rect })
		if err := json.Unmarshal(response, reply); err != nil {
			return nil, err
		}
		return &base.Point{X: reply.Value.X, Y: reply.Value.Y}, nil
	}

	rect, err := elem.rect()
	if err != nil {
		return nil, err
	}
	return &base.Point{X: rect.X, Y: rect.Y}, nil
}

func (elem *WebElement) Location() (*base.Point, error) {
	return elem.location("")
}

func (elem *WebElement) LocationInView() (*base.Point, error) {
	return elem.location("_in_view")
}

func (elem *WebElement) Size() (*base.Size, error) {
	if !elem.wd.w3cCompatible {
		rUrl := elem.client.RequestURL("/session/%s/element/%s/size", elem.client.GetId(), elem.id)
		response, err := elem.client.Execute("GET", rUrl, nil)
		if err != nil {
			return nil, err
		}
		reply := new(struct{ Value base.Rect })
		if err := json.Unmarshal(response, reply); err != nil {
			return nil, err
		}
		return &base.Size{Width: reply.Value.Width, Height: reply.Value.Height}, nil
	}

	rect, err := elem.rect()
	if err != nil {
		return nil, err
	}

	return &base.Size{Width: rect.Width, Height: rect.Height}, nil
}

func (elem *WebElement) ScrollTo(deltaX int, deltaY int) error {
	loc, err := elem.Location()
	if err != nil {
		return err
	}
	elem.wd.StoreWheelActions("wheel1", base.CreateWheelAction(int(loc.X), int(loc.Y), deltaX, deltaY))
	if err = elem.wd.PerformActions(); err != nil {
		return err
	}
	if err = elem.wd.ReleaseActions(); err != nil {
		return err
	}
	return nil
}

func (elem *WebElement) Print() {
	text, _ := elem.Text()
	size, _ := elem.Size()
	isDisplayed, _ := elem.IsDisplayed()
	isEnabled, _ := elem.IsEnabled()
	isSelected, _ := elem.IsSelected()
	tagName, _ := elem.TagName()
	fmt.Printf("text: %s\n", text)
	fmt.Printf("size: w: %f h: %f\n", size.Width, size.Height)
	fmt.Printf("isEnabled: %v\n", isEnabled)
	fmt.Printf("isDisplayed: %v\n", isDisplayed)
	fmt.Printf("isSelected: %v\n", isSelected)
	fmt.Printf("tagName: %s\n", tagName)
}
