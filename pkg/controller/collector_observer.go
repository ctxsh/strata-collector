package controller

import (
	"context"

	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ObserveredCollector struct {
	collector *v1beta1.Collector
	// observeTime time.Time
}

type CollectorObserver struct {
	Client  client.Client
	Request ctrl.Request
	Context context.Context
}

func (o *CollectorObserver) observe(observed *ObserveredCollector) error {
	observedCollector := new(v1beta1.Collector)
	err := o.observeCollector(o.Request.NamespacedName, observedCollector)
	if err != nil {
		return err
	}

	// default everything here...

	observed.collector = observedCollector
	return nil
}

func (o *CollectorObserver) observeCollector(key types.NamespacedName, collector *v1beta1.Collector) error {
	return o.Client.Get(o.Context, key, collector)
}
