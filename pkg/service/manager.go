package service

import (
	"context"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ManagerOpts struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Manager struct {
	registry *Registry

	logger  logr.Logger
	metrics *strata.Metrics

	cache  cache.Cache
	client client.Client
}

func NewManager(mgr ctrl.Manager, opts *ManagerOpts) *Manager {
	return &Manager{
		registry: NewRegistry(),
		logger:   opts.Logger,
		metrics:  opts.Metrics,
		cache:    mgr.GetCache(),
		client:   mgr.GetClient(),
	}
}

// TODO: don't need key
func (m *Manager) AddDiscoveryService(ctx context.Context, obj *v1beta1.Discovery) error {
	key := types.NamespacedName{
		Namespace: obj.Namespace,
		Name:      obj.Name,
	}

	svc := NewDiscovery(obj, &DiscoveryOpts{
		Cache:    m.cache,
		Client:   m.client,
		Logger:   m.logger.WithValues("discovery", key),
		Registry: m.registry,
	})

	return m.registry.AddDiscoveryService(key, svc)
}

func (m *Manager) DeleteDiscoveryService(key types.NamespacedName) error {
	return m.registry.DeleteDiscoveryService(key)
}

func (m *Manager) AddCollectionPool(ctx context.Context, key types.NamespacedName, obj v1beta1.Collector) error {
	collector := NewCollectionPool(obj, &CollectionPoolOpts{
		Logger:   m.logger.WithValues("collector", key),
		Metrics:  m.metrics,
		Registry: m.registry,
	})

	return m.registry.AddCollectionPool(key, collector)
}

func (m *Manager) DeleteCollectionPool(key types.NamespacedName) error {
	return m.registry.DeleteCollectionPool(key)
}
