package service

import (
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

type CollectionPoolOpts struct {
	NumWorkers int64
	Discard    bool
	Logger     logr.Logger
	Metrics    *strata.Metrics
}

type CollectionPool struct {
	namespacedName types.NamespacedName
	numWorkers     int64
	workers        []*CollectionWorker
	recvChan       chan resource.Resource
	logger         logr.Logger
	metrics        *strata.Metrics

	discard bool

	stopOnce sync.Once
	sync.Mutex
}

func NewCollectionPool(namespace, name string, opts *CollectionPoolOpts) *CollectionPool {
	return &CollectionPool{
		namespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
		discard:    opts.Discard,
		numWorkers: opts.NumWorkers,
		workers:    make([]*CollectionWorker, opts.NumWorkers),
		logger:     opts.Logger,
		metrics:    opts.Metrics,
		// TODO: make this configurable
		recvChan: make(chan resource.Resource, 10000),
	}
}

func (p *CollectionPool) Start() {
	for i := int64(0); i < p.numWorkers; i++ {
		p.workers[i] = NewCollectionWorker(&CollectionWorkerOpts{
			Logger: p.logger.WithValues("worker", i),
		})
		p.workers[i].Start(p.recvChan)
	}
}

func (p *CollectionPool) Stop() {
	p.stopOnce.Do(func() {
		close(p.recvChan)
	})
}

func (p *CollectionPool) SendChan() chan<- resource.Resource {
	return p.recvChan
}

func (p *CollectionPool) NamespacedName() types.NamespacedName {
	return p.namespacedName
}

var _ Collector = &CollectionPool{}
