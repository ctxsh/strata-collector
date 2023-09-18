package collector

import (
	"net/http"
	"time"

	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
)

const (
	DefaultTimeout time.Duration = 2 * time.Second
)

type WorkerOpts struct {
	Logger logr.Logger
}

type Worker struct {
	httpClient http.Client
	logger     logr.Logger
}

func NewWorker(opt *WorkerOpts) *Worker {
	return &Worker{
		// TODO: better client creation and config
		httpClient: http.Client{
			Timeout: DefaultTimeout,
		},
		logger: opt.Logger,
	}
}

func (w *Worker) Start(recvChan chan *resource.Resource) {
	go w.start(recvChan)
}

func (w *Worker) start(recvChan chan *resource.Resource) {
	for r := range recvChan {
		w.collectAndSend(r)
	}

	w.logger.V(8).Info("worker shutting down")
}

func (w *Worker) collectAndSend(r *resource.Resource) {
	w.logger.V(8).Info("collecting resource", "resource", r)
	if err := w.collect(); err != nil {
		w.logger.Error(err, "failed to collect resource", "resource", r)
		return
	}

	if err := w.send(); err != nil {
		w.logger.Error(err, "failed to send resource", "resource", r)
		return
	}
}

func (w *Worker) collect() error {
	return nil
}

func (w *Worker) send() error {
	return nil
}
