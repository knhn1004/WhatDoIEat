// internal/web/routes.go
package web

import (
	"github.com/gin-gonic/gin"

	"github.com/knhn1004/WhatDoIEat/internal/web/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/survey", controllers.SurveyHandler)
}
