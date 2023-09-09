package collectors

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CollectorOpts struct {
	Name     string
	Enabled  bool
	Workers  int
	Client   client.Client
	Log      logr.Logger
	Interval int64
}

type Collector struct {
	name      string
	cancel    context.CancelFunc
	enabled   bool
	client    client.Client
	log       logr.Logger
	interval  int64
	startChan chan error
	stopChan  chan struct{}
	stopOnce  sync.Once
	workers   int
	sync.Mutex
}

// NewCollector returns a new collection resource
func NewCollector(opts *CollectorOpts) *Collector {
	// No need to default the opts as would be the usual step,
	// since we are already going be defaulting them.
	return &Collector{
		enabled:   opts.Enabled,
		client:    opts.Client,
		interval:  opts.Interval,
		log:       opts.Log,
		workers:   opts.Workers,
		startChan: make(chan error),
		stopChan:  make(chan struct{}),
	}
}

// Start begins the collection process.
func (c *Collector) Start(ctx context.Context) <-chan error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	workChan := make(chan Resource, c.workers)

	go func() {
		c.startDiscovery(ctx, workChan)
		c.startWorkers(ctx)
		c.startStatus(ctx)
	}()

	return c.startChan
}

func (c *Collector) startDiscovery(ctx context.Context, work chan Resource) error {
	d := NewDiscovery(&DiscoveryOpts{
		Client:   c.client,
		Interval: time.Duration(c.interval) * time.Second,
		Log:      c.log.WithValues("name", c.name),
	})

	if err := <-d.Start(ctx, work); err != nil {
		return err
	}

	return nil
}

func (c *Collector) startWorkers(ctx context.Context) {
	for i := 0; i < c.workers; i++ {
		// start workers
	}
}

func (c *Collector) startStatus(ctx context.Context) {

}

// Stop ends the collection process.
func (c *Collector) Stop() {
	c.stopOnce.Do(func() {
		c.cancel()
	})
}
