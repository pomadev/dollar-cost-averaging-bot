package p

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pomadev/dollar-cost-averaging-bot/client"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func Crypto(ctx context.Context, m PubSubMessage) error {
	exchange := os.Getenv("EXCHANGE")
	btc, err := strconv.ParseInt(os.Getenv("BTC_YEN"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse BTC_YEN: %s\n", err)
	}
	eth, err := strconv.ParseInt(os.Getenv("ETH_YEN"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse ETH_YEN: %s\n", err)
	}

	c, err := client.CreateClient(exchange)
	if err != nil {
		log.Fatalf("Failed to create client: %s\n", err)
	}
	err = c.OrderBTC(btc)
	if err != nil {
		log.Fatalf("Failed to order BTC: %s\n", err)
	}
	err = c.OrderETH(eth)
	if err != nil {
		log.Fatalf("Failed to order ETH: %s\n", err)
	}
	return nil
}
