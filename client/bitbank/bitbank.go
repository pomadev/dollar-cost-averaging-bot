package bitbank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pomadev/dollar-cost-averaging-bot/client/common"
)

const (
	PUBLIC_API_URL  = "https://public.bitbank.cc"
	PRIVATE_API_URL = "https://api.bitbank.cc/v1"
	BTC_JPY         = "btc_jpy"
	ETH_JPY         = "eth_jpy"
)

type BitbankClient struct {
	AccessKey string
	ApiSecret string
	nonce     int64
}

func (c *BitbankClient) getNonce() string {
	if c.nonce == 0 {
		c.nonce = time.Now().Unix()
	}
	c.nonce++
	return fmt.Sprintf("%d", c.nonce)
}

func (c *BitbankClient) OrderBTC(yen int64) (string, string, error) {
	return c.order(BTC_JPY, yen)
}

func (c *BitbankClient) OrderETH(yen int64) (string, string, error) {
	return c.order(ETH_JPY, yen)
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
		Code           int    `json:"code"`
		ExecutedAmount string `json:"executed_amount"`
		AveragePrice   string `json:"average_price"`
	} `json:"data"`
}

func (c *BitbankClient) order(pair string, yen int64) (string, string, error) {
	price, err := getPrice(pair)
	if err != nil {
		return "", "", fmt.Errorf("Failed to get price: %s", err)
	}
	amount := common.CalcAmount(price, yen, 10000)
	if amount == 0 {
		log.Print("Amount is too small")
		return "", "", nil
	}
	requestBody := orderRequest{
		Pair:   pair,
		Amount: strconv.FormatFloat(amount, 'f', 4, 64),
		Side:   "buy",
		Type:   "market",
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", fmt.Errorf("Failed to marshal request body: %s", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/spot/order", PRIVATE_API_URL), bytes.NewReader(requestBodyJson))
	if err != nil {
		return "", "", fmt.Errorf("Failed to create request: %s", err)
	}
	nonce := c.getNonce()
	req.Header.Set("ACCESS-KEY", c.AccessKey)
	req.Header.Set("ACCESS-NONCE", nonce)
	req.Header.Set("ACCESS-SIGNATURE", common.MakeSign(nonce+string(requestBodyJson), c.ApiSecret))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("Failed to order: %s", err)
	}
	defer res.Body.Close()

	var order orderResponse
	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return "", "", fmt.Errorf("Failed to decode order: %s", err)
	}
	if order.Success != 1 {
		return "", "", fmt.Errorf("Failed to order success code: %d", order.Data.Code)
	}
	return order.Data.ExecutedAmount, order.Data.AveragePrice, nil
}

type tickerResponse struct {
	Success int `json:"success"`
	Data    struct {
		Buy  string `json:"buy"`
		Code int    `json:"code"`
	} `json:"data"`
}

func getPrice(pair string) (int64, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s/ticker", PUBLIC_API_URL, pair))
	if err != nil {
		return 0, fmt.Errorf("Failed to get ticker: %s", err)
	}
	defer res.Body.Close()

	var ticker tickerResponse
	if err := json.NewDecoder(res.Body).Decode(&ticker); err != nil {
		return 0, fmt.Errorf("Failed to decode ticker: %s", err)
	}
	if ticker.Success != 1 {
		return 0, fmt.Errorf("Failed to get ticker success code: %d", ticker.Data.Code)
	}
	p, err := strconv.ParseInt(ticker.Data.Buy, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse buy price: %s", err)
	}
	return p, nil
}
