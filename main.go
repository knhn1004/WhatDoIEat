package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/minoplhy/duckduckgo-images-api"
	supa "github.com/nedpals/supabase-go"

	"github.com/knhn1004/WhatDoIEat/internal/bot"
	"github.com/knhn1004/WhatDoIEat/internal/config"
	"github.com/knhn1004/WhatDoIEat/internal/models"
	"github.com/knhn1004/WhatDoIEat/internal/services"
)

const yelpAPIEndpoint = "https://api.yelp.com/v3/businesses/search"

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	botToken := config.BotToken
	appToken := config.AppToken

	bot.InitializeBot(botToken, appToken)
	services.StartServices()

	supabase := supa.CreateClient(config.SupabaseURL, config.SupabaseURL)
	fmt.Printf("supabase: %v\n", supabase) // TODO: remove this

	recipes, err := services.GenRecipeOpenAI()
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
	}

	/* var businesses []models.Restaurant
	options := map[string]string{
		"radius":  "10000", // 10km
		"sort_by": "rating",
		"limit":   "5",
	}
	businesses, err = getRestaurants("Japanese", "", options)
	if err != nil {
		fmt.Printf("getRestaurants error: %v\n", err)
	} else {
		for _, business := range businesses {
			fmt.Println("Name: ", business.Name)
			fmt.Println("Phone: ", business.Phone)
			fmt.Println("URL: ", business.URL)
			fmt.Println("Rating: ", business.Rating)
			fmt.Println("Location: ", business.Location)
			fmt.Println("Map: ", genGoogleMapsURL(business.Location))
			fmt.Println("Image: ", business.ImageURL)
		}
	} */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Start(ctx)
}

func getRestaurants(restaurantType, location string, options map[string]string) ([]models.Restaurant, error) {
	if location == "" {
		location = "San Jose"
	}

	params := url.Values{}
	params.Add("term", restaurantType)
	params.Add("location", location)

	for key, value := range options {
		params.Add(key, value)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", yelpAPIEndpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+config.YelpAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s. Response body: %s", resp.StatusCode, resp.Status, buf.String())
	}

	var yelpResp models.YelpResponse
	err = json.Unmarshal(buf.Bytes(), &yelpResp)
	if err != nil {
		return nil, err
	}

	return yelpResp.Businesses, nil
}

func genGoogleMapsURL(loc models.Location) string {
	baseURL := "https://www.google.com/maps/search/?api=1&query="

	// Concatenate address details
	address := loc.Address1 + " " + loc.Address2 + " " + loc.Address3 + " " + loc.City + " " + loc.ZipCode + " " + loc.State + " " + loc.Country

	// URL encode the address
	encodedAddress := url.QueryEscape(address)

	return baseURL + encodedAddress
}

func getImageUrlForRecipe(recipeName string, supabase *supa.Client) (string, error) {
	var img string
	serp := goduckgo.Search(goduckgo.Query{Keyword: recipeName})
	if len(serp.Results) > 0 {
		img = serp.Results[0].Image
	}
	if img == "" {
		return "", fmt.Errorf("no images found for recipe: %s", recipeName)
	}
	// download image data
	resp, err := http.Get(img)
	if err != nil {
		return "", fmt.Errorf("error downloading image: %v", err)
	}
	defer resp.Body.Close()

	var imgBuffer bytes.Buffer
	_, err = io.Copy(&imgBuffer, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading image data: %v", err)
	}

	// get file extension
	imgParts := strings.Split(img, ".")
	imgType := imgParts[len(imgParts)-1] // get last part, which should be the file extension

	// random filename
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), imgType)
	supabase.Storage.From("recipe-images").Upload(filename, &imgBuffer)

	supabaseUrl := os.Getenv("SUPABASE_URL")
	url := fmt.Sprintf("%s/storage/v1/object/public/recipe-images/%s", supabaseUrl, filename)

	return url, nil
}
