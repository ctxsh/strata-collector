package registry

import (
	"context"
	"fmt"

	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getCollecctor returns the collector pool for a discovery service.
func (r *Registry) getCollector(ctx context.Context, spec v1beta1.DiscoverySpec) (v1beta1.Collector, error) {
	var list v1beta1.CollectorList
	err := r.client.List(ctx, &list, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(spec.Collector.MatchLabels),
	})
	if err != nil {
		return v1beta1.Collector{}, err
	}

	if len(list.Items) == 0 {
		return v1beta1.Collector{}, fmt.Errorf("no collector found for discovery")
	} else if len(list.Items) > 1 {
		return v1beta1.Collector{}, fmt.Errorf("multiple collectors found for discovery, use additional labels to narrow the search to one")
	}

	return list.Items[0], nil
}

// namespacedName creates a NamespacedName type from an object interface.
func namespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}
