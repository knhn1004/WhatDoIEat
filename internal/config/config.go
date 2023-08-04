package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	// Slack
	BotToken string
	AppToken string

	// Supabase
	SupabaseURL string
	SupabaseKey string

	// Cohere
	CohereKey string

	// Yelp
	YelpAPIKey string

	// OpenAI
	OpenAIKey string
)

// Load reads config from .env file
func Load() error {
	// Load .env file
	if err := godotenv.Load(".env.local"); err != nil {
		return err
	}

	// Bot token
	if v, ok := os.LookupEnv("SLACK_BOT_TOKEN"); !ok {
		return errors.New("SLACK_BOT_TOKEN required")
	} else {
		BotToken = v
	}

	// App token
	if v, ok := os.LookupEnv("SLACK_APP_TOKEN"); !ok {
		return errors.New("SLACK_APP_TOKEN required")
	} else {
		AppToken = v
	}

	// Supabase URL
	if v, ok := os.LookupEnv("SUPABASE_URL"); !ok {
		return errors.New("SUPABASE_URL required")
	} else {
		SupabaseURL = v
	}

	// Supabase key
	if v, ok := os.LookupEnv("SUPABASE_ADMIN_KEY"); !ok {
		return errors.New("SUPABASE_ADMIN_KEY required")
	} else {
		SupabaseKey = v
	}

	// Cohere key
	if v, ok := os.LookupEnv("COHERE_API_KEY"); !ok {
		return errors.New("COHERE_API_KEY required")
	} else {
		CohereKey = v
	}

	// Yelp API key
	if v, ok := os.LookupEnv("YELP_API_KEY"); !ok {
		return errors.New("YELP_API_KEY required")
	} else {
		YelpAPIKey = v
	}

	// OpenAI API key
	if v, ok := os.LookupEnv("OPENAI_API_KEY"); !ok {
		return errors.New("OPENAI_API_KEY required")
	} else {
		OpenAIKey = v
	}

	return nil
}
