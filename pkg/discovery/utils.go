package discovery

import (
	"fmt"
)

// getAnnotations returns the scrape annotations for a given resource.  By default
// we support the common prometheus annotations using the prefix "prometheus.io" as
// to be a drop in replacement for the prometheus operator.  The prefix can be changed
// to support other annotations as well.  In the future, we will support other filters
// to include or exclude labels based on annotations.
func getAnnotations(annotations map[string]string, prefix string) map[string]string {
	var scrape string = "false"
	var scheme string = "http"
	var port string = "9090"
	var path string = "/metrics"

	scrapeAnnotation := fmt.Sprintf("%s/scrape", prefix)
	if a, ok := annotations[scrapeAnnotation]; ok {
		scrape = a
	}

	schemeAnnotation := fmt.Sprintf("%s/scheme", prefix)
	if a, ok := annotations[schemeAnnotation]; ok {
		scheme = a
	}

	portAnnotation := fmt.Sprintf("%s/port", prefix)
	if a, ok := annotations[portAnnotation]; ok {
		port = a
	}

	pathAnnotation := fmt.Sprintf("%s/path", prefix)
	if a, ok := annotations[pathAnnotation]; ok {
		path = a
	}

	return map[string]string{
		"scrape": scrape,
		"scheme": scheme,
		"port":   port,
		"path":   path,
	}
}
