package output

type Output interface {
	Connect() error
	Send(data []byte) error
	Close()
}
