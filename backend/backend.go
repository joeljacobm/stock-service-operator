package backend

import (
	"net/http"

	"github.com/go-logr/logr"
)

type StockReport struct {
	log    logr.Logger
	Symbol string
	req    *http.Request
	client http.Client
}

type Backend interface {
	IsValidSymbol() bool
	GetStockPrice() (float64, error)
}

func GetBackend(backend string, symbol string, log logr.Logger) Backend {
	if backend == "vantage" {
		return getVantageBackEnd(symbol, log)
	}
	return getYahooBackEnd(symbol, log)
}
