// Copyright 2023 Rob Lyon <rob@ctxswitch.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"ctx.sh/strata-collector/pkg/encoder"
	"ctx.sh/strata-collector/pkg/filter"
	"ctx.sh/strata-collector/pkg/metric"
	"ctx.sh/strata-collector/pkg/output"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
)

const (
	DefaultTimeout time.Duration = 2 * time.Second
)

type CollectionWorkerOpts struct {
	Logger  logr.Logger
	Output  output.Output
	Encoder encoder.Encoder
	Filters *filter.Filter
}

type CollectionWorker struct {
	httpClient http.Client
	output     output.Output
	logger     logr.Logger
	encoder    encoder.Encoder
	filters    *filter.Filter
}

func NewCollectionWorker(opts *CollectionWorkerOpts) *CollectionWorker {
	return &CollectionWorker{
		// TODO: better client creation and config
		httpClient: http.Client{
			Timeout: DefaultTimeout,
		},
		encoder: opts.Encoder,
		output:  opts.Output,
		logger:  opts.Logger,
		filters: opts.Filters,
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

func (w *CollectionWorker) collect(r resource.Resource) ([]*metric.Metric, error) {
	// TODO: Add timeout, scheme, and other options - this is just a quick
	// implementation to get things working.
	url := fmt.Sprintf("%s://%s:%s%s", r.Scheme, r.IP, r.Port, r.Path)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// TODO: determine tags and fields from our config and use the resource struct
	// to actually return them.
	m, err := metric.FromPrometheusMetric(time.Now(), buf)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (w *CollectionWorker) send(metrics []*metric.Metric) error {

	for _, m := range metrics {
		if w.filters.Do(m) {
			continue
		}

		data, err := w.encoder.Encode(m)
		if err != nil {
			// TODO: collect errors and send them at the end.
			w.logger.Error(err, "failed to encode metric", "metric", m)
			continue
		}

		if err = w.output.Send(data); err != nil {
			// TODO: collect errors and send them at the end.
			return err
		}
	}
	return nil
}
