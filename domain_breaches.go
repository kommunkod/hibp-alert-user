package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DomainBreachesResponse map[string][]string

func QueryDomainBreaches(domain string, config Config) (DomainBreachesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://haveibeenpwned.com/api/v3/breacheddomain/%s", domain), nil)
	if err != nil {
		log.Fatal(err)
		return DomainBreachesResponse{}, err
	}

	req.Header.Set("User-Agent", "HIBP Breach Alert")
	req.Header.Set("hibp-api-key", config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return DomainBreachesResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode != 429 {
			log.Fatalf("Received status code %d", resp.StatusCode)
			return DomainBreachesResponse{}, fmt.Errorf("Received status code %d", resp.StatusCode)
		}

		waitFor := resp.Header.Get("Retry-After")
		log.Printf("Rate limited. Waiting for %s seconds", waitFor)
		ival, err := strconv.ParseInt(waitFor, 10, 64)
		if err != nil {
			log.Fatal(err)
			return DomainBreachesResponse{}, err
		}
		time.Sleep(time.Second * time.Duration(ival+10))

		return QueryDomainBreaches(domain, config)
	}

	var domainBreachResponse DomainBreachesResponse

	if err := json.NewDecoder(resp.Body).Decode(&domainBreachResponse); err != nil {
		log.Fatal(err)
		return DomainBreachesResponse{}, err
	}

	return domainBreachResponse, nil
}
