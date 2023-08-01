package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

var (
	botToken    string
	appToken    string
	supabaseUrl string
	supabaseKey string
)

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
	fmt.Printf("supabase: %v\n", supabase)

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx *slacker.CommandContext) {
			t1, _ := ctx.Response().Reply("about to be replaced 🚧️")

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

	err := bot.Listen(ctx)
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

	botToken = os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		msg := "SLACK_BOT_TOKEN is required"
		log.Fatal(msg)
		return &ErrorType{message: msg}
	}

	appToken = os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		msg := "SLACK_APP_TOKEN is required"
		log.Fatal(msg)
		return &ErrorType{message: msg}
	}

	supabaseUrl = os.Getenv("SUPABASE_URL")
	if supabaseUrl == "" {
		msg := "SUPABASE_URL is required"
		log.Fatal(msg)
		return &ErrorType{message: msg}
	}

	supabaseKey = os.Getenv("SUPABASE_ADMIN_KEY")
	if supabaseKey == "" {
		msg := "SUPABASE_KEY is required"
		log.Fatal(msg)
		return &ErrorType{message: msg}
	}

	return nil
}
