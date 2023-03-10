package bitflyer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pomadev/dollar-cost-averaging-bot/client/common"
)

const (
	API_URL = "https://api.bitflyer.com/v1/"
	BTC_JPY = "BTC_JPY"
	ETH_JPY = "ETH_JPY"
)

type BitflyerClient struct {
	AccessKey string
	ApiSecret string
	nonce     int64
}

func (c *BitflyerClient) getNonce() string {
	if c.nonce == 0 {
		c.nonce = time.Now().Unix()
	}
	c.nonce++
	return fmt.Sprintf("%d", c.nonce)
}

func (c *BitflyerClient) OrderBTC(yen int64) error {
	return c.order(BTC_JPY, yen)
}

func (c *BitflyerClient) OrderETH(yen int64) error {
	return c.order(ETH_JPY, yen)
}

type orderRequest struct {
	ProductCode    string  `json:"product_code"`
	ChildOrderType string  `json:"child_order_type"`
	Side           string  `json:"side"`
	Size           float64 `json:"size"`
}

func (c *BitflyerClient) order(pair string, yen int64) error {
	price, err := getPrice(pair)
	if err != nil {
		return fmt.Errorf("Failed to get price: %s", err)
	}
	amount := common.CalcAmount(price, yen, 100000000)
	requestBody := orderRequest{
		ProductCode:    pair,
		ChildOrderType: "MARKET",
		Side:           "BUY",
		Size:           amount,
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("Failed to marshal request body: %s", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/me/sendchildorder", API_URL), bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return fmt.Errorf("Failed to create request: %s", err)
	}
	nonce := c.getNonce()
	req.Header.Set("ACCESS-KEY", c.AccessKey)
	req.Header.Set("ACCESS-TIMESTAMP", nonce)
	req.Header.Set("ACCESS-SIGN", common.MakeSign(nonce+"POST"+"/v1/me/sendchildorder"+string(requestBodyJson), c.ApiSecret))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send request: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to send request HTTP Status Code: %s", res.Status)
	}
	return nil
}

type tickerResponse struct {
	State string  `json:"state"`
	Ltp   float64 `json:"ltp"`
}

func getPrice(pair string) (int64, error) {
	res, err := http.Get(fmt.Sprintf("%s/ticker?product_code=%s", API_URL, pair))
	if err != nil {
		return 0, fmt.Errorf("Failed to get ticker: %s", err)
	}
	defer res.Body.Close()

	var ticker tickerResponse
	if err := json.NewDecoder(res.Body).Decode(&ticker); err != nil {
		return 0, fmt.Errorf("Failed to decode ticker: %s", err)
	}
	if ticker.State != "RUNNING" {
		return 0, fmt.Errorf("Ticker state is not RUNNING: %s", ticker.State)
	}

	return int64(ticker.Ltp), nil
}
