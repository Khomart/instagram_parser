package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/khomart/instagram_recipe_parser/internal/config"
)

type Router struct {
	*gin.Engine
}

func NewRouter(c *config.Config, requestHandler *RequestHandler) (*Router, error) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/parse", AuthMiddleware(c.AllowedEmails), requestHandler.ParseURL)

	return &Router{
		r,
	}, nil
}

func (r *Router) Serve(addr string) error {
	return r.Run(addr)
}
