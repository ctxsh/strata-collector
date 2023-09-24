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

package collector

import (
	"context"

	v1beta1 "ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Observed struct {
	collector *v1beta1.Collector
	// observeTime time.Time
}

type Observer struct {
	Client  client.Client
	Request ctrl.Request
	Context context.Context
}

func (o *Observer) observe(observed *Observed) error {
	observedCollector := new(v1beta1.Collector)
	err := o.observeCollector(o.Request.NamespacedName, observedCollector)
	if err != nil {
		return client.IgnoreNotFound(err)
	}

	v1beta1.Defaulted(observedCollector)

	observed.collector = observedCollector
	return nil
}

func (o *Observer) observeCollector(key types.NamespacedName, collector *v1beta1.Collector) error {
	return o.Client.Get(o.Context, key, collector)
}
