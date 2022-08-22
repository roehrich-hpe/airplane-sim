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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AirplaneSpec defines the desired state of Airplane
type AirplaneSpec struct {
	// TailNumber is "N-number" registration on our tail. We support only: Nxxxxx, where X is a digit or an uppercase letter.
	// +kubebuilder:validation:Pattern:="^N[A-Z\\d]{5}$"
	TailNumber string `json:"tailNumber"`
}

// AirplaneStatus defines the observed state of Airplane
type AirplaneStatus struct {
	// Rudder names the rudder resource
	Rudder corev1.ObjectReference `json:"rudder,omitempty"`

	// Pedals names the pedals resource
	Pedals corev1.ObjectReference `json:"pedals,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="TAILNUMBER",type="string",JSONPath=".spec.tailNumber",description="N-Number registration"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// Airplane is the Schema for the airplanes API
type Airplane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AirplaneSpec   `json:"spec"`
	Status AirplaneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AirplaneList contains a list of Airplane
type AirplaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Airplane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Airplane{}, &AirplaneList{})
}
