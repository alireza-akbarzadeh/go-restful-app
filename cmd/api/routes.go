package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	g := gin.Default()

	// Health check endpoint
	g.GET("/api/hello", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	v1 := g.Group("/api/v1")
	{
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getAllEvent)
		v1.GET("/events/:id", app.getEvent)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events/:id", app.deleteEvent)

		v1.POST("/auth/register", app.registerUser)
	}

	return g
}
