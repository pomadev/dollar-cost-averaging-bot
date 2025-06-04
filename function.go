package p

import (
	"context"
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/pomadev/dollar-cost-averaging-bot/client"
	"github.com/pomadev/dollar-cost-averaging-bot/config"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func Crypto(ctx context.Context, m PubSubMessage) error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}

	c, err := client.CreateClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %s\n", err)
	}
	bamount, bprice, err := c.OrderBTC(cfg.BTCAmount)
	if err != nil {
		log.Fatalf("Failed to order BTC: %s\n", err)
	}
	eamount, eprice, err := c.OrderETH(cfg.ETHAmount)
	if err != nil {
		log.Fatalf("Failed to order ETH: %s\n", err)
	}
	bot, err := linebot.New(cfg.LineSecret, cfg.LineAccessToken)
	if err != nil {
		log.Fatalf("Failed to create client: %s\n", err)
	}
	if _, err := bot.PushMessage(cfg.LineUserID, linebot.NewTextMessage(fmt.Sprintf("以下を購入しました。\n\n[BTC]\n価格: %s円\n数量: %sBTC\n\n[ETH]\n価格: %s円\n数量: %sETH", bprice, bamount, eprice, eamount))).Do(); err != nil {
		log.Fatalf("Failed to push message: %s\n", err)
	}
	return nil
}
