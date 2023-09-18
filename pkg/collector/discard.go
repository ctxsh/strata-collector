package collector

import (
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
)

// DiscardOpts are the options for the discard collector.
type DiscardOpts struct {
	Discard bool
	Logger  logr.Logger
	Metrics *strata.Metrics
}

// Discard is a collector that discards all resources.  It's primary purpose is
// to provide a collector that will be used when either: a collector is not found,
// the collector is disabled, or there is pressure on the confifured collector pool
// and the metrics need to be discarded.  This collector is always enabled and tracks
// the number of discarded resources.
type Discard struct {
	recvChan chan resource.Resource
	logger   logr.Logger
	metrics  *strata.Metrics
	stopOnce sync.Once
}

// NewDiscard creates a new discard collector.
func NewDiscard(opts *DiscardOpts) *Discard {
	return &Discard{
		logger:   opts.Logger,
		metrics:  opts.Metrics,
		recvChan: make(chan resource.Resource),
	}
}

// Start starts the discard collector.
func (d *Discard) Start() {
	go d.start()
}

// start is the main loop for the discard collector.
func (d *Discard) start() {
	for r := range d.recvChan {
		// metrics and discard
		d.logger.V(8).Info("discarding resource", "resource", r)
	}
}

// Stop stops the discard collector.
func (d *Discard) Stop() {
	d.stopOnce.Do(func() {
		close(d.recvChan)
	})
}

// SendChan implements the Collector interface.  It returns the channel that
// resources should be sent to.
func (d *Discard) SendChan() chan<- resource.Resource {
	return d.recvChan
}
