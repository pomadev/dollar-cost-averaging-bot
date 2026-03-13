package bitbank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pomadev/dollar-cost-averaging-bot/client/common"
)

const (
	PUBLIC_API_URL  = "https://public.bitbank.cc"
	PRIVATE_API_URL = "https://api.bitbank.cc/v1"
	BTC_JPY         = "btc_jpy"
)

type Client struct {
	AccessKey  string
	ApiSecret  string
	HTTPClient *http.Client
	nonce      int64
}

type OrderResult struct {
	AmountBTC float64
	PriceYen  int64
	DryRun    bool
}

func NewClient(accessKey, apiSecret string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{
		AccessKey:  accessKey,
		ApiSecret:  apiSecret,
		HTTPClient: httpClient,
	}
}

func (c *Client) getNonce() string {
	if c.nonce == 0 {
		c.nonce = time.Now().Unix()
	}
	c.nonce++
	return fmt.Sprintf("%d", c.nonce)
}

type orderRequest struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
	Side   string `json:"side"`
	Type   string `json:"type"`
}

type orderResponse struct {
	Success int `json:"success"`
	Data    struct {
		Code int `json:"code"`
	} `json:"data"`
}

func (c *Client) BuyBTC(ctx context.Context, yen int64, dryRun bool) (OrderResult, error) {
	price, err := c.getPrice(ctx, BTC_JPY)
	if err != nil {
		return OrderResult{}, fmt.Errorf("get BTC price: %w", err)
	}
	amount := common.CalcAmount(price, yen, 10000)
	if amount == 0 {
		return OrderResult{}, fmt.Errorf("amount is too small for BTC_YEN=%d at price=%d", yen, price)
	}
	if dryRun {
		return OrderResult{
			AmountBTC: amount,
			PriceYen:  price,
			DryRun:    true,
		}, nil
	}

	requestBody := orderRequest{
		Pair:   BTC_JPY,
		Amount: strconv.FormatFloat(amount, 'f', 4, 64),
		Side:   "buy",
		Type:   "market",
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return OrderResult{}, fmt.Errorf("marshal request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/user/spot/order", PRIVATE_API_URL), bytes.NewReader(requestBodyJson))
	if err != nil {
		return OrderResult{}, fmt.Errorf("create request: %w", err)
	}
	nonce := c.getNonce()
	req.Header.Set("ACCESS-KEY", c.AccessKey)
	req.Header.Set("ACCESS-NONCE", nonce)
	req.Header.Set("ACCESS-SIGNATURE", common.MakeSign(nonce+string(requestBodyJson), c.ApiSecret))
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return OrderResult{}, fmt.Errorf("submit order: %w", err)
	}
	defer res.Body.Close()

	var order orderResponse
	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return OrderResult{}, fmt.Errorf("decode order: %w", err)
	}
	if order.Success != 1 {
		return OrderResult{}, fmt.Errorf("bitbank order failed with code: %d", order.Data.Code)
	}
	return OrderResult{
		AmountBTC: amount,
		PriceYen:  price,
		DryRun:    false,
	}, nil
}

type tickerResponse struct {
	Success int `json:"success"`
	Data    struct {
		Buy  string `json:"buy"`
		Code int    `json:"code"`
	} `json:"data"`
}

func (c *Client) getPrice(ctx context.Context, pair string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/ticker", PUBLIC_API_URL, pair), nil)
	if err != nil {
		return 0, fmt.Errorf("create ticker request: %w", err)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("get ticker: %w", err)
	}
	defer res.Body.Close()

	var ticker tickerResponse
	if err := json.NewDecoder(res.Body).Decode(&ticker); err != nil {
		return 0, fmt.Errorf("decode ticker: %w", err)
	}
	if ticker.Success != 1 {
		return 0, fmt.Errorf("bitbank ticker failed with code: %d", ticker.Data.Code)
	}
	p, err := strconv.ParseInt(ticker.Data.Buy, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse buy price: %w", err)
	}
	return p, nil
}
