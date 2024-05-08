package adapter

// LocatorStrategy Strategy for searching element on the page
type LocatorStrategy string

const (
	LocatorStrategyCSSSelector     LocatorStrategy = "css selector"
	LocatorStrategyLinkText        LocatorStrategy = "link text"
	LocatorStrategyPartialLinkText LocatorStrategy = "partial link text"
	LocatorStrategyTagName         LocatorStrategy = "tag name"
	LocatorStrategyXPath           LocatorStrategy = "xpath"
)
