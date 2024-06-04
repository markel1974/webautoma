package executor

type ConfigCommand struct {
	Id      string `json:"id"`
	Target  string `json:"target"`
	Command string `json:"command"`
	Until   string `json:"until"`
	Value   string `json:"value"`

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
	Quit                    *bool        `json:"quit"`
	ProfileCapture          []string     `json:"profileCapture"`
	ProfileSupportedMethods []string     `json:"profileSupportedMethods"`
	Tests                   []ConfigTest `json:"tests"`
}
