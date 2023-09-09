package collectors

// Resource represents a discovered kubernetes resource and contains
// the information needed to collect and fileter metrics.
type Resource struct {
	Enabled bool
	IP      string
	Port    string
	Path    string
	// labels
	// tags
}
