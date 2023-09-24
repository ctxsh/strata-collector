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

package resource

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Metadata represents the object metadata information we will
// pass on as tags.
type Metadata struct {
	Kind            string
	ResourceVersion string
	Namespace       string
}

// NewMetadata creates a new metadata object using information found from a client.Object
// interface.
func NewMetadata(obj client.Object) Metadata {
	return Metadata{
		Kind:            obj.GetObjectKind().GroupVersionKind().GroupKind().String(),
		Namespace:       obj.GetNamespace(),
		ResourceVersion: obj.GetResourceVersion(),
	}
}

// NewMetadataFromRef creates a new metadata object using information found from
// a v1 ObjectReference.
func NewMetadataFromRef(obj corev1.ObjectReference) Metadata {
	return Metadata{
		Kind:            obj.GetObjectKind().GroupVersionKind().GroupKind().String(),
		Namespace:       obj.Namespace,
		ResourceVersion: obj.ResourceVersion,
	}
}
