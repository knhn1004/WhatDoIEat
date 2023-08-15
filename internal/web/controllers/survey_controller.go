package controllers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/knhn1004/WhatDoIEat/internal/config"
	"github.com/knhn1004/WhatDoIEat/internal/services"
)

var preferencesList = []string{
	"Spicy Food",
	"Vegetarian Meals",
	"Dairy-Free",
	"Seafood",
	"Gluten-Free",
}

func SurveyHandler(c *gin.Context) {
	// get jwt token from query param
	tokenStr := c.Query("token")

	// validate
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppToken), nil
	})

	// Check for token validation errors and the token's expiration time
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var userId string
	// retrieve claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		userId = claims["userId"].(string)
	}
	fmt.Println(userId)

	pref := make(map[string]bool)
	for _, p := range preferencesList {
		pref[p] = false
	}

	// load user preferences from DB
	sp := services.GetSupabaseClient()

	var res interface{}
	err = sp.DB.From("users").Select("*").Eq("slack_id", userId).Execute(res)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

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
