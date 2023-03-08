package executor

type NetworkMessage struct {
	Method string `json:"method"`
	Params struct {
		Response struct {
			Url     string                 `json:"url"`
			Timing  interface{}            `json:"timing"`
			Status  interface{}            `json:"status"`
			Headers map[string]interface{} `json:"headers"`
		} `json:"response"`
	} `json:"params"`
}
