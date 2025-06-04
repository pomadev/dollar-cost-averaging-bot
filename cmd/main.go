package main

import p "github.com/pomadev/dollar-cost-averaging-bot"

func main() {
	if err := p.Crypto(nil, p.PubSubMessage{}); err != nil {
		panic(err)
	}
}
