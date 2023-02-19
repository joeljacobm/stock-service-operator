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
	"fmt"

	"edb.com/stock-service/backend"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var stockreportlog = logf.Log.WithName("stockreport-resource")

func (r *StockReport) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-stock-service-edb-com-v1-stockreport,mutating=true,failurePolicy=fail,sideEffects=None,groups=stock-service.edb.com,resources=stockreports,verbs=create;update,versions=v1,name=mstockreport.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &StockReport{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *StockReport) Default() {
	stockreportlog.Info("default", "name", r.Name)

}

//+kubebuilder:webhook:path=/validate-stock-service-edb-com-v1-stockreport,mutating=false,failurePolicy=fail,sideEffects=None,groups=stock-service.edb.com,resources=stockreports,verbs=create;update,versions=v1,name=vstockreport.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &StockReport{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *StockReport) ValidateCreate() error {
	stockreportlog.Info("validate create", "name", r.Name)

	if !backend.GetBackend(r.Spec.Api, r.Spec.Symbol, stockreportlog).IsValidSymbol() {
		return fmt.Errorf("error creating %s due to invalid symbol %s", r.Name, r.Spec.Symbol)
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *StockReport) ValidateUpdate(old runtime.Object) error {
	stockreportlog.Info("validate update", "name", r.Name)

	if !backend.GetBackend(r.Spec.Api, r.Spec.Symbol, stockreportlog).IsValidSymbol() {
		return fmt.Errorf("error updating %s due to invalid symbol %s", r.Name, r.Spec.Symbol)
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *StockReport) ValidateDelete() error {
	return nil
}
