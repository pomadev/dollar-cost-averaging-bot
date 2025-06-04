package client

import (
	"fmt"
	"time"
)

type BaseClient struct {
	AccessKey string
	ApiSecret string
	nonce     int64
}

func (c *BaseClient) GetNonce() string {
	if c.nonce == 0 {
		c.nonce = time.Now().Unix()
	}
	c.nonce++
	return fmt.Sprintf("%d", c.nonce)
}