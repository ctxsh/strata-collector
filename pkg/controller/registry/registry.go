package registry

import (
	"context"
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/collector"
	"ctx.sh/strata-collector/pkg/discovery"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RegistryOpts are the options for the registry.
type RegistryOpts struct {
	Cache   cache.Cache
	Client  client.Client
	Logger  logr.Logger
	Metrics *strata.Metrics
}

// Registry is the discovery service and collector retistry.  It is responsible for
// managing the relationship between discovery services and collectors.  It creates
// and manages the discovery services and collectors.  It also manages the channels
// that are used by the both services.
type Registry struct {
	cache       cache.Cache
	client      client.Client
	logger      logr.Logger
	metrics     *strata.Metrics
	discoveries map[types.NamespacedName]*discovery.Service
	collectors  map[types.NamespacedName]collector.Collector

	sync.Mutex
}

func New(mgr ctrl.Manager, opts *RegistryOpts) *Registry {
	return &Registry{
		cache:       opts.Cache,
		client:      opts.Client,
		logger:      opts.Logger,
		metrics:     opts.Metrics,
		discoveries: make(map[types.NamespacedName]*discovery.Service),
		collectors:  make(map[types.NamespacedName]collector.Collector),
	}
}

// TODO: don't need key
func (r *Registry) AddDiscoveryService(ctx context.Context, key types.NamespacedName, obj *v1beta1.Discovery) error {
	r.Lock()
	defer r.Unlock()

	// Check to see if we already have a discovery service for this key and if so, stop it.
	if s, ok := r.discoveries[key]; ok {
		r.logger.V(8).Info("updating existing discovery service")
		s.Stop()
	}

	svc := discovery.NewService(obj, &discovery.ServiceOpts{
		Cache:      r.cache,
		Client:     r.client,
		Logger:     r.logger.WithValues("discovery", key),
		Collectors: r.collectors,
	})

	r.discoveries[key] = svc
	svc.Start()

	return nil
}

func (r *Registry) DeleteDiscoveryService(key types.NamespacedName) error {
	r.Lock()
	defer r.Unlock()

	if o, ok := r.discoveries[key]; ok {
		o.Stop()
		delete(r.discoveries, key)
	}

	return nil
}

func (r *Registry) GetDiscoveryService(key types.NamespacedName) (o *discovery.Service, ok bool) {
	o, ok = r.discoveries[key]
	return
}

func (r *Registry) AddCollectionPool(ctx context.Context, key types.NamespacedName, obj v1beta1.Collector) error {
	r.Lock()
	defer r.Unlock()

	// Check to see if we already have a collector for this key and if so, stop it.
	if c, ok := r.collectors[key]; ok {
		r.logger.V(8).Info("updating existing collector")
		c.Stop()
	}

	collector := collector.NewPool(obj.Namespace, obj.Name, &collector.PoolOpts{
		NumWorkers: *obj.Spec.Workers,
		Logger:     r.logger.WithValues("collector", key),
		Metrics:    r.metrics,
	})

	r.collectors[key] = collector
	collector.Start()

	return nil
}

func (r *Registry) DeleteCollectionPool(key types.NamespacedName) error {
	r.Lock()
	defer r.Unlock()

	if c, ok := r.collectors[key]; ok {
		c.Stop()
		delete(r.collectors, key)
	}

	return nil
}

func (r *Registry) GetCollectionPool(key types.NamespacedName) (o collector.Collector, ok bool) {
	o, ok = r.collectors[key]
	return
}

// TODO: Notifiers and registries.

// RegisterToCollector registers a discovery service to a collector.
// SendChan returns the work channel for the discovery service.
// NewProcessingChan creates and returns a new processing channel.
// ProcessingChanFrom(types.NamespacedName) returns the processing channel for a collector.

// I should keep the channel creation here.  This way it is independent of the
// the collector.  The collector can be restarted without affecting the discovery
// service.  Think about this a bit more since that would impact shutdown synchronization
// as well.

// Allow for a discard channel.  This would be used when a collector has been removed or
// is not present.  Status would show that the discovery service is discarding data.

// Consider moving the mangers back up here for a single place to manage processes.
