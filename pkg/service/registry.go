package service

import (
	"fmt"
	"sync"

	"ctx.sh/strata-collector/pkg/resource"
	"k8s.io/apimachinery/pkg/types"
)

// Registry is the discovery service and collector retistry.  It is responsible for
// managing the relationship between discovery services and collectors.  It creates
// and manages the discovery services and collectors.  It also manages the channels
// that are used by the both services.
type Registry struct {
	discoveries map[types.NamespacedName]*Discovery
	collectors  map[types.NamespacedName]Collector

	sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		discoveries: make(map[types.NamespacedName]*Discovery),
		collectors:  make(map[types.NamespacedName]Collector),
	}
}

func (r *Registry) GetCollectionPool(key types.NamespacedName) (o Collector, ok bool) {
	r.RLock()
	defer r.RUnlock()

	o, ok = r.collectors[key]
	return
}

func (r *Registry) AddCollectionPool(key types.NamespacedName, obj Collector) error {
	r.Lock()
	defer r.Unlock()

	if o, ok := r.collectors[key]; ok {
		o.Stop()
		delete(r.collectors, key)
	}

	obj.Start()
	r.collectors[key] = obj
	return nil
}

func (r *Registry) DeleteCollectionPool(key types.NamespacedName) error {
	r.Lock()
	defer r.Unlock()

	if c, ok := r.collectors[key]; ok {
		c.Stop()
		delete(r.collectors, key)
		return nil
	}

	return fmt.Errorf("collection pool not found for %s", key)
}

func (r *Registry) GetDiscoveryService(key types.NamespacedName) (o *Discovery, ok bool) {
	r.RLock()
	defer r.RUnlock()

	o, ok = r.discoveries[key]
	return
}

func (r *Registry) AddDiscoveryService(key types.NamespacedName, obj *Discovery) error {
	r.Lock()
	defer r.Unlock()

	if o, ok := r.discoveries[key]; ok {
		o.Stop()
		delete(r.discoveries, key)
	}

	obj.Start()
	r.discoveries[key] = obj
	return nil
}

func (r *Registry) DeleteDiscoveryService(key types.NamespacedName) error {
	r.Lock()
	defer r.Unlock()

	if o, ok := r.discoveries[key]; ok {
		o.Stop()
		delete(r.discoveries, key)
		return nil
	}

	return fmt.Errorf("discovery service not found for %s", key)
}

// SendResources sends the resources to the collector.
func (r *Registry) SendResources(key types.NamespacedName, resources []resource.Resource) error {
	r.RLock()
	defer r.RUnlock()

	c, ok := r.collectors[key]
	if !ok {
		return fmt.Errorf("collection pool not found for %s", key)
	}

	for i := 0; i < len(resources); i++ {
		sendChan := c.SendChan()
		if sendChan != nil {
			sendChan <- resources[i]
		}
	}

	return nil
}

// GetInFlightResources returns the number of resources that are currently awaiting processing
// in the send channel.
func (r *Registry) GetInFlightResources(key types.NamespacedName) (int, error) {
	r.RLock()
	defer r.RUnlock()

	c, ok := r.collectors[key]
	if !ok {
		return 0, fmt.Errorf("collection not found for %s", key)
	}

	return len(c.SendChan()), nil
}
