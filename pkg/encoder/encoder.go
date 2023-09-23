package encoder

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}
