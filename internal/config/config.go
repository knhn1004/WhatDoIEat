package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Target  *string
	KeyName string
}

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

	configs = []Configuration{
		{&BotToken, "SLACK_BOT_TOKEN"},
		{&AppToken, "SLACK_APP_TOKEN"},
		{&SupabaseURL, "SUPABASE_URL"},
		{&SupabaseKey, "SUPABASE_ADMIN_KEY"},
		{&CohereKey, "COHERE_API_KEY"},
		{&YelpAPIKey, "YELP_API_KEY"},
		{&OpenAIKey, "OPENAI_API_KEY"},
	}
)

// Load reads config from .env file
func Load() error {
	// Load .env file
	if err := godotenv.Load(".env.local"); err != nil {
		return err
	}

	for _, config := range configs {
		if value, ok := os.LookupEnv(config.KeyName); !ok {
			return errors.New(config.KeyName + " required")
		} else {
			*config.Target = value
		}
	}

	return nil
}
