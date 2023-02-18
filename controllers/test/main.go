package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	Raw float64 `json:"raw`
}

func main() {
	symbol := "GOOGL" // Replace with the stock symbol you want to check

	url := fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=price", symbol)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Invalid stock symbol or error with API")
		return
	}

	var yahooFinanceResponse YahooFinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooFinanceResponse); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if len(yahooFinanceResponse.QuoteSummary.Result) == 0 {
		fmt.Println("Invalid stock symbol")
		return
	}

	price := yahooFinanceResponse.QuoteSummary.Result[0].Price
	fmt.Printf("The stock price of %s is $%.2f", symbol, price.RegularMarketPrice.Raw)
}
