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
	"time"

	"github.com/cohere-ai/cohere-go"
	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
	openai "github.com/sashabaranov/go-openai"
	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
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
	Name    string  `json:"name"`
	Address string  `json:"address1"`
	City    string  `json:"city"`
	State   string  `json:"state"`
	ZipCode string  `json:"zip_code"`
	Rating  float64 `json:"rating"`
}

type ErrorType struct {
	message string
}

func (e *ErrorType) Error() string {
	return e.message
}

func main() {
	loadEnv()

	bot := slacker.NewClient(botToken, appToken,
		slacker.WithBotMode(slacker.BotModeIgnoreApp),
	)

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
		"limit":   "10",
	}
	businesses, err = getRestaurants("Japanese", "", options)
	if err != nil {
		fmt.Printf("getRestaurants error: %v\n", err)
	} else {
		fmt.Printf("businesses: %v\n", businesses)
	}

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx *slacker.CommandContext) {
			t1, _ := ctx.Response().Reply("about to be replaced 🚧️")
			fmt.Println(ctx.Event().UserID)
			fmt.Println(ctx.Event().UserProfile)

			time.Sleep(time.Second)

			ctx.Response().Reply("pong🏓️", slacker.WithReplace(t1))
		},
	})
	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "upload <sentence>",
		Description: "Upload a sentence!",
		Handler: func(ctx *slacker.CommandContext) {
			sentence := ctx.Request().Param("sentence")
			api := ctx.SlackClient()
			event := ctx.Event()

			api.PostMessage(event.ChannelID, slack.MsgOptionText("📄️ Uploading file ...", false))
			_, err := api.UploadFile(slack.FileUploadParameters{Content: sentence, Channels: []string{event.ChannelID}})
			if err != nil {
				ctx.Response().ReplyError(err)
			}
		},
	})
	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "echo {word}",
		Description: "Echo a word!",
		Handler: func(ctx *slacker.CommandContext) {
			word := ctx.Request().Param("word")

			attachments := []slack.Attachment{}
			attachments = append(attachments, slack.Attachment{
				Color:      "good",
				AuthorName: "👨️ Raed Shomali",
				Title:      "Attachment Title",
				Text:       "Attachment Text",
			})

			ctx.Response().Reply(word, slacker.WithAttachments(attachments))
		},
	})

	HelpDef := &slacker.CommandDefinition{
		Command:     "help",
		Description: "help!",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("Your own help function...")
		},
	}

	bot.Help(HelpDef)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
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
