package backend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

type Yahoo struct {
	StockReport
}
type YahooFinanceResponse struct {
	QuoteSummary struct {
		Result []struct {
			Price struct {
				RegularMarketPrice RegularMarketPrice `json:"regularMarketPrice"`
			} `json:"price"`
		} `json:"result"`
	} `json:"quoteSummary"`
}

type RegularMarketPrice struct {
	Raw float64 `json:"raw"`
}

func getYahooBackEnd(symbol string, log logr.Logger) *Yahoo {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=price", symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	return &Yahoo{StockReport{Symbol: symbol, log: log, req: req, client: http.Client{}}}
}

func (y *Yahoo) IsValidSymbol() bool {

	resp, err := y.fetchStock()
	if err != nil {
		return false
	}

	if len(resp.QuoteSummary.Result) == 0 {
		return false
	}

	return true
}

func (y *Yahoo) GetStockPrice() (float64, error) {

	resp, err := y.fetchStock()
	if err != nil {
		return 0, err
	}
	return resp.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw, nil

}

func (y *Yahoo) fetchStock() (YahooFinanceResponse, error) {

	resp, err := y.client.Do(y.req)
	if err != nil {
		return YahooFinanceResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return YahooFinanceResponse{}, fmt.Errorf("invalid stock symbol or error with API")
	}

	var yahooFinanceResponse YahooFinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooFinanceResponse); err != nil {
		return YahooFinanceResponse{}, fmt.Errorf("error decoding response: %s", err.Error())
	}
	return yahooFinanceResponse, nil
}
