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

package discovery

import (
	"context"
	"time"

	"ctx.sh/strata-collector/pkg/service"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler struct {
	client   client.Client
	log      logr.Logger
	observed Observed
	recorder record.EventRecorder
	services *service.Manager
}

var requeueResult reconcile.Result = ctrl.Result{
	Requeue:      true,
	RequeueAfter: 30 * time.Second,
}

func (r *Reconciler) reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	r.log.V(8).Info("request received", "request", request)

	if r.observed.discovery == nil {
		r.log.V(8).Info("discovery service not found, deleting registry entry")
		err := r.services.DeleteDiscoveryService(request.NamespacedName)
		if err != nil {
			r.log.Error(err, "unable to delete discovery service")
			// TODO: Need to think through the conditions here and whether we should requeue the request.
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	err := r.services.AddDiscoveryService(ctx, r.observed.discovery.DeepCopy())
	if err != nil {
		r.log.Error(err, "unable to add discovery service")
		return requeueResult, err
	}

	r.log.V(8).Info("reconciling discovery")
	return ctrl.Result{}, nil
}
