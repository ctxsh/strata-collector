package collector

import (
	"context"

	"ctx.sh/strata-collector/pkg/controller/registry"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client   client.Client
	log      logr.Logger
	observed Observed
	recorder record.EventRecorder
	registry *registry.Registry
}

// var requeueResult reconcile.Result = ctrl.Result{
// 	Requeue:      true,
// 	RequeueAfter: 30 * time.Second,
// }

func (r *Reconciler) reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	if r.observed.collector == nil {
		return ctrl.Result{}, nil
	}

	r.log.V(8).Info("reconciling collector")
	return ctrl.Result{}, nil
}
