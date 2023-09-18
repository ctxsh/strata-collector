package registry

import (
	"context"

	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getCollecctor returns the collector pool for a discovery service.
func (r *Registry) getCollector(ctx context.Context, refs []corev1.ObjectReference) []v1beta1.Collector {
	collectors := make([]v1beta1.Collector, len(refs))

	for i, ref := range refs {
		var collector v1beta1.Collector
		err := r.client.Get(ctx, types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}, &collector)
		if err != nil {
			r.logger.Error(err, "unable to get collector", "collector", ref)
			continue
		}

		collectors[i] = collector
	}

	return collectors
}

// namespacedName creates a NamespacedName type from an object interface.
func namespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}
