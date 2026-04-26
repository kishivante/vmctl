package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
	Node    string `json:"node"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
