package controllers

import "github.com/gin-gonic/gin"

func SurveyHandler(c *gin.Context) {
	preferences := map[string]bool{
		"Spicy Food":       true,
		"Vegetarian Meals": false,
		"Dairy-Free":       false,
		"Seafood":          true,
		"Gluten-Free":      false,
	}

	c.HTML(200, "survey.html", gin.H{
		"Title":       "Survey Page",
		"Preferences": preferences,
	})
}
