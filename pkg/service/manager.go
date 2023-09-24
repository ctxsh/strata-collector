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

	return m.registry.AddCollectionPool(key, collector, *obj.Spec.BufferSize)
}

func (m *Manager) DeleteCollectionPool(key types.NamespacedName) error {
	return m.registry.DeleteCollectionPool(key)
}
