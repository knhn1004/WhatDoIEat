package main

import (
	"context"
	"log"

	"github.com/knhn1004/WhatDoIEat/internal/bot"
	"github.com/knhn1004/WhatDoIEat/internal/config"
	"github.com/knhn1004/WhatDoIEat/internal/services"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	bot.InitializeBot(config.BotToken, config.AppToken)
	services.StartServices()

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

	/* var businesses []models.Restaurant
	options := map[string]string{
		"radius":  "10000", // 10km
		"sort_by": "rating",
		"limit":   "5",
	}
	businesses, err = services.GetRestaurants("Japanese", "", options)
	if err != nil {
		fmt.Printf("getRestaurants error: %v\n", err)
	} else {
		for _, business := range businesses {
			fmt.Println("Name: ", business.Name)
			fmt.Println("Phone: ", business.Phone)
			fmt.Println("URL: ", business.URL)
			fmt.Println("Rating: ", business.Rating)
			fmt.Println("Location: ", business.Location)
			fmt.Println("Map: ", services.GenGoogleMapsURL(business.Location))
			fmt.Println("Image: ", business.ImageURL)
		}
	} */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Start(ctx)
}
