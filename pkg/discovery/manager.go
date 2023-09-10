package discovery

import (
	"context"
	"fmt"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/collector"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ManagerOpts struct {
	Client  client.Client
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Manager struct {
	client      client.Client
	logger      logr.Logger
	metrics     *strata.Metrics
	discoveries map[types.NamespacedName]*Service
	collectors  *collector.Manager
}

func NewManager(opts *ManagerOpts) *Manager {
	return &Manager{
		client:      opts.Client,
		logger:      opts.Logger,
		metrics:     opts.Metrics,
		discoveries: make(map[types.NamespacedName]*Service),
	}
}

func (m *Manager) WithCollector(collectors *collector.Manager) *Manager {
	m.collectors = collectors

	return m
}

func (m *Manager) Add(ctx context.Context, key types.NamespacedName, spec v1beta1.DiscoverySpec) error {
	if s, ok := m.discoveries[key]; ok {
		s.Stop()
	}

	var list v1beta1.CollectorList
	err := m.client.List(ctx, &list, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(spec.Collection.MatchLabels),
	})
	if err != nil {
		return err
	}

	if len(list.Items) == 0 {
		return fmt.Errorf("no collector found for discovery")
	} else if len(list.Items) > 1 {
		return fmt.Errorf("multiple collectors found for discovery, use additional labels to narrow the search to one")
	}

	collector := list.Items[0]
	collectorKey := types.NamespacedName{
		Namespace: collector.Namespace,
		Name:      collector.Name,
	}
	pool, ok := m.collectors.Get(collectorKey)
	if !ok {
		return fmt.Errorf("collector not found for discovery")
	}

	sendChan := pool.SendChan()

	// Query the collector by the selectors to get the channel.

	// check the collector to see if it exists.

	if s, ok := m.discoveries[key]; ok {
		s.Stop()
	}

	svc := NewService(&ServiceOpts{
		Client:          m.client,
		Enabled:         *spec.Enabled,
		IntervalSeconds: *spec.IntervalSeconds,
		Selector:        spec.Selector,
		Logger:          m.logger.WithValues("discovery", key, "collector", collectorKey.String()),
		Metrics:         m.metrics,
	})
	m.discoveries[key] = svc
	svc.Start(sendChan)

	return nil
}
