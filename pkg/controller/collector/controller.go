package collector

import (
	"context"

	v1beta1 "ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/service"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	ctrl "sigs.k8s.io/controller-runtime"
)

// Controller Interface implementation
type Controller struct {
	Cache    cache.Cache
	Client   client.Client
	Log      logr.Logger
	Mgr      ctrl.Manager
	Registry *service.Registry
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
		registry: r.Registry,
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
			// Only update the object if the resource generation has changed.
			return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
	}
}
