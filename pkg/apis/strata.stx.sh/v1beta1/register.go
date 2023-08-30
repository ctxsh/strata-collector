package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Version specifies the API version
const Version = "v1beta1"

// SchemaGroupVersion is the group version that wil be used to register the objects
var SchemaGroupVersion = schema.GroupVersion{
	Group:   strata.GroupName,
	Version: Version,
}

// Kind takes an unqualified kind and returns back a Group and qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemaGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns back a Group and qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemaGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
)

// addKnownTypes adds a list of known types to the scheme
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Collector{},
		&CollectorList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
