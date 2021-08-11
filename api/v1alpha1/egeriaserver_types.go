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

// Desired state for Egeria Server
type EgeriaServerSpec struct {
	// server name (as known by Egeria)
	Name string `json:"name"`
	// number of replicas
	Size int32 `json:"replicas"`
	// Reference to location of server config document (exact details not yet known)
	Config string `json:"config"`
}

// Observed state of Egeria Server
type EgeriaServerStatus struct {
	// current status - unknown, inactive, active
	Status string `json:"status"`
	// Name of the platform (flavour) we are running on
	Platformname string `json:"platformname"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EgeriaServer is the Schema for the egeriaservers API
type EgeriaServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EgeriaServerSpec   `json:"spec,omitempty"`
	Status EgeriaServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EgeriaServerList contains a list of EgeriaServer
type EgeriaServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EgeriaServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EgeriaServer{}, &EgeriaServerList{})
}
