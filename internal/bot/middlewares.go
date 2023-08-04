package bot

import (
	"log"

	"github.com/shomali11/slacker/v2"
)

func logUserMiddleware() slacker.CommandMiddlewareHandler {
	return func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx *slacker.CommandContext) {
			// Log user info
			userId := ctx.Event().UserID
			profile := ctx.Event().UserProfile

			log.Printf("User %s (ID: %s) sent command", profile.DisplayName, userId)

			next(ctx)
		}
	}
}
