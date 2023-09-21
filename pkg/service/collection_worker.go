package service

import (
	"fmt"
	"net/http"
	"time"

	"ctx.sh/strata-collector/pkg/resource"
	"ctx.sh/strata-collector/pkg/sink"
	"github.com/go-logr/logr"
)

const (
	DefaultTimeout time.Duration = 2 * time.Second
)

type CollectionWorkerOpts struct {
	Logger logr.Logger
	Output sink.Sink
}

type CollectionWorker struct {
	httpClient http.Client
	output     sink.Sink
	logger     logr.Logger
}

func NewCollectionWorker(opt *CollectionWorkerOpts) *CollectionWorker {
	return &CollectionWorker{
		// TODO: better client creation and config
		httpClient: http.Client{
			Timeout: DefaultTimeout,
		},
		output: opt.Output,
		logger: opt.Logger,
	}
}

func (w *CollectionWorker) Start(recvChan <-chan resource.Resource) {
	go w.start(recvChan)
}

func (w *CollectionWorker) start(recvChan <-chan resource.Resource) {
	for r := range recvChan {
		w.collectAndSend(r)
	}

	w.logger.V(8).Info("worker shutting down")
}

func (w *CollectionWorker) collectAndSend(r resource.Resource) {
	w.logger.V(8).Info("collecting resource", "resource", r)
	metrics, err := w.collect(r)
	if err != nil {
		w.logger.Error(err, "failed to collect resource", "resource", r)
		return
	}

	err = w.send(metrics)
	if err != nil {
		w.logger.Error(err, "failed to send resource", "resource", r)
		return
	}
}

func (w *CollectionWorker) collect(r resource.Resource) ([]*Metric, error) {
	pm := NewPrometheusScraper(w.httpClient, fmt.Sprintf("%s://%s:%s%s", r.Scheme, r.IP, r.Port, r.Path))

	m, err := pm.Get(map[string]string{})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (w *CollectionWorker) send(metrics []*Metric) error {
	for _, m := range metrics {
		err := w.output.Send(m.Bytes())
		if err != nil {
			// TODO: collect errors and send them at the end.
			return err
		}
	}
	return nil
}
