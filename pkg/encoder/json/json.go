package json

import "encoding/json"

// JsonEncoder is an encoder that encodes a generic interface into a JSON
// marshalled byte array.
type JsonEncoder struct{}

func New() *JsonEncoder {
	return &JsonEncoder{}
}

func (e *JsonEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
