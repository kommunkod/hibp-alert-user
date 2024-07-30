package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type LatestBreachResponse struct {
	Name               string   `json:"Name"`
	Title              string   `json:"Title"`
	Domain             string   `json:"Domain"`
	BreachDate         string   `json:"BreachDate"`
	AddedDate          string   `json:"AddedDate"`
	ModifiedDate       string   `json:"ModifiedDate"`
	PwnCount           int      `json:"PwnCount"`
	Description        string   `json:"Description"`
	LogoPath           string   `json:"LogoPath"`
	DataClasses        []string `json:"DataClasses"`
	IsVerified         bool     `json:"IsVerified"`
	IsFabricated       bool     `json:"IsFabricated"`
	IsSensitive        bool     `json:"IsSensitive"`
	IsRetired          bool     `json:"IsRetired"`
	IsSpamList         bool     `json:"IsSpamList"`
	IsMalware          bool     `json:"IsMalware"`
	IsSubscriptionFree bool     `json:"IsSubscriptionFree"`
}

func QueryLatestBreach() (LatestBreachResponse, error) {
	req, err := http.NewRequest("GET", "https://haveibeenpwned.com/api/v3/latestbreach", nil)
	if err != nil {
		log.Fatal(err)
		return LatestBreachResponse{}, err
	}

	req.Header.Set("User-Agent", "HIBP Breach Alert")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return LatestBreachResponse{}, err
	}

	defer resp.Body.Close()

	var latestBreachResponse LatestBreachResponse

	if err := json.NewDecoder(resp.Body).Decode(&latestBreachResponse); err != nil {
		log.Fatal(err)
		return LatestBreachResponse{}, err
	}

	return latestBreachResponse, nil
}
