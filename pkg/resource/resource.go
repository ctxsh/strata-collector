package resource

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DefaultScrapeAnnotation bool   = false
	DefaultSchemeAnnotation string = "http"
	DefaultPathAnnotation   string = "/metrics"
	DefaultPortAnnotation   string = "9090"
	DefaultIncludeMeta      bool   = false
	DefaultPrefix           string = "prometheus.io"
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
	Metadata           Metadata
	IncludeLabels      []string
	Labels             Labels
	IncludeAnnotations []string
	Annotations        Annotations
	Timestamp          time.Time
}

// New returns a new defaulted resource.  The scrape annotations are used initially.
func New(a map[string]string, prefix string) *Resource {
	return defaulted(a, prefix)
}

// WithAnnotations sets the annotations of the resource
func (r *Resource) WithAnnotations(a map[string]string) *Resource {
	r.Annotations = a
	return r
}

// WithIP sets the IP address of the resource.
func (r *Resource) WithIP(ip string) *Resource {
	r.IP = ip
	return r
}

// WithLabels sets the labels of the resource.
func (r *Resource) WithLabels(labels map[string]string) *Resource {
	r.Labels = labels
	return r
}

// WithMetadata creates a new metadata object containing metadata information
// that will be used for tags in the collection process.
func (r *Resource) WithMetadata(obj client.Object) *Resource {
	r.Metadata = NewMetadata(obj)
	return r
}

// WithMetadataRef creates an new metadata object using the target references.  It's
// primarily (only) used for the endpoint resource creation off of headless services
// since they point back to existing pods and not the parent service.
func (r *Resource) WithMetadataRef(obj *corev1.ObjectReference) *Resource {
	r.Metadata = NewMetadataFromRef(*obj)
	return r
}

// defaulted returns a new resources defaulted with the scrap annotations. By default
// we support the common prometheus annotations using the prefix "prometheus.io" as
// to be a drop in replacement for the prometheus operator.  The prefix can be changed
// to support other annotations as well.  In the future, we will support other filters
// to include or exclude labels based on annotations.
//
// The current annotations are supported:
//
// <prefix>/scrape
// The scrape annotations controls the enablement of the scrape service.  It represents
// a boolean value.  Though it will validate any string, anything other than the 'true'
// string will result in a false value.
//
// <prefix>/scheme
// The scheme annotations control the http scheme used to communicate with the prometheus
// scrape endpoint.  Valid values are 'http' or 'https'.  If https has been selected you
// will need to add TLS information in the Collector manifest.
//
// <prefix>/port
// The port annotation is used to override the default scrape port of 9090.
//
// <prefix>/path
// The path annotation is used to overrid the default scrape path which is set to
// '/metrics' by default.
//
// TODO:
// Warning, remember that tags can explode cardinality in certain systems which can degrade
// performance signifcantly and increase cost - be it from self managed or vendor solutions.
// These are added as a way to consolidate observability requirements outside of the metrics
// clients.  It is usually considered bad practice to throw everything you think you would
// ever possibly need at any observability system.
//
// <prefix>/includeLabels
// A comma seperated string representing the names of labels to include as tags in the metrics
// that are generated.
//
// <prefix>/includeAnnotations
// A comma seperated string representing the names of any annotations to include as tags in the
// metrics.
//
// <prefix>/includeMetadata
// A comma seperated list of valid metadata fields to include as tags in the metrics.
func defaulted(annotations map[string]string, prefix string) *Resource {
	res := &Resource{
		Scrape:             DefaultScrapeAnnotation,
		Scheme:             DefaultSchemeAnnotation,
		Port:               DefaultPortAnnotation,
		Path:               DefaultPathAnnotation,
		IncludeMetadata:    false,
		IncludeLabels:      make([]string, 0),
		IncludeAnnotations: make([]string, 0),
		Timestamp:          time.Now(),
	}

	scrapeAnnotation := fmt.Sprintf("%s/scrape", prefix)
	if a, ok := annotations[scrapeAnnotation]; ok {
		res.Scrape = a == "true"
	}

	schemeAnnotation := fmt.Sprintf("%s/scheme", prefix)
	if a, ok := annotations[schemeAnnotation]; ok {
		res.Scheme = a
	}

	portAnnotation := fmt.Sprintf("%s/port", prefix)
	if a, ok := annotations[portAnnotation]; ok {
		res.Port = a
	}

	pathAnnotation := fmt.Sprintf("%s/path", prefix)
	if a, ok := annotations[pathAnnotation]; ok {
		res.Path = a
	}

	// TODO: add annotation for metadata inclusion.  Right now we only allow this
	// in the manifest, but that's an all or nothing approach and it would be better
	// just to add an annotation as an option to selectively enable metadata scraping.
	// This could also be a list of items to selectively control which metadata fields
	// are passed along as tags.

	return res
}
