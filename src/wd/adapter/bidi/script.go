package bidi

type Source struct {
	Realm   string `json:"realm"`
	Context string `json:"context"`
}

type StackFrame struct {
	ColumnNumber int    `json:"columnNumber"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
	URL          string `json:"url"`
}

type StackTrace struct {
	CallFrames []*StackFrame `json:"callFrames"`
}
