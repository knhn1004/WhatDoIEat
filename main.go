package main

import (
	"context"
	"log"

	"github.com/knhn1004/WhatDoIEat/internal/bot"
	"github.com/knhn1004/WhatDoIEat/internal/config"
	"github.com/knhn1004/WhatDoIEat/internal/services"
	"github.com/knhn1004/WhatDoIEat/internal/web"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	bot.InitializeBot(config.BotToken, config.AppToken)
	services.StartServices()
	web.InitWeb()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Start(ctx)
}
