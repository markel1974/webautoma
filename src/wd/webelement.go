package wd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/markel1974/webautoma/src/wd/base"
	"io/ioutil"
)

type WebElement struct {
	parent *WebDriver
	id     string
}

func (elem *WebElement) Click() error {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/click", elem.id)
	return elem.parent.voidCommand(urlTemplate, nil)
}

func (elem *WebElement) SendKeys(keys string) error {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/value", elem.id)
	return elem.parent.voidCommand(urlTemplate, elem.parent.processKeyString(keys))
}

func (wd *WebDriver) processKeyString(keys string) interface{} {
	if !wd.w3cCompatible {
		chars := make([]string, len(keys))
		for i, c := range keys {
			chars[i] = string(c)
		}
		return map[string][]string{"value": chars}
	}
	return map[string]string{"text": keys}
}

func (elem *WebElement) TagName() (string, error) {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/name", elem.id)
	return elem.parent.stringCommand(urlTemplate)
}

func (elem *WebElement) Text() (string, error) {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/text", elem.id)
	return elem.parent.stringCommand(urlTemplate)
}

func (elem *WebElement) Submit() error {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/submit", elem.id)
	return elem.parent.voidCommand(urlTemplate, nil)
}

func (elem *WebElement) Clear() error {
	urlTemplate := fmt.Sprintf("/session/%%s/element/%s/clear", elem.id)
	return elem.parent.voidCommand(urlTemplate, nil)
}

func (elem *WebElement) MoveTo(xOffset float64, yOffset float64) error {
	return elem.parent.voidCommand("/session/%s/moveto", map[string]interface{}{
		"element": elem.id,
		"xoffset": xOffset,
		"yoffset": yOffset,
	})
}

func (elem *WebElement) FindElement(by string, value string) (base.IWebElement, error) {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/element", elem.id)
	response, err := elem.parent.find(by, value, "", rUrl)
	if err != nil {
		return nil, err
	}

	return elem.parent.DecodeElement(response)
}

func (elem *WebElement) FindElements(by string, value string) ([]base.IWebElement, error) {
	rUrl := fmt.Sprintf("/session/%%s/element/%s/element", elem.id)
	response, err := elem.parent.find(by, value, "s", rUrl)
	if err != nil {
		return nil, err
	}
	return elem.parent.DecodeElements(response)
}

func (elem *WebElement) boolQuery(urlTemplate string) (bool, error) {
	return elem.parent.boolCommand(fmt.Sprintf(urlTemplate, elem.id))
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
	urlTemplate := fmt.Sprintf(template, elem.id, name)
	return elem.parent.stringCommand(urlTemplate)
}

func (elem *WebElement) GetAttribute(name string) (string, error) {
	template := "/session/%%s/element/%s/attribute/%s"
	urlTemplate := fmt.Sprintf(template, elem.id, name)
	return elem.parent.stringCommand(urlTemplate)
}

// rect implements the "Get Element Rect" method of the W3C standard.
func (elem *WebElement) rect() (*base.Rect, error) {
	wd := elem.parent
	rUrl := wd.requestURL("/session/%s/element/%s/rect", wd.id, elem.id)
	response, err := wd.execute("GET", rUrl, nil)
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
	wd := elem.parent
	return wd.stringCommand(fmt.Sprintf("/session/%%s/element/%s/css/%s", elem.id, name))
}

func (elem *WebElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"ELEMENT":            elem.id,
		webElementIdentifier: elem.id,
	})
}

func (elem *WebElement) Screenshot( /* scroll */ _ bool) ([]byte, error) {
	data, err := elem.parent.stringCommand(fmt.Sprintf("/session/%%s/element/%s/screenshot", elem.id))
	if err != nil {
		return nil, err
	}
	buf := []byte(data)
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(buf))
	return ioutil.ReadAll(decoder)
}

func (elem *WebElement) location(suffix string) (*base.Point, error) {
	if !elem.parent.w3cCompatible {
		wd := elem.parent
		rPath := "/session/%s/element/%s/location" + suffix
		rUrl := wd.requestURL(rPath, wd.id, elem.id)
		response, err := wd.execute("GET", rUrl, nil)
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
	if !elem.parent.w3cCompatible {
		wd := elem.parent
		rUrl := wd.requestURL("/session/%s/element/%s/size", wd.id, elem.id)
		response, err := wd.execute("GET", rUrl, nil)
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
