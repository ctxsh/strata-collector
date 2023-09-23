package stdout

import (
	"fmt"

	"ctx.sh/strata-collector/pkg/output"
)

type Stdout struct {
}

func New() *Stdout {
	return &Stdout{}
}

func (s *Stdout) Connect() error {
	return nil
}

func (s *Stdout) Send(data []byte) error {
	fmt.Println(string(data))
	return nil
}

func (s *Stdout) Close() {
}

var _ output.Output = &Stdout{}
