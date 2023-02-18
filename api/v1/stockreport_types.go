/*
Copyright 2023.

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

// StockReportSpec defines the desired state of StockReport
type StockReportSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of StockReport. Edit stockreport_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// StockReportStatus defines the observed state of StockReport
type StockReportStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StockReport is the Schema for the stockreports API
type StockReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StockReportSpec   `json:"spec,omitempty"`
	Status StockReportStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StockReportList contains a list of StockReport
type StockReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StockReport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StockReport{}, &StockReportList{})
}
