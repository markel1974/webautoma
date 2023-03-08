package executor

type ConfigCommand struct {
	Id      string `json:"id"`
	Target  string `json:"target"`
	Command string `json:"command"`
	Until   string `json:"until"`
	Value   string `json:"value"`
}

type ConfigTest struct {
	Commands []ConfigCommand `json:"commands"`
}

type Config struct {
	Url     string       `json:"url"`
	Timeout *int         `json:"timeout"`
	Quit    *bool        `json:"quit"`
	Tests   []ConfigTest `json:"tests"`
}
