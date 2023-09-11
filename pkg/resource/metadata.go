package resource

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Metadata represents the object metadata information we will
// pass on as tags.
type Metadata struct {
	Kind            string
	ResourceVersion string
	Namespace       string
}

// NewMetadata creates a new metadata object using information found from a client.Object
// interface.
func NewMetadata(obj client.Object) Metadata {
	return Metadata{
		Kind:            obj.GetObjectKind().GroupVersionKind().GroupKind().String(),
		Namespace:       obj.GetNamespace(),
		ResourceVersion: obj.GetResourceVersion(),
	}
}

// NewMetadataFromRef creates a new metadata object using information found from
// a v1 ObjectReference.
func NewMetadataFromRef(obj corev1.ObjectReference) Metadata {
	return Metadata{
		Kind:            obj.GetObjectKind().GroupVersionKind().GroupKind().String(),
		Namespace:       obj.Namespace,
		ResourceVersion: obj.ResourceVersion,
	}
}
