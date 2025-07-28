package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"tetraaa/goback/finance"
	"time"

	"github.com/robfig/cron/v3"
)

func RegisterCronJobs() {
	crontab := cron.New()
	crontab.AddFunc("@hourly", UpdatePortfolioHistory)
	crontab.Start()
}

func UpdatePortfolioHistory() {
	portfolio_history, err := os.ReadFile("data/portfolio_history.json")
	if err != nil {
		fmt.Println("Unable to update portfolio history : json file does not exist.")
		return
	}

	newValue, err := finance.GetUpdatedPortfolioValue()
	if err != nil {
		fmt.Println("Unable to update portfolio history : unable to compute current portfolio value")
		return
	}
	response := finance.PortfolioHistory{}
	json.Unmarshal(portfolio_history, &response)
	response.PortfolioHistory = append(response.PortfolioHistory, finance.PortfolioElement{Timestamp: time.Now().UnixMilli(), PortfolioValue: *newValue})
	updatedFileContent, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Unable to update portfolio history : file content was invalid.")
		return
	}
	err = os.WriteFile("data/portfolio_history.json", updatedFileContent, 0644)
	if err != nil {
		fmt.Println("Unable to update portfolio history : file update failed")
		return
	}

}
