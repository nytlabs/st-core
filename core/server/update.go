package server

type Update struct {
	Action string      `json:"action"`
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
}
