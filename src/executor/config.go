package executor

// ConfigCommand represents a configuration command with fields for command details, coordinates, and window properties.
type ConfigCommand struct {
	Id      string `json:"id"`
	Target  string `json:"target"`
	Command string `json:"command"`
	Until   string `json:"until"`
	Value   string `json:"value"`

	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	OffsetX float64 `json:"offsetX"`
	OffsetY float64 `json:"offsetY"`

	WindowHandleName string `json:"windowHandleName"`
	WindowTimeout    int    `json:"windowTimeout"`
	OpensWindow      bool   `json:"opensWindow"`
}

// ConfigTest represents a test configuration containing a list of commands to be executed as part of the test.
type ConfigTest struct {
	Commands []ConfigCommand `json:"commands"`
}

// Config represents the configuration structure used to manage application settings and operational parameters.
type Config struct {
	Url                     string       `json:"url"`
	Timeout                 *int         `json:"timeout"`
	HumanWait               *int         `json:"humanWait"`
	RetryInterval           *int         `json:"retryInterval"`
	MaxScreenshotLength     *int         `json:"maxScreenshotLength"`
	DownloadPath            string       `json:"downloadPath"`
	Quit                    *bool        `json:"quit"`
	ProfileCapture          []string     `json:"profileCapture"`
	ProfileSupportedMethods []string     `json:"profileSupportedMethods"`
	Tests                   []ConfigTest `json:"tests"`
	ProbeId                 string       `json:"probeId"`
	ThreadName              string       `json:"threadName"`
	Debug                   bool         `json:"debug"`
}
