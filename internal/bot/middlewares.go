package bot

import (
	"log"

	"github.com/shomali11/slacker/v2"

	"github.com/knhn1004/WhatDoIEat/internal/services"
)

func logUserMiddleware() slacker.CommandMiddlewareHandler {
	return func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx *slacker.CommandContext) {
			// Log user info
			userId := ctx.Event().UserID
			profile := ctx.Event().UserProfile

			log.Printf("User %s (ID: %s) sent command", profile.DisplayName, userId)

			sp := services.GetSupabaseClient()

			// Check if user exists in DB
			var userResults []map[string]interface{}
			err := sp.DB.From("users").Select("slack_id").Eq("slack_id", userId).Execute(&userResults)
			if err != nil {
				log.Printf("Error checking user in DB: %v", err)
			}

			// If user doesn't exist in DB, add them
			if len(userResults) == 0 {
				row := map[string]interface{}{
					"slack_id":  userId,
					"food_pref": []string{},
				}
				var insertResults []map[string]interface{}
				err := sp.DB.From("users").Insert(row).Execute(&insertResults)
				if err != nil {
					log.Printf("Error adding user to DB: %v", err)
					return
				}
				log.Printf("Added user %s to DB", profile.DisplayName)
			}

			next(ctx)
		}
	}
}
