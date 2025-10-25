package stream

type Message struct {
	Topic   string
	Key     string
	Headers map[string]string
	Payload []byte
}

type Stream interface {
	Start() chan Message
}
