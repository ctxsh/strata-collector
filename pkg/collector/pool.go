package collector

import (
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

type PoolOpts struct {
	NumCollectors int64
	Discard       bool
	Logger        logr.Logger
	Metrics       *strata.Metrics
}

type Pool struct {
	namespacedName types.NamespacedName
	numCollectors  int64
	workers        []*Worker
	recvChan       chan resource.Resource

	discard bool

	stopOnce sync.Once
}

func NewPool(namespace, name string, opts *PoolOpts) *Pool {
	return &Pool{
		namespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
		discard:       opts.Discard,
		numCollectors: opts.NumCollectors,
		workers:       make([]*Worker, opts.NumCollectors),
		// TODO: make this configurable
		recvChan: make(chan resource.Resource, 10000),
	}
}

// TODO: send the channel.
func (p *Pool) Start() {
	for i := int64(0); i < p.numCollectors; i++ {
		p.workers[i] = NewWorker(&WorkerOpts{})
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
