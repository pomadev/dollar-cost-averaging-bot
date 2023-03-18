package client

import (
	"fmt"
	"os"

	"github.com/pomadev/dollar-cost-averaging-bot/client/bitbank"
	"github.com/pomadev/dollar-cost-averaging-bot/client/bitflyer"
)

type client interface {
	OrderBTC(int64) (string, string, error)
	OrderETH(int64) (string, string, error)
}

func CreateClient(exchange string) (client, error) {
	accessKey := os.Getenv("ACCESS_KEY")
	apiSecret := os.Getenv("API_SECRET")
	if accessKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("ACCESS_KEY or API_SECRET is not set")
	}

	switch exchange {
	case "bitbank":
		return &bitbank.BitbankClient{AccessKey: accessKey, ApiSecret: apiSecret}, nil
	case "bitflyer":
		return &bitflyer.BitflyerClient{AccessKey: accessKey, ApiSecret: apiSecret}, nil
	default:
		return nil, fmt.Errorf("Unknown exchange: %s", exchange)
	}
}
