package discovery

import (
	"context"

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
}

// var requeueResult reconcile.Result = ctrl.Result{
// 	Requeue:      true,
// 	RequeueAfter: 30 * time.Second,
// }

func (r *Reconciler) reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	if r.observed.discovery == nil {
		return ctrl.Result{}, nil
	}

	r.log.V(8).Info("reconciling discovery")
	return ctrl.Result{}, nil
}
