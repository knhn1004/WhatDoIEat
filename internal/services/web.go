package services

import (
	"github.com/gin-gonic/gin"

	"github.com/knhn1004/WhatDoIEat/internal/web"
)

var r *gin.Engine

func InitWeb() {
	r = gin.Default()

	// Loading HTML templates
	r.LoadHTMLGlob("web/templates/*")

	// Setting up static file serving
	r.Static("/static", "web/static")

	// Routes
	web.SetupRoutes(r)

	go r.Run(":8080")
}
