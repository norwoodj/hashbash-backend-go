package rabbit

import "encoding/json"

type MessageSerializer interface {
	GetContentType() string
	SerializeMessage(interface{}) ([]byte, error)
}

type JsonMessageSerializer struct {}

func (JsonMessageSerializer) GetContentType() string {
	return "application/json"
}

func (JsonMessageSerializer) SerializeMessage(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}
