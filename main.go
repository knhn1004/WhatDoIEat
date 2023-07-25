package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("SLACK_BOT_TOKEN is required")
	}
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		log.Fatal("SLACK_APP_TOKEN")
	}

	bot := slacker.NewClient(botToken, appToken,
		slacker.WithBotMode(slacker.BotModeIgnoreApp),
	)
	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx *slacker.CommandContext) {
			t1, _ := ctx.Response().Reply("about to be replaced")

			time.Sleep(time.Second)

			ctx.Response().Reply("pong", slacker.WithReplace(t1))
		},
	})
	bot.AddCommand(&slacker.CommandDefinition{
		Command:     "upload <sentence>",
		Description: "Upload a sentence!",
		Handler: func(ctx *slacker.CommandContext) {
			sentence := ctx.Request().Param("sentence")
			slackClient := ctx.SlackClient()
			event := ctx.Event()

			slackClient.PostMessage(event.ChannelID, slack.MsgOptionText("Uploading file ...", false))
			_, err := slackClient.UploadFile(slack.FileUploadParameters{Content: sentence, Channels: []string{event.ChannelID}})
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
				AuthorName: "Raed Shomali",
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
