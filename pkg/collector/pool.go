package collector

import (
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

type PoolOpts struct {
	NumWorkers int64
	Discard    bool
	Logger     logr.Logger
	Metrics    *strata.Metrics
}

type Pool struct {
	namespacedName types.NamespacedName
	numWorkers     int64
	workers        []*Worker
	recvChan       chan resource.Resource
	logger         logr.Logger
	metrics        *strata.Metrics

	discard bool

	stopOnce sync.Once
}

func NewPool(namespace, name string, opts *PoolOpts) *Pool {
	return &Pool{
		namespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
		discard:    opts.Discard,
		numWorkers: opts.NumWorkers,
		workers:    make([]*Worker, opts.NumWorkers),
		logger:     opts.Logger,
		metrics:    opts.Metrics,
		// TODO: make this configurable
		recvChan: make(chan resource.Resource, 10000),
	}
}

func (p *Pool) Start() {
	for i := int64(0); i < p.numWorkers; i++ {
		p.workers[i] = NewWorker(&WorkerOpts{
			Logger: p.logger.WithValues("worker", i),
		})
	}
}

func (p *Pool) Stop() {
	p.stopOnce.Do(func() {
		close(p.recvChan)
	})
}

func (p *Pool) SendChan() chan<- resource.Resource {
	return p.recvChan
}

func (p *Pool) NamespacedName() types.NamespacedName {
	return p.namespacedName
}
