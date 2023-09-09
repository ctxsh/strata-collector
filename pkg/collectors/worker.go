package collectors

import (
	"context"
	"net/http"
	"sync"
	"time"

	"ctx.sh/strata-collector/pkg/metrics"
	"github.com/go-logr/logr"
)

const (
	DefaultTimeout time.Duration = 2 * time.Second
)

type WorkerOpts struct {
	Log logr.Logger
}

type Worker struct {
	httpClient http.Client
	log        logr.Logger
	startChan  chan error
	stopChan   chan struct{}
	stopOnce   sync.Once
}

func NewWorker(opt *WorkerOpts) *Worker {
	// TODO: better client creation and config
	return &Worker{
		httpClient: http.Client{
			Timeout: DefaultTimeout,
		},
		log:       opt.Log,
		startChan: make(chan error),
		stopChan:  make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context, work <-chan Resource, sink chan<- metrics.Metric) <-chan error {
	go func() {
		w.start(ctx, work, sink)
	}()

	return w.startChan
}

func (w *Worker) start(ctx context.Context, work <-chan Resource, sink chan<- metrics.Metric) {
	for {
		select {
		case <-w.stopChan:
			w.log.V(8).Info("worker received stop")
			return
		case <-ctx.Done():
			w.log.V(8).Info("worker is done")
		case r := <-work:
			m, err := w.collect(ctx, r)
			if err != nil {
				w.log.Error(err, "collection failed")
				continue
			}
			w.send(ctx, sink, m)
		}
	}
}

func (w *Worker) collect(ctx context.Context, res Resource) (metrics.Metric, error) {
	return metrics.Metric{}, nil
}

func (w *Worker) send(ctx context.Context, sink chan<- metrics.Metric, m metrics.Metric) {
	// TODO: deal with timeout so we don't end up infinitely blocking
	// TODO: check len and drop
	sink <- m
}

func (w *Worker) Stop() {
	w.stopOnce.Do(func() {
		close(w.stopChan)
	})
}
