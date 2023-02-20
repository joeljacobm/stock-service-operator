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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StockReportSpec defines the desired state of StockReport
type StockReportSpec struct {
	// Symbol is the stock symbol eg :- AAPL, GOOGL
	Symbol string `json:"symbol"`
	// +optional
	// +kubebuilder:default="60s"
	// RefreshInterval indicates the interval to fetch the latest stock price.
	// The format should match go duration type eg :- 1ms,1s,1m,1h,1d
	// Default is set to 60s
	RefreshInterval string `json:"refreshInterval,omitempty"`
	// +optional
	// +kubebuilder:validation:Enum:=yahoo;vantage
	// +kubebuilder:default=yahoo
	// Api is the finance api used to fetch the stock prices.
	// Default is to yahoo finance
	Api string `json:"api,omitempty"`
}

// StockReportStatus defines the observed state of StockReport
type StockReportStatus struct {
	LastRefreshed v1.Time `json:"lastRefreshed,omitempty"`
	Status        string  `json:"status,omitempty"`
	ConfigMap     string  `json:"configMap,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Config Map",type="string",JSONPath=".status.configMap"
// +kubebuilder:printcolumn:name="Last Refreshed",type="date",JSONPath=".status.lastRefreshed"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"

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
