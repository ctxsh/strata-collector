package controller

import (
	"ctx.sh/strata"
	"ctx.sh/strata-collector/pkg/collector"
	"ctx.sh/strata-collector/pkg/discovery"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type RegistryOpts struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Registry struct {
	Collections *collector.Manager
	Discoveries *discovery.Manager

	// processChan map[types.NamespacedName]chan collector.Resource
}

func NewRegistry(mgr ctrl.Manager, opts *RegistryOpts) *Registry {
	return &Registry{
		Collections: collector.NewManager(&collector.ManagerOpts{
			Logger:  opts.Logger,
			Metrics: opts.Metrics,
		}),
		Discoveries: discovery.NewManager(&discovery.ManagerOpts{
			Logger:  opts.Logger,
			Metrics: opts.Metrics,
			Client:  mgr.GetClient(),
		}),
	}
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
