package discovery

import (
	"context"

	v1beta1 "ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Observed struct {
	discovery *v1beta1.Discovery
}

type Observer struct {
	Client  client.Client
	Request ctrl.Request
	Context context.Context
}

func (o *Observer) observe(observed *Observed) error {
	observedDiscovery := new(v1beta1.Discovery)
	err := o.observeDiscovery(o.Request.NamespacedName, observedDiscovery)
	if err != nil {
		return err
	}

	// default everything here...

	observed.discovery = observedDiscovery
	return nil
}

func (o *Observer) observeDiscovery(key types.NamespacedName, discovery *v1beta1.Discovery) error {
	return o.Client.Get(o.Context, key, discovery)
}
