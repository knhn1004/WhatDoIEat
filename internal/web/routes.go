// internal/web/routes.go
package web

import (
	"github.com/gin-gonic/gin"

	"github.com/knhn1004/WhatDoIEat/internal/web/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/pref", controllers.SurveyHandler)
	r.POST("/save-pref", controllers.SavePrefHandler)
}
