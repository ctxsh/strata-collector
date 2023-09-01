package controller

import (
	"context"

	v1beta1 "ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Reconciler struct {
	Client client.Client
	Log    logr.Logger
	Mgr    ctrl.Manager
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Collector{}).
		WithEventFilter(r.predicates()).
		Complete(r)
}

func (r *Reconciler) predicates() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
}

func (r *Reconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}
