// Copyright 2023 Rob Lyon <rob@ctxswitch.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	channels    map[types.NamespacedName]chan resource.Resource

	sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		discoveries: make(map[types.NamespacedName]*Discovery),
		collectors:  make(map[types.NamespacedName]Collector),
		channels:    make(map[types.NamespacedName]chan resource.Resource),
	}
}

func (r *Registry) GetCollectionPool(key types.NamespacedName) (o Collector, ok bool) {
	r.RLock()
	defer r.RUnlock()

	o, ok = r.collectors[key]
	return
}

func (r *Registry) AddCollectionPool(key types.NamespacedName, obj Collector, bufferSize int64) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.channels[key]; !ok {
		r.channels[key] = make(chan resource.Resource, bufferSize)
	}

	if o, ok := r.collectors[key]; ok {
		o.Stop()
		delete(r.collectors, key)
	}

	obj.Start(r.channels[key])
	r.collectors[key] = obj
	return nil
}

func (r *Registry) DeleteCollectionPool(key types.NamespacedName) error {
	r.Lock()
	defer r.Unlock()

	ch, ok := r.channels[key]
	if !ok {
		return fmt.Errorf("unable to delete, channel for collection not found for %s", key)
	}

	close(ch)
	delete(r.channels, key)

	co, ok := r.collectors[key]
	if !ok {
		return fmt.Errorf("unable to delete, collection pool not found for %s", key)
	}

	co.Stop()
	delete(r.collectors, key)

	return nil
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

func (r *Registry) RegisteredWithCollector(nn types.NamespacedName) int64 {
	r.RLock()
	defer r.RUnlock()

	var count int64
	for _, o := range r.discoveries {
		if o.HasCollector(nn) {
			count++
		}
	}

	return count
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

	c, ok := r.channels[key]
	if !ok {
		return fmt.Errorf("channel for collection not found for %s", key)
	}

	for i := 0; i < len(resources); i++ {
		if c != nil {
			c <- resources[i]
		}
	}

	return nil
}

// GetInFlightResources returns the number of resources that are currently awaiting processing
// in the send channel.
func (r *Registry) GetInFlightResources(key types.NamespacedName) (int64, error) {
	r.RLock()
	defer r.RUnlock()

	c, ok := r.channels[key]
	if !ok {
		return 0, fmt.Errorf("collection not found for %s", key)
	}

	return int64(len(c)), nil
}

// GetChannel returns the send channel for the collector.
func (r *Registry) GetChannel(key types.NamespacedName) (chan resource.Resource, error) {
	r.RLock()
	defer r.RUnlock()

	c, ok := r.channels[key]
	if !ok {
		return nil, fmt.Errorf("channel for collection not found for %s", key)
	}

	return c, nil
}

// func (r *Registry) CloseChannel(key types.NamespacedName) error {
// 	r.Lock()
// 	defer r.Unlock()

// 	c, ok := r.channels[key]
// 	if !ok {
// 		return fmt.Errorf("channel for collection not found for %s", key)
// 	}

// 	close(c)
// 	delete(r.channels, key)
// 	return nil
// }
