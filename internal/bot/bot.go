package bot

import (
	"context"

	"github.com/shomali11/slacker/v2"
)

var bot *slacker.Slacker

func InitializeBot(botToken, appToken string) {
	bot = slacker.NewClient(botToken, appToken,
		slacker.WithBotMode(slacker.BotModeIgnoreApp))
	registerHandlers()
}

func registerHandlers() {
	// Register command defs
	bot.AddCommand(PingCmd)
	bot.AddCommand(UploadCmd)
	bot.AddCommand(EchoCmd)
}

func Start(ctx context.Context) error {
	return bot.Listen(ctx)
}
