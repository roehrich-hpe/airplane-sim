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

// PedalsSpec defines the desired state of Pedals
type PedalsSpec struct {
	// Pressed indicates which pedal is pressed
	// +kubebuilder:validation:Enum=none;left;right
	// +kubebuilder:default:=none
	Pressed string `json:"pressed,omitempty"`
}

// PedalsStatus defines the observed state of Pedals
type PedalsStatus struct {
	// LinkagePosition indicates where the pedal linkage is currently
	// +kubebuilder:validation:Enum=neutral;left;right
	// +kubebuilder:default:=neutral
	LinkagePosition string `json:"linkagePosition,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="PRESSED",type="string",JSONPath=".spec.pressed",description="Indicates which pedal is pressed"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// Pedals is the Schema for the pedals API
type Pedals struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PedalsSpec   `json:"spec"`
	Status PedalsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PedalsList contains a list of Pedals
type PedalsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pedals `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pedals{}, &PedalsList{})
}
