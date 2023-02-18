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
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stockservicev1 "edb.com/stock-service/api/v1"
	v1 "edb.com/stock-service/api/v1"
	"edb.com/stock-service/backend"
)

// StockReportReconciler reconciles a StockReport object
type StockReportReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stock-service.edb.com,resources=stockreports,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stock-service.edb.com,resources=stockreports/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stock-service.edb.com,resources=stockreports/finalizers,verbs=update

func (r *StockReportReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling")

	stockReport := &v1.StockReport{}

	if err := r.Get(ctx, req.NamespacedName, stockReport); err != nil {
		log.Error(err, "unable to fetch StockReport")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if stockReport.Spec.Symbol == "" {
		return ctrl.Result{}, errors.New("stock symbol cannot be empty")
	}
	backend := backend.GetBackend(stockReport.Spec.Api, stockReport.Spec.Symbol, log)
	price, err := backend.GetStockPrice()
	if err != nil {
		log.Error(err, "failed fetching stock price", "symbol", stockReport.Spec.Symbol)
	}
	log.Info("successflly fetched stock price","symbol",stockReport.Spec.Symbol,"price",price)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StockReportReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stockservicev1.StockReport{}).
		Complete(r)
}
