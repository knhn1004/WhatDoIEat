// internal/web/controllers/survey_controller.go
package controllers

import "github.com/gin-gonic/gin"

func SurveyHandler(c *gin.Context) {
	c.HTML(200, "survey.html", gin.H{
		"Title":   "Survey Page",
		"Content": "Some content you want to render...",
	})
}
