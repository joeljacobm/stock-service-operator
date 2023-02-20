package backend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

const (
	APIKey         = "your_vantage_api_key"
)

type Vantage struct {
	StockReport
}

type VantageSearchResponse struct {
	GlobalQuote GlobalQuote `json:"Global Quote"`
}

type GlobalQuote struct {
	Price string `json:"05. price"`
}

func getVantageBackEnd(symbol string, log logr.Logger) *Vantage {

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, APIKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	return &Vantage{StockReport{Symbol: symbol, log: log, req: req}}
}

func (y *Vantage) IsValidSymbol() bool {
	response, err := y.fetchStock()
	if err != nil {
		return false
	}
	if response.GlobalQuote.Price == "" {
		return false
	}

	return true
}

func (y *Vantage) GetStockPrice() (string, error) {

	response, err := y.fetchStock()
	if err != nil {
		return "", nil
	}
	return response.GlobalQuote.Price, nil
}

func (v *Vantage) fetchStock() (VantageSearchResponse, error) {
	resp, err := v.client.Do(v.req)
	if err != nil {
		return VantageSearchResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return VantageSearchResponse{}, err
	}

	var result VantageSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return VantageSearchResponse{}, err
	}

	return result, nil
}
