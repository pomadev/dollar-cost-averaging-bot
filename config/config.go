package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Exchange        string
	BTCAmount       int64
	ETHAmount       int64
	AccessKey       string
	APISecret       string
	LineSecret      string
	LineAccessToken string
	LineUserID      string
}

func Load() (*Config, error) {
	cfg := &Config{
		Exchange:        os.Getenv("EXCHANGE"),
		AccessKey:       os.Getenv("ACCESS_KEY"),
		APISecret:       os.Getenv("API_SECRET"),
		LineSecret:      os.Getenv("LINE_SECRET"),
		LineAccessToken: os.Getenv("LINE_ACCESS_TOKEN"),
		LineUserID:      os.Getenv("LINE_USER_ID"),
	}

	btc, err := strconv.ParseInt(os.Getenv("BTC_YEN"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BTC_YEN: %w", err)
	}
	cfg.BTCAmount = btc

	eth, err := strconv.ParseInt(os.Getenv("ETH_YEN"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ETH_YEN: %w", err)
	}
	cfg.ETHAmount = eth

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Exchange == "" {
		return fmt.Errorf("EXCHANGE is not set")
	}
	if c.Exchange != "bitbank" && c.Exchange != "bitflyer" {
		return fmt.Errorf("unknown exchange: %s", c.Exchange)
	}
	if c.AccessKey == "" || c.APISecret == "" {
		return fmt.Errorf("ACCESS_KEY or API_SECRET is not set")
	}
	if c.LineSecret == "" || c.LineAccessToken == "" || c.LineUserID == "" {
		return fmt.Errorf("LINE configuration is incomplete")
	}
	if c.BTCAmount <= 0 {
		return fmt.Errorf("BTC_YEN must be positive")
	}
	if c.ETHAmount <= 0 {
		return fmt.Errorf("ETH_YEN must be positive")
	}
	return nil
}