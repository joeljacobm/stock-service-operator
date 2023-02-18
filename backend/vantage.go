package backend

import "github.com/go-logr/logr"

const (
	VantageBaseUrl = "https://www.alphavantage.co/query?function"
	APIKey         = ""
)

type Vantage struct {
	StockReport
}

func getVantageBackEnd(symbol string, log logr.Logger) *Vantage {
	return &Vantage{StockReport{Symbol: symbol, log: log}}
}

func (y *Vantage) IsValidSymbol() bool {

	return false
}

func (y *Vantage) GetStockPrice() (string, error) {

	return "", nil
}
