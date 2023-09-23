package service

import (
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/encoder"
	"ctx.sh/strata-collector/pkg/output"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

type CollectionPoolOpts struct {
	Discard  bool
	Logger   logr.Logger
	Metrics  *strata.Metrics
	Registry *Registry
}

type CollectionPool struct {
	namespacedName types.NamespacedName
	numWorkers     int64
	encoder        encoder.Encoder
	workers        []*CollectionWorker
	registry       *Registry
	logger         logr.Logger
	metrics        *strata.Metrics
	output         output.Output

	stopOnce sync.Once
	sync.Mutex
}

func NewCollectionPool(obj v1beta1.Collector, opts *CollectionPoolOpts) *CollectionPool {
	return &CollectionPool{
		namespacedName: types.NamespacedName{
			Namespace: obj.GetNamespace(),
			Name:      obj.GetName(),
		},
		registry:   opts.Registry,
		output:     OutputFactory(*obj.Spec.Output),
		encoder:    EncoderFactory(*obj.Spec.Encoder),
		numWorkers: *obj.Spec.Workers,
		workers:    make([]*CollectionWorker, *obj.Spec.Workers),
		logger:     opts.Logger,
		metrics:    opts.Metrics,
	}
}

func (p *CollectionPool) Start(ch <-chan resource.Resource) {
	for i := int64(0); i < p.numWorkers; i++ {
		p.workers[i] = NewCollectionWorker(&CollectionWorkerOpts{
			Logger:  p.logger.WithValues("worker", i),
			Output:  p.output,
			Encoder: p.encoder,
		})
		p.workers[i].Start(ch)
	}
}

func (p *CollectionPool) Stop() {
	p.logger.V(8).Info("stopping collection pool")
	p.stopOnce.Do(func() {
		// Do I need something here?  Channels will now be closed by the
		// registry.
	})
}

func (p *CollectionPool) NamespacedName() types.NamespacedName {
	return p.namespacedName
}

var _ Collector = &CollectionPool{}
