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

	/* recipes, err := services.GenRecipeOpenAI()
	if err != nil {
		fmt.Printf("genRecipeOpenAI error: %v\n", err)
	} else {
		for _, recipe := range recipes {
			fmt.Println("Meal: ", recipe.Meal)
			fmt.Println("Name: ", recipe.Name)
			fmt.Println("Short Description: ", recipe.ShortDescription)
			fmt.Println("Ingredients: ", recipe.Ingredients)
			fmt.Println("Steps: ", recipe.Steps)
		}
	} */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Start(ctx)
}
