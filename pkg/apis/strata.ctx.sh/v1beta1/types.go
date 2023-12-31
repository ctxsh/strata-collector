// Copyright 2023 Rob Lyon <rob@ctxswitch.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DiscoveryResources represets the resources that will be included in
// discovery.
type DiscoveryResources struct {
	// +optional
	Pods *bool `json:"pods,omitempty"`
	// +optional
	Services *bool `json:"services,omitempty"`
	// +optional
	Endpoints *bool `json:"endpoints,omitempty"`
}

// DiscoverySpec represents the parameters for the discovery service.
type DiscoverySpec struct {
	// +required
	// Collector is the label selector used to identify the collector
	// pool that will be used for processing.
	Collectors []corev1.ObjectReference `json:"collector"`
	// +optional
	// Selector is the label selector used to filter the resources
	// used by the discovery service.  If not set, then all resources will
	// evaluated.
	Selector metav1.LabelSelector `json:"selector"`
	// +optional
	// Enabled is a flag to enable or disable the discovery worker.
	Enabled *bool `json:"enabled"`
	// +optional
	// IntervalSeconds is the interval in seconds that the discovery worker
	// rediscover resources and send them to the processing channel.
	IntervalSeconds *int64 `json:"intervalSeconds"`
	// +optional
	// Prefix is the annotation prefix used to gather scrape information
	// from discovered resources.  By default it is set to "prometheus.io".
	Prefix *string `json:"prefix"`
	// +optional
	// Resources represents whether or not a resource will be included during
	// discovery.  By default all resources will be included.
	Resources *DiscoveryResources `json:"resources"`
}

// DiscoveryStatus represents the status of a discovery service.
type DiscoveryStatus struct {
	// DiscoveredResourcesCount is the number of resources that have been discovered
	// by the discovery service in a single run.
	DiscoveredResourcesCount int64 `json:"discoveredResourcesCount"`
	// LastDiscovered is the last time that the discovery service
	// ran and discovered resources.
	LastDiscovered metav1.Time `json:"lastDiscovered"`
	// ReadyCollectors is the number of upstream collectors that are connected and ready
	// to recieved the discovered resources.
	ReadyCollectors int64 `json:"readyCollectors"`
	// TotalCollectors is the total number of configured collectors.
	TotalCollectors int64 `json:"totalCollectors"`
	// InFlightResources is the number of resources waiting on the collectors for processing
	InFlightResources int64 `json:"inFlightResources"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=dx,singular=discovery
// +kubebuilder:printcolumn:name="Enabled",type="boolean",JSONPath=".spec.enabled"
// +kubebuilder:printcolumn:name="Ready Collectors",type="integer",JSONPath=".status.readyCollectors"
// +kubebuilder:printcolumn:name="Total Collectors",type="integer",JSONPath=".status.totalCollectors",priority=1
// +kubebuilder:printcolumn:name="Discovered",type="integer",JSONPath=".status.discoveredResourcesCount",priority=1
// +kubebuilder:printcolumn:name="In Flight",type="integer",JSONPath=".status.inFlightResources",priority=1
// +kubebuilder:printcolumn:name="Last",type="date",JSONPath=".status.lastDiscovered"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

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

// TLS represents the configurations needed to establish a TLS connection
// to a scrape endpoint.  This will probably change a bit when I start working
// on the collector and setting up the http client.  We should allow the service
// to pull in certs from k8s as well, but this will allow mounting those secrets
// into the pod.  TBH, I don't know how often this will be used since most scrape
// endpoints that I've seen have not been encrypted.
type TLS struct {
	// +optional
	// Path to the CA certificate
	CA *string `json:"ca,omitempty"`
	// +optional
	// Path to the certificate file
	Cert *string `json:"cert,omitempty"`
	// +optional
	// Path to the private key
	Key *string `json:"key,omitempty"`
	// +optional
	// InsecureSkipVerify enables/disables certificate verification between the collector and
	// the scrape endpoint.
	InsecureSkipVerify *bool `json:"inseccureSkipVerify,omitempty"`
}

// Stdout represents the configuration for the stdout data sink.
type Stdout struct{}

// Nats represents the configuration for the nats data sink.
type Nats struct {
	// +optional
	// Port is the port that the nats server is listening on.
	Port *int32 `json:"port,omitempty"`
	// +optional
	// Subject is the subject that the collector will publish to.
	Subject *string `json:"subject,omitempty"`
	// +optional
	// URL is the url of the nats server.
	URL *string `json:"url,omitempty"`
}

// CollectorOutput represents the configuration for the data sink that will
// receive the collected metrics.  Currently only supports a single output,
// but in the future we will consider supporting multiple outputs.
type CollectorOutput struct {
	// +optional
	// Nats is the configuration for the nats data sink.
	Nats *Nats `json:"nats,omitempty"`
	// +optional
	// Stdout is the configuration for the stdout data sink.
	Stdout *Stdout `json:"stdout,omitempty"`
}

// CollectorClipFilter represents the configuration for the clip filter.
type CollectorClipFilter struct {
	// +optional
	// Max is the maximum value that will be allowed.
	Max *float64 `json:"max,omitempty"`
	// +optional
	// Min is the minimum value that will be allowed.
	Min *float64 `json:"min,omitempty"`
	// +optional
	// Inclusive specifies whether or not the max and min values are inclusive
	// when evaluating.
	Inclusive *bool `json:"inclusive,omitempty"`
}

// CollectorExcludeFilter represents the configuration for the exclude filter.
type CollectorExcludeFilter struct {
	// +optional
	// Values is a list of values that will be excluded.
	Values []float64 `json:"values,omitempty"`
}

// CollectorFilters represents the filters that will be used to filter the
// metrics prior to sending them to the data output.
type CollectorFilters struct {
	// +optional
	// +nullable
	// Clip is a filter function that removes metric values that are outside
	// the max and min values.
	Clip *CollectorClipFilter `json:"clip,omitempty"`
	// +optional
	// +nullable
	// Exclude is a filter function that removes metric values that listed.
	Exclude *CollectorExcludeFilter `json:"exclude,omitempty"`
}

// CollectorSpec represents the parameters for the collector service.
type CollectorSpec struct {
	// +optional
	// BufferSize is the size of the buffer that will be used to queue
	// resources for processing.  If not set, then the default buffer size
	// will be used.
	BufferSize *int64 `json:"bufferSize"`
	// +optional
	// Encoder is the encoding that will be used to encode the metrics
	// that are sent to the data sink.  If not set, then the default
	// encoding will be used.
	Encoder *string `json:"encoder"`
	// +optional
	// Enabled is a flag to enable or disable the collector pool.
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
	// the metadata will not be included.  If set to true, then the namespace,
	// resource kind, and resource version will be added as tags.
	IncludeMetadata *bool `json:"includeMetadata"`
	// +optional
	// Workers is the number of workers in the collection pool that will
	// be used to collect metrics.
	Workers *int64 `json:"workers"`
	// +optional
	// CollectorOutput is the configuration for the data sink that will
	// receive the collected metrics.
	Output *CollectorOutput `json:"output"`
	// +optional
	// Filters is a list of filters that will be used to filter the metrics
	// prior to sending them to the data output.
	Filters *CollectorFilters `json:"filters"`
}

// CollectorStatus represents the status of a collector pool.
type CollectorStatus struct {
	// ID is the unique identifier for the collector pool.  Initially we can use it to
	// track the processing channels, but I think it would be beneficial to use it to
	// potentially add to the metrics that are collected as a reference back to the
	// pool that it was collected from.  For the most part this won't grow too large
	// and impact cardinality, however in restart conditions the id would be reset...
	// so it would only really be useful for short term correlations.  It's going to
	// be a uuid represented as a string.
	ID string `json:"id"`
	// RegisteredDiscoveries is the number of discovery services that are
	// registered to the collector.
	RegisteredDiscoveries int64 `json:"registeredDiscoveries"`
	// InFlightResources is the number of queued resources that are ready to
	// be processed.
	InFlightResources int64 `json:"inFlightResources"`
	// TotalSent is the number of metrics that have been sent to the output
	// successfully.
	TotalSent int64 `json:"totalSent"`
	// TotalErrors is the umber of metrics that have failed to be sent to the
	// output.
	TotalErrors int64 `json:"totalErrors"`
	// TotalFiltered is the number of metrics that have been filtered out by
	// the collector.
	TotalFiltered int64 `json:"totalFiltered"`
	// MetricsCollected is the number of metrics collected by the collector.
	MetricsCollected int64 `json:"metricsCollected"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=cx,singular=collector
// +kubebuilder:printcolumn:name="Enabled",type="boolean",JSONPath=".spec.enabled"
// +kubebuilder:printcolumn:name="Collected",type="integer",JSONPath=".status.metricsCollected"
// +kubebuilder:printcolumn:name="Sent",type="integer",JSONPath=".status.totalSent",priority=1
// +kubebuilder:printcolumn:name="Errors",type="integer",JSONPath=".status.totalErrors",priority=1
// +kubebuilder:printcolumn:name="Filtered",type="integer",JSONPath=".status.totalFiltered",priority=1
// +kubebuilder:printcolumn:name="Registered",type="integer",JSONPath=".status.registeredDiscoveries",priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Collector represents a pool of collection workers that will collect metrics
// from pods, services, and endpoints provided by the discovery service.  The
// collector will then send metrics to the configured data sink.
type Collector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CollectorSpec `json:"spec"`
	// +optional
	Status CollectorStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CollectorList represents a list of managed collector pools.
type CollectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Collector `json:"items"`
}
