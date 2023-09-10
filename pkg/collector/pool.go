package collector

import (
	"sync"
)

type PoolOpts struct {
	NumCollectors int64
}

type Pool struct {
	numCollectors int64
	workers       []*Worker
	recvChan      chan Resource

	stopOnce sync.Once
}

func NewPool(opts *PoolOpts) *Pool {
	return &Pool{
		numCollectors: opts.NumCollectors,
		workers:       make([]*Worker, opts.NumCollectors),
		// TODO: make this configurable
		recvChan: make(chan Resource, 10000),
	}
}

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

func (p *Pool) SendChan() chan<- Resource {
	return p.recvChan
}
