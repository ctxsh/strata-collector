package collector

import (
	"context"

	"ctx.sh/strata-collector/pkg/service"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Handler struct {
	client   client.Client
	log      logr.Logger
	recorder record.EventRecorder
	observed Observed
	services *service.Manager
}

func (h *Handler) reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	h.log.V(8).Info("request received", "request", request)

	observer := &Observer{
		Client:  h.client,
		Request: request,
		Context: ctx,
	}

	if err := observer.observe(&h.observed); err != nil {
		h.log.Error(err, "unable to observe current state")
		return ctrl.Result{}, err
	}

	reconciler := &Reconciler{
		client:   h.client,
		log:      h.log,
		recorder: h.recorder,
		observed: h.observed,
		services: h.services,
	}

	result, err := reconciler.reconcile(ctx, request)
	if err != nil {
		h.log.Error(err, "unable to reconcile request")
	}

	return result, err
}
