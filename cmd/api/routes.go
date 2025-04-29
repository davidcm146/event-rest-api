package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)
		v1.GET("/events/:id/attendees", app.getAttendeesByEvent)
		v1.GET("/attendees/:id/events", app.getEventsByAttendee)

		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.loginUser)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.authMiddleware())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)
		authGroup.DELETE("/events/:id/attendees/:userId", app.removeAttendeeFromEvent)

	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return g
}
