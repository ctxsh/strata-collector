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
	"context"
	"sync"
	"time"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/encoder"
	"ctx.sh/strata-collector/pkg/filter"
	"ctx.sh/strata-collector/pkg/output"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CollectionPoolOpts struct {
	Cache    cache.Cache
	Client   client.Client
	Discard  bool
	Logger   logr.Logger
	Metrics  *strata.Metrics
	Registry *Registry
}

type CollectionPool struct {
	name       string
	namespace  string
	cache      cache.Cache
	client     client.Client
	numWorkers int64
	encoder    encoder.Encoder
	filters    *filter.Filter
	workers    []*CollectionWorker
	registry   *Registry
	logger     logr.Logger
	metrics    *strata.Metrics
	output     output.Output
	stats      *CollectionStats
	obj        *v1beta1.Collector

	stopChan chan struct{}
	stopOnce sync.Once
	sync.Mutex
}

func NewCollectionPool(obj *v1beta1.Collector, opts *CollectionPoolOpts) *CollectionPool {
	return &CollectionPool{
		name:       obj.GetName(),
		namespace:  obj.GetNamespace(),
		client:     opts.Client,
		cache:      opts.Cache,
		registry:   opts.Registry,
		output:     OutputFactory(obj.Spec.Output),
		encoder:    EncoderFactory(*obj.Spec.Encoder),
		obj:        obj,
		filters:    FilterFactory(obj.Spec.Filters),
		numWorkers: *obj.Spec.Workers,
		workers:    make([]*CollectionWorker, *obj.Spec.Workers),
		logger:     opts.Logger,
		metrics:    opts.Metrics,
		stats:      NewCollectionStats(),
		stopChan:   make(chan struct{}),
	}
}

func (p *CollectionPool) Start(ch <-chan resource.Resource) {
	ctx := context.Background()

	for i := int64(0); i < p.numWorkers; i++ {
		p.workers[i] = NewCollectionWorker(&CollectionWorkerOpts{
			Logger:  p.logger.WithValues("worker", i),
			Output:  p.output,
			Encoder: p.encoder,
			Filters: p.filters,
			Stats:   p.stats,
		})
		p.workers[i].Start(ch)
	}

	go p.status(ctx)
}

func (p *CollectionPool) Stop() {
	p.logger.V(8).Info("stopping collection pool")
	p.stopOnce.Do(func() {
		close(p.stopChan)
	})
}

func (p *CollectionPool) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: p.namespace,
		Name:      p.name,
	}
}

func (p *CollectionPool) status(ctx context.Context) {
	ticker := time.NewTicker(DefaultStatusInterval)
	for {
		select {
		case <-p.stopChan:
			return
		case <-ticker.C:
			err := p.updateStatus(ctx)
			if err != nil {
				p.logger.Error(err, "unable to update status")
			}
		}
	}
}

func (p *CollectionPool) updateStatus(ctx context.Context) error {
	p.Lock()
	defer p.Unlock()

	var obj v1beta1.Collector
	err := p.cache.Get(ctx, p.NamespacedName(), &obj)
	if err != nil {
		return err
	}

	inFlight, err := p.registry.GetInFlightResources(p.NamespacedName())
	if err != nil {
		return err
	}

	obj.Status = v1beta1.CollectorStatus{
		RegisteredDiscoveries: p.registry.RegisteredWithCollector(p.NamespacedName()),
		InFlightResources:     int64(inFlight),
		TotalSent:             p.stats.TotalSent.Load(),
		TotalErrors:           p.stats.TotalErrors.Load(),
		TotalFiltered:         p.stats.TotalFiltered.Load(),
		MetricsCollected:      p.stats.MetricsCollected.Load(),
	}

	p.logger.V(8).Info("updating collector status", "status", obj.Status)

	// Not sure if I want to reset here. The numbers are going to get quite
	// large, but if we reset, there's a period of time on lower volume installs
	// where everything will be zeroed out for a time.  It's not the best UX.
	// p.stats.Reset()
	return p.client.Status().Update(ctx, &obj)
}

var _ Collector = &CollectionPool{}
