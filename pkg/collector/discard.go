package collector

import (
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
)

type DiscardOpts struct {
	NumCollectors int64
	Discard       bool
	Logger        logr.Logger
	Metrics       *strata.Metrics
}

type Discard struct {
	recvChan chan resource.Resource
	logger   logr.Logger
	metrics  *strata.Metrics

	stopOnce sync.Once
}

func NewDiscard(opts *DiscardOpts) *Discard {
	return &Discard{
		logger:   opts.Logger,
		metrics:  opts.Metrics,
		recvChan: make(chan resource.Resource),
	}
}

func (d *Discard) Start() {
	go d.start()
}

func (d *Discard) start() {
	for r := range d.recvChan {
		// metrics and discard
		d.logger.V(8).Info("discarding resource", "resource", r)
	}
}

func (d *Discard) Stop() {
	d.stopOnce.Do(func() {
		close(d.recvChan)
	})
}

func (d *Discard) SendChan() chan<- resource.Resource {
	return d.recvChan
}
