package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var preferencesList = []string{
	"Spicy Food",
	"Vegetarian Meals",
	"Dairy-Free",
	"Seafood",
	"Gluten-Free",
}

func SurveyHandler(c *gin.Context) {
	pref := make(map[string]bool)
	for _, p := range preferencesList {
		pref[p] = false
	}

	userPref := []string{"Spicy Food", "Seafood"} // TODO: update with data from DB
	for _, p := range userPref {
		pref[p] = true
	}

	c.HTML(200, "survey.html", gin.H{
		"Title":       "Your Food Preferences",
		"Preferences": pref,
	})
}

func SavePrefHandler(c *gin.Context) {
	// get from form
	data := c.PostFormArray("preferences")

	// TODO: store to DB
	fmt.Println(data)
}
