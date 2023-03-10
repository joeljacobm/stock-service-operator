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

package controllers

import (
	"context"

	"fmt"
	"time"

	stockservicev1 "edb.com/stock-service/api/v1"
	v1 "edb.com/stock-service/api/v1"
	"edb.com/stock-service/backend"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// StockReportReconciler reconciles a StockReport object
type StockReportReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stock-service.edb.com,resources=stockreports,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stock-service.edb.com,resources=stockreports/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stock-service.edb.com,resources=stockreports/finalizers,verbs=update
//+kubebuilder:rbac:groups=*,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=configmaps/status,verbs=get

func (r *StockReportReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling")

	stockReport := &v1.StockReport{}

	if err := r.Get(ctx, req.NamespacedName, stockReport); err != nil {
		log.Error(err, "unable to fetch StockReport")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return r.UpdateErrorStatus(stockReport, client.IgnoreNotFound(err))
	}

	if stockReport.Spec.Symbol == "" {
		return ctrl.Result{}, fmt.Errorf("stock symbol cannot be empty")
	}
	backend := backend.GetBackend(stockReport.Spec.Api, stockReport.Spec.Symbol, log)
	price, err := backend.GetStockPrice()
	if err != nil {
		return r.UpdateErrorStatus(stockReport, fmt.Errorf("failed fetching stock price for symbol %s with error: %s", stockReport.Spec.Symbol, err.Error()))
	}
	if price == 0 {
		return r.UpdateErrorStatus(stockReport, fmt.Errorf("invalid price for symbol %s", stockReport.Spec.Symbol))
	}
	log.Info("successfully retrieved stock price", "symbol", stockReport.Spec.Symbol, "price", price)

	duration, err := time.ParseDuration(stockReport.Spec.RefreshInterval)
	if err != nil {
		return r.UpdateErrorStatus(stockReport, err)
	}

	cmName := stockReport.Name + "-cm"

	// Create a new ConfigMap with the stock price
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: stockReport.Namespace,
		},
		Data: map[string]string{
			"updateTime":   time.Now().String(),
			"price":        fmt.Sprintf("%.3f", price),
			"stock_symbol": stockReport.Spec.Symbol,
			"api":          stockReport.Spec.Api,
		},
	}
	configMap.SetOwnerReferences([]metav1.OwnerReference{{APIVersion: stockReport.APIVersion, Name: stockReport.Name, Kind: stockReport.Kind, UID: stockReport.GetUID()}})

	// Check if the ConfigMap already exists
	foundConfigMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, foundConfigMap)

	switch {
	case err != nil && errors.IsNotFound(err):
		// Create the ConfigMap if it doesn't exist
		if err = r.Create(ctx, configMap); err != nil {
			return r.UpdateErrorStatus(stockReport, err)

		}
	case err == nil:
		// Update the ConfigMap if it already exists
		foundConfigMap.Data = configMap.Data
		err = r.Update(ctx, foundConfigMap)
		if err != nil {
			return r.UpdateErrorStatus(stockReport, err)
		}
	default:
		// Return an error if we couldn't fetch the ConfigMap
		return r.UpdateErrorStatus(stockReport, err)
	}

	stockReport.Status = v1.StockReportStatus{
		LastRefreshed: metav1.Now(),
		Status:        "Ready",
		ConfigMap:     cmName,
	}
	if err := r.Status().Update(context.Background(), stockReport); err != nil {
		return r.UpdateErrorStatus(stockReport, err)
	}

	return ctrl.Result{RequeueAfter: duration}, nil
}

func (r *StockReportReconciler) UpdateErrorStatus(obj *v1.StockReport, err error) (ctrl.Result, error) {
	obj.Status.Status = "Error"
	if updateErr := r.Status().Update(context.Background(), obj); updateErr != nil {
		return ctrl.Result{}, updateErr
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *StockReportReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stockservicev1.StockReport{}).
		Complete(r)
}
