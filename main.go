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

	"github.com/cohere-ai/cohere-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/minoplhy/duckduckgo-images-api"
	supa "github.com/nedpals/supabase-go"
	openai "github.com/sashabaranov/go-openai"

	"github.com/knhn1004/WhatDoIEat/internal/bot"
)

var (
	botToken    string
	appToken    string
	supabaseUrl string
	supabaseKey string
	cohereKey   string
	yelpAPIKey  string
)

type Ingredient struct {
	Ingredient string `json:"ingredient"`
	Quantity   string `json:"quantity"`
}

type Recipe struct {
	Meal             string       `json:"meal"`
	Name             string       `json:"name"`
	ShortDescription string       `json:"short_description"`
	UserId           string       `json:"user_id"`
	ImageURL         string       `json:"image_url"`
	Date             string       `json:"date"`
	Ingredients      []Ingredient `json:"ingredients"`
	Steps            []string     `json:"steps"`
}

const yelpAPIEndpoint = "https://api.yelp.com/v3/businesses/search"

type YelpResponse struct {
	Businesses []Business `json:"businesses"`
}

type Business struct {
	Location Location `json:"location"`
	Name     string   `json:"name"`
	State    string   `json:"state"`
	Phone    string   `json:"phone"`
	URL      string   `json:"url"`
	ImageURL string   `json:"image_url"`
	Rating   float64  `json:"rating"`
}

type Location struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

type ErrorType struct {
	message string
}

func (e *ErrorType) Error() string {
	return e.message
}

func main() {
	loadEnv()

	bot.InitializeBot(botToken, appToken)

	supabase := supa.CreateClient(supabaseUrl, supabaseKey)
	fmt.Printf("supabase: %v\n", supabase) // TODO: remove this

	co, err := cohere.CreateClient(cohereKey)
	if err != nil {
		fmt.Printf("cohere error: %v\n", err)
	}
	fmt.Printf("cohere: %v\n", co) // TODO: remove this

	/* res, err := co.Classify(cohere.ClassifyOptions{
		Examples: []cohere.Example{
			{
				Text:  "I love you",
				Label: "positive",
			},
			{
				Text:  "I like you",
				Label: "positive",
			},
			{
				Text:  "I hate you",
				Label: "negative",
			},
			{
				Text:  "I don't like bananas",
				Label: "negative",
			},
		},
		Inputs: []string{
			"I like bananas",
			"I love milk",
		},
	})

	if err != nil {
		fmt.Printf("Classify error: %v\n", err)
	} else {
		for _, classification := range res.Classifications {
			if classification.Prediction == "positive" {
				fmt.Println("👍️")
			} else {
				fmt.Println("👎️")
			}
		}
	} */

	// genRecipeOpenAI()

	var businesses []Business
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
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Start(ctx)
}

func loadEnv() error {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	var errMsg string
	var ok bool

	botToken, ok = os.LookupEnv("SLACK_BOT_TOKEN")
	if !ok {
		errMsg = "SLACK_BOT_TOKEN is required"
		log.Fatal(errMsg)
		return &ErrorType{message: errMsg}
	}

	appToken, ok = os.LookupEnv("SLACK_APP_TOKEN")
	if !ok {
		errMsg = "SLACK_APP_TOKEN is required"
		log.Fatal(errMsg)
		return &ErrorType{message: errMsg}
	}

	supabaseUrl, ok = os.LookupEnv("SUPABASE_URL")
	if !ok {
		errMsg = "SUPABASE_URL is required"
		log.Fatal(errMsg)
		return &ErrorType{message: errMsg}
	}

	supabaseKey, ok = os.LookupEnv("SUPABASE_ADMIN_KEY")
	if !ok {
		errMsg = "SUPABASE_KEY is required"
		log.Fatal(errMsg)
		return &ErrorType{message: errMsg}
	}

	cohereKey, ok = os.LookupEnv("COHERE_API_KEY")
	if !ok {
		errMsg = "COHERE_API_KEY is required"
		log.Fatal(errMsg)
		return &ErrorType{message: errMsg}
	}

	yelpAPIKey, ok = os.LookupEnv("YELP_API_KEY")
	if !ok {
		errMsg = "YELP_API_KEY is required"
		log.Fatal(errMsg)
		return &ErrorType{message: errMsg}
	}

	return nil
}

func genRecipeOpenAI() ([]Recipe, error) {
	openaiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		log.Fatal("OPENAI_API_KEY is required")
	}
	client := openai.NewClient(openaiKey)

	sysPrompt := `
            assistant is an AI nutritionist
            given a user data, provide the next meals recipe a user can have today
            the output should be in the following json format
            ` + "````" + `
            [
              {
              "meal": "breakfast",
              "name": "",
              "short_description": "",
              "ingredients": [
                {"ingredient": "ingrdient name", "quantity": "quantity unit"},
                ...
              ],
              "steps": []
              },
              ...
            ]
            ` + "````"
	type Preferences struct {
		Styles    []string `json:"styles"`
		Diets     []string `json:"diets"`
		Allergies []string `json:"allergies"`
		Dislikes  []string `json:"dislikes"`
		Likes     []string `json:"likes"`
	}
	// preferences for food and diet
	pref := Preferences{
		Styles:    []string{"mediterranean", "vegan"},
		Diets:     []string{"keto", "paleo"},
		Allergies: []string{"dairy", "peanuts"},
		Dislikes:  []string{"broccoli", "carrots"},
		Likes:     []string{"chicken", "beef"},
	}
	prefJson, err := json.Marshal(pref)
	if err != nil {
		log.Fatal(err)
	}
	userPrompt := "user preferences: " + "```\n" + string(prefJson) + "```\n" + " Please only output the JSON, don't respond to me in any human language, just output nothing but JSON"
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: sysPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Temperature: 0,
		},
	)
	if err != nil {
		return nil, err
	}
	jsonStr := resp.Choices[0].Message.Content
	jsonStr = strings.Trim(jsonStr, "`")
	fmt.Println(jsonStr)
	var recipes []Recipe
	err = json.Unmarshal([]byte(jsonStr), &recipes)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return recipes, nil
}

func getRestaurants(restaurantType, location string, options map[string]string) ([]Business, error) {
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

	req.Header.Add("Authorization", "Bearer "+yelpAPIKey)

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

	var yelpResp YelpResponse
	err = json.Unmarshal(buf.Bytes(), &yelpResp)
	if err != nil {
		return nil, err
	}

	return yelpResp.Businesses, nil
}

func genGoogleMapsURL(loc Location) string {
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
