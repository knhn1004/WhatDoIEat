package bot

import (
	"time"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack"
)

// Command definitions
var (
	PingCmd = &slacker.CommandDefinition{
		Command: "ping",
		Handler: handlePing,
		Middlewares: []slacker.CommandMiddlewareHandler{
			logUserMiddleware(),
		},
	}

	UploadCmd = &slacker.CommandDefinition{
		Command: "upload",
		Handler: handleUpload,
	}

	EchoCmd = &slacker.CommandDefinition{
		Command: "echo",
		Handler: handleEcho,
	}
)

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
