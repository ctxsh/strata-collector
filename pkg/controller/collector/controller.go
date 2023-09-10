package collector

import (
	"context"

	v1beta1 "ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Controller Interface implementation
type Controller struct {
	Client client.Client
	Log    logr.Logger
	Mgr    ctrl.Manager
}

// SetupWithManager creates a new controller for the supplied manager which
// watches Collectors.
func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	r.Mgr = mgr

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Collector{}).
		WithEventFilter(r.predicates()).
		Complete(r)
}

// +kubebuilder:rbac:groups=strata.ctx.sh,resources=collectors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=strata.ctx.sh,resources=collectors/status,verbs=get;update;patch

// Reconcile ensures that the existing state of a resource matches requested state.
func (r *Controller) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	handler := Handler{
		client:   r.Mgr.GetClient(),
		log:      r.Log.WithValues("name", request.Name, "namespace", request.Namespace),
		recorder: r.Mgr.GetEventRecorderFor("StrataCollector"),
	}
	return handler.reconcile(ctx, request)
}

// predicates returns a map of predicate functions which determine the conditions
// of whether or not to reconcile an object.
func (r *Controller) predicates() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld == nil || e.ObjectNew == nil {
				return false
			}
			// Only update the object if the resource version has been modified.  This
			// ensures that we are not trying to reconcile the object if only the status
			// has changed.
			return e.ObjectNew.GetResourceVersion() != e.ObjectOld.GetResourceVersion()
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
	}
}
