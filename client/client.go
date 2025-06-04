package client

import (
	"fmt"

	"github.com/pomadev/dollar-cost-averaging-bot/client/bitbank"
	"github.com/pomadev/dollar-cost-averaging-bot/client/bitflyer"
	"github.com/pomadev/dollar-cost-averaging-bot/config"
)

type client interface {
	OrderBTC(int64) (string, string, error)
	OrderETH(int64) (string, string, error)
}

func CreateClient(cfg *config.Config) (client, error) {
	switch cfg.Exchange {
	case "bitbank":
		return &bitbank.BitbankClient{AccessKey: cfg.AccessKey, ApiSecret: cfg.APISecret}, nil
	case "bitflyer":
		return &bitflyer.BitflyerClient{AccessKey: cfg.AccessKey, ApiSecret: cfg.APISecret}, nil
	default:
		return nil, fmt.Errorf("Unknown exchange: %s", cfg.Exchange)
	}
}
