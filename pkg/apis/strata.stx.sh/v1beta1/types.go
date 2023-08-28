package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CollectorSpec struct {
	Enabled bool `json:"enabled"`
}

type CollectorStatus struct {
	Enabled bool `json:"enabled"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresources:status
// +kubebuilder:resource:scope=Namespaced,shortName=coll,singular=collector
// +kubebuilder:printcolumn:name="Enabled",type="boolean",JSONPath=".status.enabled"
type Collector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CollectorSpec `json:"spec"`
	// +optional
	Status CollectorStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CollectorList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []Collector `json:"items"`
}
