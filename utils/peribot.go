package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetPeribotStatus() (PeribotResponse, error) {
	fallbackValue := PeribotResponse{Status: "offline", Uptime: 0, MemoryUsed: 0, CachedChannels: 0, TotalAttachments: 0}
	resp, err := http.Get("http://localhost:3000")
	if err != nil {
		return fallbackValue, err
	}
	var peribotResponse PeribotResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fallbackValue, err
	}
	err = json.Unmarshal(bodyBytes, &peribotResponse)
	if err != nil {
		return fallbackValue, err
	}
	return peribotResponse, nil
}
