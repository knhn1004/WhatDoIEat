package bot

import (
	"github.com/shomali11/slacker/v2"
)

// Send message to channel
func SendMessage(channel, text string) {
	// Use bot client to send message
}

// Reply to command context
func Reply(ctx *slacker.CommandContext, text string) {
	// Use ctx to send reply
}

// Reply with error
func ReplyError(ctx *slacker.CommandContext, err error) {
	// Send error reply
}
