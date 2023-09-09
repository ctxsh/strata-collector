package collectors

import (
	"errors"

	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

type Manager struct {
	collectors map[types.NamespacedName]*Collector
}

func New() *Manager {
	return &Manager{
		collectors: make(map[types.NamespacedName]*Collector),
	}
}

func (m *Manager) Add(key types.NamespacedName, spec v1beta1.CollectorSpec) error {
	// TODO: better validation inside of the collector.
	m.collectors[key] = &Collector{
		enabled: *spec.Enabled,
	}

	// start the collector and return any errors from the startup
	return nil
}

func (m *Manager) Delete(key types.NamespacedName) error {
	if _, ok := m.collectors[key]; ok {
		delete(m.collectors, key)
		return nil
	}

	// TODO: better error?
	return errors.New("not found")
}
