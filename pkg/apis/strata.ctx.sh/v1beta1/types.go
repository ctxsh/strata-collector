package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CollectorSpec defines the collector options.
type CollectorSpec struct {
	// +optional
	Enabled *bool `json:"enabled"`
}

// CollectorStatus represents the status of an individual collector.
type CollectorStatus struct {
	Enabled bool `json:"enabled"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresources:status
// +kubebuilder:resource:scope=Namespaced,shortName=co,singular=collector
// +kubebuilder:printcolumn:name="Enabled",type="boolean",JSONPath=".status.enabled"

// Collector represents a metrics collection worker for the strata metrics service.
type Collector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CollectorSpec `json:"spec"`
	// +optional
	Status CollectorStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CollectorList represents a list of collection workers.
type CollectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Collector `json:"items"`
}