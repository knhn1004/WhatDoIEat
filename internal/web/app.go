package web

import "github.com/gin-gonic/gin"

var r *gin.Engine

func InitWeb() {
	r = gin.Default()

	// Loading HTML templates
	r.LoadHTMLGlob("web/templates/*")

	// Setting up static file serving
	r.Static("/static", "web/static")

	// Routes
	SetupRoutes(r)

	go r.Run(":8080")
}
