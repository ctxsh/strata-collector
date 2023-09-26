package v1beta1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// ValidateCreate implements webhook Validator.
func (d *Discovery) ValidateCreate() (admission.Warnings, error) {
	return d.validate()
}

// ValidateUpdate implements webhook Validator.
func (d *Discovery) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return d.validate()
}

// ValidateDelete implements webhook Validator.
func (d *Discovery) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (d *Discovery) validate() (admission.Warnings, error) {
	warn := make(admission.Warnings, 0)
	if *d.Spec.IntervalSeconds < 1 {
		warn = append(warn, "IntervalSeconds must be greater than 0")
	}

	for _, r := range d.Spec.Collectors {
		if r.Name == "" {
			warn = append(warn, "Collector name must be set")
		}
	}

	if len(warn) > 0 {
		return warn, fmt.Errorf("invalid discovery")
	}

	return nil, nil
}

// ValidateCreate implements webhook Validator.
func (c *Collector) ValidateCreate() (admission.Warnings, error) {
	return c.validate()
}

// ValidateUpdate implements webhook Validator.
func (c *Collector) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return c.validate()
}

// ValidateDelete implements webhook Validator.
func (c *Collector) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (c *Collector) validate() (admission.Warnings, error) {
	warn := make(admission.Warnings, 0)

	if *c.Spec.Workers <= 0 {
		warn = append(warn, "Workers must be greater than or equal to 0")
	}

	if len(warn) > 0 {
		return warn, fmt.Errorf("invalid collector")
	}

	return nil, nil
}
