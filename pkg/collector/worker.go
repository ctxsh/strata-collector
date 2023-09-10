package collector

import (
	"net/http"
	"time"

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
}

func NewWorker(opt *WorkerOpts) *Worker {
	return &Worker{
		// TODO: better client creation and config
		httpClient: http.Client{
			Timeout: DefaultTimeout,
		},
		log: opt.Log,
	}
}

func (w *Worker) Start(recvChan chan *Resource) {
	go w.start(recvChan)
}

func (w *Worker) start(recvChan chan *Resource) {
	for r := range recvChan {
		w.collectAndSend(r)
	}

	w.log.V(8).Info("worker shutting down")
}

func (w *Worker) collectAndSend(r *Resource) {
	if err := w.collect(); err != nil {
		w.log.Error(err, "failed to collect resource", "resource", r)
		return
	}

	if err := w.send(); err != nil {
		w.log.Error(err, "failed to send resource", "resource", r)
		return
	}
}

func (w *Worker) collect() error {
	return nil
}

func (w *Worker) send() error {
	return nil
}
