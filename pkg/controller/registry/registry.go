package registry

import (
	"context"
	"fmt"
	"sync"

	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/collector"
	"ctx.sh/strata-collector/pkg/discovery"
	"ctx.sh/strata-collector/pkg/resource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RegistryOpts struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Registry struct {
	// Collections *collector.Manager
	// Discoveries *discovery.Manager

	client client.Client
	ctx    context.Context

	metrics *strata.Metrics
	logger  logr.Logger

	discard     *collector.Discard
	discoveries map[types.NamespacedName]*discovery.Service
	collectors  map[types.NamespacedName]*collector.Pool
	channels    map[types.NamespacedName]chan<- resource.Resource

	sync.Mutex
}

func New(mgr ctrl.Manager, opts *RegistryOpts) *Registry {
	// Proabaly want to managed the discard collector differently.
	discard := collector.NewDiscard(&collector.DiscardOpts{
		Logger:  opts.Logger.WithValues("collector", "discard"),
		Metrics: opts.Metrics,
	})
	discard.Start()

	return &Registry{
		client: mgr.GetClient(),
		ctx:    context.Background(),

		metrics: opts.Metrics,
		logger:  opts.Logger,

		discard:     discard,
		discoveries: make(map[types.NamespacedName]*discovery.Service),
		collectors:  make(map[types.NamespacedName]*collector.Pool),
		channels:    make(map[types.NamespacedName]chan<- resource.Resource),
	}
}

func (r *Registry) AddDiscoveryService(ctx context.Context, key types.NamespacedName, obj v1beta1.Discovery) error {
	r.Lock()
	defer r.Unlock()

	// Check to see if we already have a discovery service for this key and if so, stop it.
	if s, ok := r.discoveries[key]; ok {
		r.logger.V(8).Info("updating existing discovery service")
		s.Stop()
	}

	svc := discovery.NewService(obj.Namespace, obj.Name, &discovery.ServiceOpts{
		Client:          r.client,
		Enabled:         *obj.Spec.Enabled,
		IntervalSeconds: *obj.Spec.IntervalSeconds,
		Selector:        obj.Spec.Selector,
		Prefix:          *obj.Spec.Prefix,
		Logger:          r.logger.WithValues("discovery", key),
	})

	var sendChan chan<- resource.Resource

	collectorObj, err := r.getCollector(ctx, obj.Spec)
	if err != nil {
		// If we haven't found a collector then it either has not been deployed, or it has been deleted.
		// In either case, the discovery service should start up and discard data until the collector is
		// available.
		r.logger.Info("selected collector was not found, registering discard collector", "discovery", key)
		sendChan = r.discard.SendChan()
	} else {
		collector, ok := r.collectors[namespacedName(&collectorObj)]
		if !ok {
			// This would probably only happen if there is a race between
			// the discovery service and the collector coming on line, however
			// with the locks in place, I don't think that is possible. So instead
			// of registering it to discard, just return an error and requeue.
			return fmt.Errorf("collector not found in the registry")
		}
		sendChan = collector.SendChan()
	}

	r.discoveries[key] = svc
	svc.Start(sendChan)

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
