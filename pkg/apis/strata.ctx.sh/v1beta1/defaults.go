package v1beta1

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// DefaultDiscoveryPrefix is the default prefix for all resources.
	DefaultDiscoveryPrefix string = "prometheus.io"
	// DefaultDiscoveryIntervalSeconds is the default interval in seconds that the discovery
	DefaultDiscoveryIntervalSeconds int64 = 10
	// DefaultDiscoveryIncludeMetadata is the default value for including metadata.
	DefaultDiscoveryIncludeMetadata bool = false
	// DefaultDiscoveryEnabled is the default value for enabling the discovery service.
	DefaultDiscoveryEnabled bool = true
)

var (
	DefaultDiscoveryResources = []string{
		"pods",
		"services",
		"endpoints",
	}
)

// Defaulted sets the resource defaults.
func Defaulted(obj client.Object) {
	switch obj := obj.(type) {
	case *Collector:
		defautledCollector(obj)
	case *Discovery:
		defautledDiscovery(obj)
	}
}

func defautledCollector(obj *Collector) {}

func defautledDiscovery(obj *Discovery) {
	if obj.Spec.Enabled == nil {
		enabled := DefaultDiscoveryEnabled
		obj.Spec.Enabled = &enabled
	}

	if obj.Spec.IntervalSeconds == nil {
		interval := DefaultDiscoveryIntervalSeconds
		obj.Spec.IntervalSeconds = &interval
	}

	if obj.Spec.Prefix == nil {
		prefix := DefaultDiscoveryPrefix
		obj.Spec.Prefix = &prefix
	}

	if obj.Spec.IncludeMetadata == nil {
		includeMetadata := DefaultDiscoveryIncludeMetadata
		obj.Spec.IncludeMetadata = &includeMetadata
	}

	if obj.Spec.Resources == nil {
		resources := DefaultDiscoveryResources
		obj.Spec.Resources = resources
	}
}
