package collector

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
	registry *service.Registry
}

var requeueResult reconcile.Result = ctrl.Result{
	Requeue:      true,
	RequeueAfter: 30 * time.Second,
}

func (r *Reconciler) reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	if r.observed.collector == nil {
		if err := r.registry.DeleteCollectionPool(request.NamespacedName); err != nil {
			r.log.Error(err, "unable to delete collection pool")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if err := r.registry.AddCollectionPool(ctx, request.NamespacedName, *r.observed.collector); err != nil {
		return requeueResult, err
	}

	r.log.V(8).Info("reconciling collector")
	return ctrl.Result{}, nil
}
