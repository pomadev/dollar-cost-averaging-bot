package p

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
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
	bamount, bprice, err := c.OrderBTC(btc)
	if err != nil {
		log.Fatalf("Failed to order BTC: %s\n", err)
	}
	eamount, eprice, err := c.OrderETH(eth)
	if err != nil {
		log.Fatalf("Failed to order ETH: %s\n", err)
	}
	bot, err := linebot.New(os.Getenv("LINE_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create client: %s\n", err)
	}
	if _, err := bot.PushMessage(os.Getenv("LINE_USER_ID"), linebot.NewTextMessage(fmt.Sprintf("以下を購入しました。\n\n[BTC]\n購入価格: %s円\n購入数量: %sBTC\n\n[ETH]\n購入価格: %s円\n購入数量: %sETH", bprice, bamount, eprice, eamount))).Do(); err != nil {
		log.Fatalf("Failed to push message: %s\n", err)
	}
	return nil
}
