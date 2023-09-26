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

// +kubebuilder:docs-gen:collapse=Apache License

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// +kubebuilder:docs-gen:collapse=Go imports

// +kubebuilder:webhook:verbs=create;update,path=/mutate-strata-ctx-sh-v1beta1-discovery,mutating=true,failurePolicy=fail,groups=strata.ctx.sh,resources=discoveries,versions=v1beta1,name=mdiscovery.strata.ctx.sh,admissionReviewVersions=v1,sideEffects=none
// +kubebuilder:webhook:verbs=create;update,path=/validate-strata-ctx-sh-v1beta1-discovery,mutating=false,failurePolicy=fail,groups=strata.ctx.sh,resources=discoveries,versions=v1beta1,name=vdiscovery.strata.ctx.sh,admissionReviewVersions=v1,sideEffects=none
// +kubebuilder:webhook:verbs=create;update,path=/mutate-strata-ctx-sh-v1beta1-collector,mutating=true,failurePolicy=fail,groups=strata.ctx.sh,resources=collectors,versions=v1beta1,name=mcollector.strata.ctx.sh,admissionReviewVersions=v1,sideEffects=none
// +kubebuilder:webhook:verbs=create;update,path=/validate-strata-ctx-sh-v1beta1-collector,mutating=false,failurePolicy=fail,groups=strata.ctx.sh,resources=collectors,versions=v1beta1,name=vcollector.strata.ctx.sh,admissionReviewVersions=v1,sideEffects=none

// SetupWebhookWithManager adds webhook for Discovery.
func (d *Discovery) SetupWebhookWithManager(mgr ctrl.Manager) error {
	// whs := mgr.GetWebhookServer()

	return ctrl.NewWebhookManagedBy(mgr).
		For(d).
		Complete()
}

var _ webhook.Defaulter = &Discovery{}

var _ webhook.Validator = &Discovery{}

// SetupWebhookWithManager adds webhook for Discovery.
func (c *Collector) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(c).
		Complete()
}

var _ webhook.Defaulter = &Collector{}

var _ webhook.Validator = &Collector{}

// +kubebuilder:docs-gen:collapse=Validate object name
