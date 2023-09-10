package collector

import (
	"errors"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

type ManagerOpts struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Manager struct {
	logger  logr.Logger
	metrics *strata.Metrics
	pools   map[types.NamespacedName]*Pool
}

func NewManager(opts *ManagerOpts) *Manager {
	return &Manager{
		logger:  opts.Logger,
		metrics: opts.Metrics,
		pools:   make(map[types.NamespacedName]*Pool),
	}
}

func (m *Manager) Add(key types.NamespacedName, spec *v1beta1.CollectorSpec) error {
	return nil
}

func (m *Manager) Delete(key types.NamespacedName) error {
	if p, ok := m.pools[key]; ok {
		p.Stop()
		delete(m.pools, key)
		return nil
	}

	return errors.New("not found")
}

func (m *Manager) Get(key types.NamespacedName) (p *Pool, ok bool) {
	p, ok = m.pools[key]
	return
}
