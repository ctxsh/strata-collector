package collectors

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DiscoveryOpts struct {
	// TODO: May need the labels to merge in
	Client   client.Client
	Log      logr.Logger
	Interval time.Duration
}

type Discovery struct {
	client    client.Client
	interval  time.Duration
	log       logr.Logger
	startChan chan error
	stopChan  chan struct{}
	stopOnce  sync.Once
}

func NewDiscovery(opts *DiscoveryOpts) *Discovery {
	return &Discovery{
		client:    opts.Client,
		interval:  opts.Interval,
		log:       opts.Log,
		startChan: make(chan error),
		stopChan:  make(chan struct{}),
	}
}

func (d *Discovery) Start(ctx context.Context, work chan<- Resource) <-chan error {
	go func() {
		d.start(ctx, work)
	}()

	return d.startChan
}

func (d *Discovery) start(ctx context.Context, work chan<- Resource) {
	d.startChan <- d.discover(ctx, work)

	ticker := time.NewTicker(d.interval)
	for {
		select {
		case <-d.stopChan:
			d.log.V(8).Info("worker received stop")
			return
		case <-ctx.Done():
			d.log.V(8).Info("worker is done")
		case <-ticker.C:
			_ = d.discover(ctx, work)
		}
	}
}

func (d *Discovery) discover(ctx context.Context, work chan<- Resource) error {
	d.log.V(6).Info("starting discovery run")
	// discover pods
	// discover services
	// store the metrics for the collector to monitor and update it's status.
	return nil
}

func (d *Discovery) Stop() {
	d.stopOnce.Do(func() {
		close(d.stopChan)
	})
}
