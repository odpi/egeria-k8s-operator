/*
Copyright 2022 Contributors to the Egeria project.

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

// EgeriaPlatformSpec : Desired State for Egeria Platform
type EgeriaPlatformSpec struct {
	// TODO: Look at moving to use the standard scaling approach https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#scale-subresource
	// Number of replicas for this platform (ie number of pods to run)
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// +kubebuilder:default:=1
	Size int32 `json:"replicas,omitempty"`
	// Secret containing TLS keys and certs
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	Security string `json:"security-secret,omitempty"`
	// Container image to use, overriding operator configuration
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:default="quay.io/odpi/egeria:latest"
	Image string `json:"image,omitempty"`
	// Configmap used for server configuration
	// +kubebuilder:validation:MaxItems=253
	// +kubebuilder:validation:MinItems=1
	// Should be unique, but cannot be set - restriction of schema validation
	ServerConfig []string `json:"serverconfig"`
}

// EgeriaPlatformStatus : Observed state of Egeria Platform
type EgeriaPlatformStatus struct {
	// Observed Egeria version from platform origin
	//TODO: Version - may be better via healthchecks
	Version string `json:"version,omitempty"`
	// list of server names that are active - this should match those configured
	Activeservers []string `json:"active,omitempty"`
	// Get info about the config being used
	ManagedService    string   `json:"service,omitempty"`
	ManagedDeployment string   `json:"deployment,omitempty"`
	Pods              []string `json:"pods,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EgeriaPlatform is the Schema for the egeriaplatforms API
// +kubebuilder:printcolumn:name="Pods",JSONPath=".status.pods",description="Current pods running",type=string
// +kubebuilder:printcolumn:name="Deployment",JSONPath=".status.deployment",description="Deployment",type=string
// +kubebuilder:printcolumn:name="Age",JSONPath=".metadata.creationTimestamp",description="Time since create",type=date
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
