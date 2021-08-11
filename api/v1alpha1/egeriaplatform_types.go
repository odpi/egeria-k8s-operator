/*
Copyright 2021 Contributors to the Egeria project.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Desired State for Egeria Platform
type EgeriaPlatformSpec struct {
	// Number of replicas for this platform (ie number of pods to run)
	Size int32 `json:"replicas"`
	// Secret containing TLS keys and certs
	Security string `json:"security"`
	// Optional container image to use, overriding operator configuration
	Image string `json:"image,omitempty"`
	// Optional Storage size ie 10G for PVC
	StorageSize string `json:"storagesize,omitempty"`
	// Optional Storage Class for PVC
	StorageClass string `json:"storageclass,omitempty"`
	// TODO: Consider including Audit Log connector & other platform config
}

// Observed state of Egeria Platform
type EgeriaPlatformStatus struct {
	// Observed version from platform origin
	Version string `json:"version"`
	// list of servers we know about
	Knownservers []string `json:"known"`
	// list of servers that are active
	Activeservers []string `json:"active"`
	// Calculated load indicator to aid in scheduling
	Workload int32 `json:"workload"`
	// TODO: What else do we 'observe' about the status of the platform
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EgeriaPlatform is the Schema for the egeriaplatforms API
type EgeriaPlatform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EgeriaPlatformSpec   `json:"spec,omitempty"`
	Status EgeriaPlatformStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EgeriaPlatformList contains a list of EgeriaPlatform
type EgeriaPlatformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EgeriaPlatform `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EgeriaPlatform{}, &EgeriaPlatformList{})
}
