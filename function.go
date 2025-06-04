package p

import (
	"context"
	"fmt"
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
		return fmt.Errorf("failed to parse BTC_YEN: %w", err)
	}
	eth, err := strconv.ParseInt(os.Getenv("ETH_YEN"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse ETH_YEN: %w", err)
	}

	c, err := client.CreateClient(exchange)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	bamount, bprice, err := c.OrderBTC(btc)
	if err != nil {
		return fmt.Errorf("failed to order BTC: %w", err)
	}
	eamount, eprice, err := c.OrderETH(eth)
	if err != nil {
		return fmt.Errorf("failed to order ETH: %w", err)
	}
	bot, err := linebot.New(os.Getenv("LINE_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		return fmt.Errorf("failed to create LINE bot client: %w", err)
	}
	if _, err := bot.PushMessage(os.Getenv("LINE_USER_ID"), linebot.NewTextMessage(fmt.Sprintf("以下を購入しました。\n\n[BTC]\n価格: %s円\n数量: %sBTC\n\n[ETH]\n価格: %s円\n数量: %sETH", bprice, bamount, eprice, eamount))).Do(); err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}
	return nil
}
