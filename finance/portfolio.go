package finance

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PortfolioHistory struct {
	PortfolioHistory []PortfolioElement `json:"portfolio_history"`
}

type PortfolioElement struct {
	Timestamp      int64   `json:"timestamp"`
	PortfolioValue float64 `json:"portfolioValue"`
}

type PortfolioPositions []struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Position string `json:"position"`
	Isin     string `json:"isin,omitempty"`
	Ticker   string `json:"ticker,omitempty"`
}

type PortfolioResponse struct {
	PortfolioPositions PortfolioPositions `json:"portfolio"`
}

type StockDataMeta struct {
	Currency           string  `json:"currency,omitempty"`
	Symbol             string  `json:"symbol,omitempty"`
	InstrumentType     string  `json:"instrumentType,omitempty"`
	RegularMarketPrice float64 `json:"regularMarketPrice,omitempty"`
	LongName           string  `json:"longName,omitempty"`
	ShortName          string  `json:"shortName,omitempty"`
}

type StockDataResponse struct {
	Chart struct {
		Result []struct {
			Meta StockDataMeta `json:"meta,omitempty"`
		} `json:"result,omitempty"`
		Error any `json:"error,omitempty"`
	} `json:"chart,omitempty"`
}

func GetPortfolio() (*PortfolioResponse, error) {
	portfolio, err := os.ReadFile("data/portfolio.json")
	if err != nil {
		return nil, err
	}

	response := PortfolioResponse{}
	json.Unmarshal(portfolio, &response)

	return &response, nil
}

func GetStockInformationFromTicker(ticker string) (*StockDataMeta, error) {
	if ticker == "" {
		return nil, errors.New("no ticker provided")
	}
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?metrics=high?&interval=1d", ticker)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36`)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	stockData := StockDataResponse{}

	json.NewDecoder(resp.Body).Decode(&stockData)
	if len(stockData.Chart.Result) <= 0 {
		fmt.Printf("Unable to get result data for ticker %s", ticker)
		return nil, err
	}
	return &stockData.Chart.Result[0].Meta, nil
}

func GetUpdatedPortfolioValue() (*float64, error) {
	portfolio, err := GetPortfolio()
	if err != nil {
		return nil, err
	}

	var updatedPortfolioValue float64 = 0

	for _, pos := range portfolio.PortfolioPositions {
		currentPosition, _ := strconv.ParseFloat(strings.Split(pos.Position, "@")[0], 64)
		avgPrice, _ := strconv.ParseFloat(strings.Split(pos.Position, "@")[1], 64)
		currentStockInfo, err := GetStockInformationFromTicker(pos.Ticker)
		if err != nil {
			updatedPortfolioValue += currentPosition * avgPrice
		} else {
			updatedPortfolioValue += currentPosition * currentStockInfo.RegularMarketPrice
		}
	}
	return &updatedPortfolioValue, nil
}
