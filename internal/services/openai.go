package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/knhn1004/WhatDoIEat/internal/config"
	"github.com/knhn1004/WhatDoIEat/internal/models"
)

// OpenAI client
var openAIClient *openai.Client

// Initialize the OpenAI client
func InitOpenAI() {
	openAIClient = openai.NewClient(config.OpenAIKey)
}

func GenRecipeOpenAI() ([]models.Recipe, error) {
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
	resp, err := openAIClient.CreateChatCompletion(
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
	var recipes []models.Recipe
	err = json.Unmarshal([]byte(jsonStr), &recipes)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return recipes, nil
}
