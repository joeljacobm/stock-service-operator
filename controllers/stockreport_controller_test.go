package controllers

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	stockservicev1 "edb.com/stock-service/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	timeout  = time.Second * 10
	interval = time.Millisecond * 250
)

var stockreport stockservicev1.StockReport

var _ = Describe("stock report controller", func() {
	It("Should retrieve the stock price from the backend api and store it in configmap", func() {
		By("creating a new stock report object", func() {
			stockreport = stockservicev1.StockReport{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "stock-service.edb.com/v1",
					Kind:       "StockReport",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample",
					Namespace: "default",
				},
				Spec: stockservicev1.StockReportSpec{
					Symbol:          "AAPL",
					RefreshInterval: "30s",
					Api:             "yahoo",
				},
			}
			Expect(k8sClient.Create(ctx, &stockreport)).Should(Succeed())
		})
		By("checking if the configmap exists")
		configMapLookupKey := types.NamespacedName{Name: "sample-cm", Namespace: "default"}
		// We'll need to retry getting this newly created config map, given that creation may not immediately happen.
		cm := &corev1.ConfigMap{}
		Eventually(func() bool {
			err := k8sClient.Get(ctx, configMapLookupKey, cm)
			if err != nil {
				return false
			}
			return true
		}, timeout, interval).Should(BeTrue())

		By("checking the stock details have been populated")
		Expect(cm.Data["stock_symbol"]).Should(Equal("AAPL"))
		Expect(cm.Data["api"]).Should(Equal("yahoo"))
		Expect(cm.Data["price"]).ShouldNot(Equal(""))

		By("checking the owner reference of the configmap")
		Expect(cm.OwnerReferences[0].Name).Should(Equal(stockreport.Name))

	})

	It("Should update the configmap when the stock symbol is changed", func() {

		By("getting the stock report object from k8s")
		stockReportLookupKey := types.NamespacedName{Name: "sample", Namespace: "default"}

		err := k8sClient.Get(ctx, stockReportLookupKey, &stockreport)
		Expect(err).NotTo(HaveOccurred())

		By("updating the stock report object")
		stockreport.Spec.Symbol = "GOOGL"
		Expect(k8sClient.Update(ctx, &stockreport)).Should(Succeed())

		By("checking if the stock symbol has been updated in the config map")
		configMapLookupKey := types.NamespacedName{Name: "sample-cm", Namespace: "default"}
		// We'll need to retry getting this newly created config map, given that creation may not immediately happen.
		cm := &corev1.ConfigMap{}
		Eventually(func() bool {
			err := k8sClient.Get(ctx, configMapLookupKey, cm)
			if err != nil {
				return false
			}
			if cm.Data["stock_symbol"] == "GOOGL" {
				return true
			}
			return false
		}, timeout, interval).Should(BeTrue())

	})

})
