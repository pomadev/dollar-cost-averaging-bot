package client

import (
	"fmt"
	"os"

	"github.com/pomadev/dollar-cost-averaging-bot/client/bitbank"
	"github.com/pomadev/dollar-cost-averaging-bot/client/bitflyer"
)

type Client interface {
	OrderBTC(int64) (string, string, error)
	OrderETH(int64) (string, string, error)
}

type Pair string

const (
	BTCJPY Pair = "btc_jpy"
	ETHJPY Pair = "eth_jpy"
)

func CreateClient(exchange string) (Client, error) {
	accessKey := os.Getenv("ACCESS_KEY")
	apiSecret := os.Getenv("API_SECRET")
	if accessKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("ACCESS_KEY or API_SECRET is not set")
	}

	switch exchange {
	case "bitbank":
		return &bitbank.BitbankClient{BaseClient: BaseClient{AccessKey: accessKey, ApiSecret: apiSecret}}, nil
	case "bitflyer":
		return &bitflyer.BitflyerClient{BaseClient: BaseClient{AccessKey: accessKey, ApiSecret: apiSecret}}, nil
	default:
		return nil, fmt.Errorf("Unknown exchange: %s", exchange)
	}
}
