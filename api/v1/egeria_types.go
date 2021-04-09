/*


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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EgeriaSpec defines the desired state of Egeria

type EgeriaServer struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:UniqueItems=true
	Name []string `json:"name"`
	Replicas  int `json:"replicas"`
}
type EgeriaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:UniqueItems=true
	Servers []EgeriaServer `json:"servers"`
}

// EgeriaStatus defines the observed state of Egeria

// First define the struct to store info on a platform
type EgeriaServerStatus struct {
	Name string `json:"name"`
	Replicas  int `json:"replicas"`
	Node string `json:"node"`
}

type EgeriaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Nodes are the names of the pods
	Platforms []EgeriaServerStatus `json:"platforms"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Egeria is the Schema for the egeria API
// +kubebuilder:subresource:status
type Egeria struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EgeriaSpec   `json:"spec,omitempty"`
	Status EgeriaStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EgeriaList contains a list of Egeria
type EgeriaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Egeria `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Egeria{}, &EgeriaList{})
}
