package base

// IWebElement defines method supported by web elements.
type IWebElement interface {

	// Click clicks on the element.
	Click() error

	// SendKeys types into the element.
	SendKeys(keys string) error

	// Submit submits the button.
	Submit() error

	// Clear clears the element.
	Clear() error

	// MoveTo moves the mouse to relative coordinates from center of element, If
	// the element is not visible, it will be scrolled into view.
	MoveTo(xOffset float64, yOffset float64) error

	// ComputedLabel ...
	ComputedLabel() (string, error)

	// FindElement finds a child element.
	FindElement(by string, value string) (IWebElement, error)

	// FindElements finds multiple children elements.
	FindElements(by string, value string) ([]IWebElement, error)

	// TagName returns the element's name.
	TagName() (string, error)

	// Text returns the text of the element.
	Text() (string, error)

	// IsSelected returns true if element is selected.
	IsSelected() (bool, error)

	// IsEnabled returns true if the element is enabled.
	IsEnabled() (bool, error)

	// IsDisplayed returns true if the element is displayed.
	IsDisplayed() (bool, error)

	// GetAttribute returns the named HTML attribute of the element.
	GetAttribute(name string) (string, error)

	// GetProperty returns the DOM property of the element. The DOM property
	// values can change (e.g. input value), the HTML attributes can't.
	GetProperty(name string) (string, error)

	// Location returns the element's location.
	Location() (*Point, error)

	// LocationInView returns the element's location once it has been scrolled
	// into view.
	LocationInView() (*Point, error)

	// Size returns the element's size.
	Size() (*Size, error)

	// CSSProperty returns the value of the specified CSS property of the
	// element.
	CSSProperty(name string) (string, error)

	// Screenshot takes a screenshot of the attribute scrolling if necessary.
	Screenshot(scroll bool) ([]byte, error)

	// Print information
	Print()

	// ScrollTo ...
	ScrollTo(x int, y int) error
}
