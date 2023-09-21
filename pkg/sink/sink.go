package sink

type Sink interface {
	Connect() error
	Send(data []byte) error
	Close()
}
