package services

import (
	"fmt"
	"log"

	"github.com/cohere-ai/cohere-go"

	"github.com/knhn1004/WhatDoIEat/internal/config"
)

// Cohere client
var co *cohere.Client

// Initialize Cohere client
func InitCohere() {
	var err error
	co, err = cohere.CreateClient(config.CohereKey)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("co: %v\n", co)
	}
}

// Classify text using Cohere
func ClassifyText() {
	// Call Cohere API
	options := cohere.ClassifyOptions{
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
	}

	// Handle response
	res, err := co.Classify(options)

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
	}
}
