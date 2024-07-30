package main

import (
	"encoding/json"
	"os"
)

type SmtpConfig struct {
	Sender   string `json:"sender"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Secure   bool   `json:"secure"`
	StartTLS bool   `json:"starttls"`
}

type IgnoreConfig struct {
	Unverified bool `json:"unverified"`
	Fabricated bool `json:"fabricated"`
	Sensitive  bool `json:"sensitive"`
	SpamList   bool `json:"spamList"`
	Malware    bool `json:"malware"`
	Retired    bool `json:"retired"`
}

type EmailConfig struct {
	Colors struct {
		Background string `json:"background"`
		Text       string `json:"text"`
	} `json:"colors"`

	Subject string `json:"subject"`
	Body    struct {
		Header              string   `json:"header"`
		Texts               []string `json:"texts"`
		PreviousBreachTexts []string `json:"previous_breach_texts"`
	} `json:"body"`
}

type Config struct {
	ApiKey           string       `json:"hibpApiKey"`
	Domains          []string     `json:"domains"`
	LatestBreach     string       `json:"latestBreach"`
	NotifiedBreaches []string     `json:"notifiedBreaches"`
	Smtp             SmtpConfig   `json:"smtp"`
	Ignore           IgnoreConfig `json:"ignore"`
	Email            EmailConfig  `json:"email"`
}

func LoadConfig() (Config, error) {
	config, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(config, &cfg)
	return cfg, err
}

func SaveConfig(cfg Config) error {
	config, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile("config.json", config, 0644)
	return err
}
