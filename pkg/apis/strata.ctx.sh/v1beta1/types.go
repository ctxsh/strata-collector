package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DiscoverySpec represents the parameters for the discovery service.
type DiscoverySpec struct {
	// +required
	// Collection is the label selector used to identify the collector
	// pool that will be used for processing.
	Collection metav1.LabelSelector `json:"collection"`
	// +optional
	// Selector is the label selector used to filter the resources
	// used by the discovery service.  If not set, then all resources will
	// evaluated.
	Selector metav1.LabelSelector `json:"selector"`
	// +optional
	// Enabled is a flag to enable or disable the discovery worker.
	Enabled *bool `json:"enabled"`
	// +optional
	// IncludeAnnotations is a list of annotations that will be added as tags
	// to the metrics that are collected.  By default no annotations will be
	// added.  If set, then the annotations will be added as tags.  Currently
	// only full string matches are supported.  In the future, wildcard matches
	// will be supported.
	IncludeAnnotations []string `json:"includeAnnotations"`
	// +optional
	// IncludeLabels is a list of labels that will be added as tags to the
	// metrics that are collected.  By default no labels will be added.  If
	// set, then the labels will be added as tags.  Currently only full string
	// matches are supported.  In the future, wildcard matches will be supported.
	IncludeLabels []string `json:"includeLabels"`
	// +optional
	// IncludeMetadata determines whether or not the metadata for the resource
	// will be added as tags to the metrics that are collected.  By default
	// the metadata will not be included.  If set to true, then the name, namespace,
	// and resourceVersion will be added as tags.
	IncludeMetadata *bool `json:"includeMetadata"`
	// +optional
	// IntervalSeconds is the interval in seconds that the discovery worker
	// rediscover resources and send them to the processing channel.
	IntervalSeconds *int64 `json:"intervalSeconds"`
	// +optional
	// Prefix is the annotation prefix used to gather scrape information
	// from discovered resources.  By default it is set to "prometheus.io".
	Prefix *string `json:"prefix"`
}

// DiscoveryStatus represents the status of a discovery service.
type DiscoveryStatus struct {
	Enabled        bool   `json:"enabled"`
	LastDiscovered string `json:"lastDiscovered"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresources:status
// +kubebuilder:resource:scope=Namespaced,shortName=discover,singular=discovery
// +kubebuilder:printcolumn:name="Enabled",type="boolean",JSONPath=".status.enabled"

// Discovery represents a discovery service that will collect pods, services, and
// endpoints from a k8s cluster.
type Discovery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DiscoverySpec `json:"spec"`
	// +optional
	Status DiscoveryStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DiscoveryList represents a list of managed discovery services.
type DiscoveryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Discovery `json:"items"`
}

// CollectorSpec represents the parameters for the collector service.
type CollectorSpec struct {
	// +optional
	// Enabled is a flag to enable or disable the collector pool.
	Enabled *bool `json:"enabled"`
	// +optional
	// Workers is the number of workers in the collection pool that will
	// be used to collect metrics.
	Workers *int64 `json:"workers"`
}

// CollectorStatus represents the status of a collector pool.
type CollectorStatus struct {
	// Enabled represents whether the collector pool is enabled or not.
	Enabled bool `json:"enabled"`
	// Discoveries represents the list of discoveries that are sending
	// discovered resources to the collector pool.
	Discoveries []string `json:"discoveries"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresources:status
// +kubebuilder:resource:scope=Namespaced,shortName=coll,singular=collector
// +kubebuilder:printcolumn:name="Enabled",type="boolean",JSONPath=".status.enabled"

// Collector represents a pool of collection workers that will collect metrics
// from pods, services, and endpoints provided by the discovery service.  The
// collector will then send metrics to the configured data sink.
type Collector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CollectorSpec   `json:"spec"`
	Status            CollectorStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CollectorList represents a list of managed collector pools.
type CollectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Collector `json:"items"`
}
