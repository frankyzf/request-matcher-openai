package mywebsocket

type wsMessage struct {
	OpType string      `json:"op_type"`
	Data   interface{} `json:"data"`
}
