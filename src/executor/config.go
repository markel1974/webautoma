package executor

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

type ConfigTest struct {
	Commands []ConfigCommand `json:"commands"`
}

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
	Debug                   bool         `json:"debug"`
}
