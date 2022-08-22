/*
Copyright 2022.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RudderSpec defines the desired state of Rudder
type RudderSpec struct {
	// Position indicates where we want the rudder to be placed
	// +kubebuilder:validation:Enum=neutral;left;right
	// +kubebuilder:default:=neutral
	Position string `json:"position,omitempty"`
}

// RudderStatus defines the observed state of Rudder
type RudderStatus struct {
	// Position indicates where the rudder is currently
	// +kubebuilder:validation:Enum=neutral;left;right
	// +kubebuilder:default:=neutral
	Position string `json:"position,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="DESIRED POSITION",type="string",JSONPath=".spec.position",description="Desired position of rudder"
//+kubebuilder:printcolumn:name="CURRENT POSITION",type="string",JSONPath=".status.position",description="Current position of rudder"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// Rudder is the Schema for the rudders API
type Rudder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RudderSpec   `json:"spec"`
	Status RudderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RudderList contains a list of Rudder
type RudderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rudder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rudder{}, &RudderList{})
}
