package main

import (
	"context"
	"log"
	"os"

	dca "github.com/pomadev/dollar-cost-averaging-bot"
)

func main() {
	if err := dca.Run(context.Background()); err != nil {
		log.Printf("job failed: %v", err)
		os.Exit(1)
	}
}
