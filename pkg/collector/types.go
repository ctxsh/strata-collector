package collector

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resource represents a discovered kubernetes resource and contains
// additional information needed to transform metrics from the metrics
// endpoint.
type Resource struct {
	IP                 string
	Scrape             bool
	Scheme             string
	Port               string
	Path               string
	IncludeMetadata    bool
	IncludeLabels      []string
	IncludeAnnotations []string
	ObjectMeta         metav1.ObjectMeta
}
