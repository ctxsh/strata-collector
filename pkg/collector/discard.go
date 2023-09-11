package collector

import (
	"sync"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
)

type DiscardOpts struct {
	NumCollectors int64
	Discard       bool
	Logger        logr.Logger
	Metrics       *strata.Metrics
}

type Discard struct {
	recvChan chan Resource

	stopOnce sync.Once
}

func NewDiscard(opts *DiscardOpts) *Discard {
	return &Discard{
		recvChan: make(chan Resource),
	}
}

func (d *Discard) Start() {
	go d.start()
}

func (d *Discard) start() {
	for range d.recvChan {
		// metrics and discard
	}
}

func (d *Discard) Stop() {
	d.stopOnce.Do(func() {
		close(d.recvChan)
	})
}

func (d *Discard) SendChan() chan<- Resource {
	return d.recvChan
}
