package dca

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/pomadev/dollar-cost-averaging-bot/client/bitbank"
)

type Config struct {
	BitbankAccessKey       string
	BitbankAPISecret       string
	LineChannelSecret      string
	LineChannelAccessToken string
	LineUserID             string
	BTCCostYen             int64
	DryRun                 bool
}

func Run(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	exchangeClient := bitbank.NewClient(cfg.BitbankAccessKey, cfg.BitbankAPISecret, httpClient)
	order, err := exchangeClient.BuyBTC(ctx, cfg.BTCCostYen, cfg.DryRun)
	if err != nil {
		return fmt.Errorf("buy BTC: %w", err)
	}

	bot, err := linebot.New(cfg.LineChannelSecret, cfg.LineChannelAccessToken, linebot.WithHTTPClient(httpClient))
	if err != nil {
		return fmt.Errorf("create LINE client: %w", err)
	}

	messagePrefix := "BTCを積み立てました。"
	if order.DryRun {
		messagePrefix = "DRY RUN: BTCの積み立て予定を計算しました。"
	}
	message := fmt.Sprintf("%s\n\n[BTC]\n価格: %d円\n数量: %.4fBTC\n投入額: %d円", messagePrefix, order.PriceYen, order.AmountBTC, cfg.BTCCostYen)
	if _, err := bot.PushMessage(cfg.LineUserID, linebot.NewTextMessage(message)).Do(); err != nil {
		return fmt.Errorf("push LINE message: %w", err)
	}
	return nil
}

func loadConfig() (Config, error) {
	btcCostYen, err := strconv.ParseInt(os.Getenv("BTC_YEN"), 10, 64)
	if err != nil {
		return Config{}, fmt.Errorf("parse BTC_YEN: %w", err)
	}
	if btcCostYen <= 0 {
		return Config{}, errors.New("BTC_YEN must be greater than 0")
	}

	dryRun, err := strconv.ParseBool(defaultValue(os.Getenv("DRY_RUN"), "false"))
	if err != nil {
		return Config{}, fmt.Errorf("parse DRY_RUN: %w", err)
	}

	cfg := Config{
		BitbankAccessKey:       strings.TrimSpace(os.Getenv("BITBANK_ACCESS_KEY")),
		BitbankAPISecret:       strings.TrimSpace(os.Getenv("BITBANK_API_SECRET")),
		LineChannelSecret:      strings.TrimSpace(os.Getenv("LINE_CHANNEL_SECRET")),
		LineChannelAccessToken: strings.TrimSpace(os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")),
		LineUserID:             strings.TrimSpace(os.Getenv("LINE_USER_ID")),
		BTCCostYen:             btcCostYen,
		DryRun:                 dryRun,
	}

	switch {
	case cfg.BitbankAccessKey == "":
		return Config{}, errors.New("BITBANK_ACCESS_KEY is not set")
	case cfg.BitbankAPISecret == "":
		return Config{}, errors.New("BITBANK_API_SECRET is not set")
	case cfg.LineChannelSecret == "":
		return Config{}, errors.New("LINE_CHANNEL_SECRET is not set")
	case cfg.LineChannelAccessToken == "":
		return Config{}, errors.New("LINE_CHANNEL_ACCESS_TOKEN is not set")
	case cfg.LineUserID == "":
		return Config{}, errors.New("LINE_USER_ID is not set")
	}

	return cfg, nil
}

func defaultValue(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
