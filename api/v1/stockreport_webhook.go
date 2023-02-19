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
	"time"

	"edb.com/stock-service/backend"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	return r.ValidateSpec()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *StockReport) ValidateUpdate(old runtime.Object) error {
	return r.ValidateSpec()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *StockReport) ValidateDelete() error {
	return nil
}

func (r *StockReport) ValidateSpec() error {
	stockreportlog.Info("validating spec", "name", r.Name)
	var allErrs field.ErrorList

	if r.Spec.RefreshInterval != "" {
		if _, err := time.ParseDuration(r.Spec.RefreshInterval); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("refreshInterval"), r.Spec.RefreshInterval, errors.Wrap(err, "invalid refreshInterval format. The format should be of type time.Duration (1ms,1s,1m,1h...)").Error()))
		}
	}

	if !backend.GetBackend(r.Spec.Api, r.Spec.Symbol, stockreportlog).IsValidSymbol() {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("symbol"), r.Spec.Symbol, fmt.Sprintf("error updating %s due to invalid symbol %s", r.Name, r.Spec.Symbol)))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: r.GroupVersionKind().Group, Kind: r.GroupVersionKind().Kind},
		r.Name, allErrs)
}
