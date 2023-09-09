package controller

import (
	"context"
	"time"

	"ctx.sh/strata-collector/pkg/collectors"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CollectorReconciler struct {
	client     client.Client
	log        logr.Logger
	observed   ObserveredCollector
	recorder   record.EventRecorder
	collectors *collectors.Manager
}

var requeueResult reconcile.Result = ctrl.Result{
	Requeue:      true,
	RequeueAfter: 30 * time.Second,
}

func (r *CollectorReconciler) reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	if r.observed.collector == nil {
		r.collectors.Delete(request.NamespacedName)
		return ctrl.Result{}, nil
	}

	r.log.V(8).Info("reconciling collector")

	err := r.collectors.Add(request.NamespacedName, *r.observed.collector.Spec.DeepCopy())
	if err != nil {
		return requeueResult, err
	}
	return ctrl.Result{}, nil
}
