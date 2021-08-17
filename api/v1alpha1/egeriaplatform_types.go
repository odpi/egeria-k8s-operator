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

type EgeriaServerSpec struct {
	// Name of this server - must match name in provided config
	Name string `json:"name"`
	// k8s secret containing only the full server config document
	ConfigSecret string `json:"config-secret"`
}

type EgeriaStorageSpec struct {
	// Optional Persistent Storage size ie 10G for PVC
	StorageSize string `json:"storagesize,omitempty"`
	// Optional Storage Class for PVC - empty will use default
	StorageClass string `json:"storageclass,omitempty"`
}

// Desired State for Egeria Platform
type EgeriaPlatformSpec struct {
	// Number of replicas for this platform (ie number of pods to run)
	Size int32 `json:"replicas,omitempty"`
	// Secret containing TLS keys and certs
	Security string `json:"security-secret,omitempty"`
	// Container image to use, overriding operator configuration
	Image string `json:"image,omitempty"`
	servers []EgeriaServerSpec `json:"servers"`
	storage EgeriaStorageSpec `json:"storage,omitempty"`
}

// Observed state of Egeria Platform
type EgeriaPlatformStatus struct {
	// Observed Egeria version from platform origin
	Version string `json:"version"`
	// list of server names that are active - this should match those configured
	Activeservers []string `json:"active"`

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
