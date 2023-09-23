package nats

import (
	"ctx.sh/strata-collector/pkg/output"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	Subject string
	conn    *nats.Conn
}

func New() *Nats {
	return &Nats{}
}

func (n *Nats) Connect() (err error) {
	n.conn, err = nats.Connect(nats.DefaultURL)
	return
}

func (n *Nats) Send(data []byte) error {
	return n.conn.Publish(n.Subject, data)
}

func (n *Nats) Close() {
	// TODO: Drain? Does this wait for consumers to finish even though we are a publisher?
	// n.conn.Drain()
	n.conn.Close()
}

var _ output.Output = &Nats{}
