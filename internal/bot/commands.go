package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"

	"github.com/knhn1004/WhatDoIEat/internal/config"
	"github.com/knhn1004/WhatDoIEat/internal/models"
	"github.com/knhn1004/WhatDoIEat/internal/services"
)

// Command definitions
var Commands = []*slacker.CommandDefinition{
	{
		Command: "ping",
		Handler: handlePing,
		Middlewares: []slacker.CommandMiddlewareHandler{
			logUserMiddleware(),
		},
	},
	{
		Command: "upload <sentence>",
		Handler: handleUpload,
	},
	{
		Command: "echo {word}",
		Handler: handleEcho,
	},
	{
		Command: "find <query>",
		Handler: handleFindRestaurant,
		Middlewares: []slacker.CommandMiddlewareHandler{
			logUserMiddleware(),
		},
	},
	{
		Command: "pref",
		Handler: handlePref,
		Middlewares: []slacker.CommandMiddlewareHandler{
			logUserMiddleware(),
		},
	},
	{
		Command: "recipe",
		Handler: handleRecipe,
		Middlewares: []slacker.CommandMiddlewareHandler{
			logUserMiddleware(),
		},
	},
}

func handlePing(ctx *slacker.CommandContext) {
	t1, _ := ctx.Response().Reply("about to be replaced 🚧️")

	time.Sleep(time.Second)

	ctx.Response().Reply("pong🏓️", slacker.WithReplace(t1))
}

func handleUpload(ctx *slacker.CommandContext) {
	sentence := ctx.Request().Param("sentence")

	api := ctx.SlackClient()
	event := ctx.Event()

	api.PostMessage(event.ChannelID, slack.MsgOptionText("📄️ Uploading file ...", false))

	_, err := api.UploadFile(slack.FileUploadParameters{
		Content:  sentence,
		Filename: "sentence.txt",
		Channels: []string{event.ChannelID},
	})
	if err != nil {
		ctx.Response().ReplyError(err)
	}
}

func handleEcho(ctx *slacker.CommandContext) {
	word := ctx.Request().Param("word")

	attachments := []slack.Attachment{
		{
			Color:      "good",
			AuthorName: "👨️ Raed Shomali",
			Title:      "Attachment Title",
			Text:       "Attachment Text",
		},
	}

	ctx.Response().Reply(word, slacker.WithAttachments(attachments))
}

func handleFindRestaurant(ctx *slacker.CommandContext) {
	query := ctx.Request().Param("query")

	// 1. Initial message
	t1, _ := ctx.Response().Reply("Searching for restaurants based on your query: " + query + " 🍔...")

	// Use a channel to retrieve the results from the goroutine
	resultChannel := make(chan []models.Restaurant, 1)
	errorChannel := make(chan error, 1)

	go func() {
		options := map[string]string{
			"radius":  "10000", // 10km
			"sort_by": "rating",
			"limit":   "5",
		}
		businesses, err := services.GetRestaurants(query, "", options)
		if err != nil {
			errorChannel <- err
			return
		}
		resultChannel <- businesses
	}()

	select {
	case businesses := <-resultChannel:
		var attachments []slack.Attachment
		for _, business := range businesses {
			// Convert the Location structure into a readable string
			locationString := formatLocation(business.Location)

			attachment := slack.Attachment{
				Color:    "good",
				Title:    business.Name,
				Text:     "Phone: " + business.Phone + "\nURL: " + business.URL + "\nRating: " + strconv.FormatFloat(business.Rating, 'f', 1, 64) + "\nLocation: " + locationString + "\nMap: " + services.GenGoogleMapsURL(business.Location) + "\nImage: " + business.ImageURL,
				ImageURL: business.ImageURL,
			}
			attachments = append(attachments, attachment)
		}
		ctx.Response().Reply("Here are some restaurants we found for you based on your query: "+query, slacker.WithReplace(t1), slacker.WithAttachments(attachments))
	case err := <-errorChannel:
		ctx.Response().Reply("Error finding restaurants: "+err.Error(), slacker.WithReplace(t1))
	}
}

func formatLocation(loc models.Location) string {
	address := loc.Address1
	if loc.Address2 != "" {
		address += ", " + loc.Address2
	}
	if loc.Address3 != "" {
		address += ", " + loc.Address3
	}
	return address + ", " + loc.City + ", " + loc.ZipCode + ", " + loc.State + ", " + loc.Country
}

func handlePref(ctx *slacker.CommandContext) {
	userId := ctx.Event().UserID

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(24 * time.Hour).Unix(), // Token expiration: 24 hours from now
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.AppToken))
	if err != nil {
		log.Fatalf("Error signing token: %v", err)
		return
	}

	preferencesEndpoint := "/pref"
	fullURL := config.ServerAddr + preferencesEndpoint

	message := slack.NewTextBlockObject(slack.MarkdownType, "Click the link below to set your preferences:\n"+fullURL+"?token="+tokenString, false, false)
	blocks := []slack.Block{}

	blocks = append(blocks, slack.NewContextBlock("preferences_block", message))

	// Send the constructed block message to Slack.
	ctx.Response().ReplyBlocks(blocks)
}

func handleRecipe(ctx *slacker.CommandContext) {
	// 1. Initial message
	t1, _ := ctx.Response().Reply("Generating your daily recipe 🍳...")

	// Use a channel to retrieve the results from the goroutine
	resultChannel := make(chan []models.Recipe, 1)
	errorChannel := make(chan error, 1)

	go func() {
		recipes, err := services.GenRecipeOpenAI()
		if err != nil {
			errorChannel <- err
			return
		}
		resultChannel <- recipes
	}()

	select {
	case recipes := <-resultChannel:
		var attachments []slack.Attachment
		for _, recipe := range recipes {
			// Convert the Recipe data into readable text
			text := "Name: " + recipe.Name + "\n" +
				"Short Description: " + recipe.ShortDescription + "\n" +
				"Ingredients: " + formatIngredients(recipe.Ingredients) + "\n" +
				"Steps: " + strings.Join(recipe.Steps, ", ")

			imgUrl, err := services.GetImageUrlFromText(recipe.Meal)
			if err != nil {
				fmt.Println(err)
			}

			attachment := slack.Attachment{
				Color:    "good",
				Title:    "Meal: " + recipe.Meal,
				Text:     text,
				ImageURL: imgUrl,
			}
			attachments = append(attachments, attachment)
		}
		ctx.Response().Reply("Here are the recipes we generated for you:", slacker.WithReplace(t1), slacker.WithAttachments(attachments))

	case err := <-errorChannel:
		ctx.Response().Reply("Error generating recipe: "+err.Error(), slacker.WithReplace(t1))
	}
}

func formatIngredients(ingredients []models.Ingredient) string {
	var formattedIngredients []string
	for _, ingredient := range ingredients {
		formattedIngredients = append(formattedIngredients, ingredient.Quantity+" of "+ingredient.Ingredient)
	}
	return strings.Join(formattedIngredients, ", ")
}
