package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter(requestHandler *RequestHandler) (*Router, error) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/test", requestHandler.Test)

	r.GET("/parse", requestHandler.ParseURL)

	return &Router{
		r,
	}, nil
}

func (r *Router) Serve(addr string) error {
	return r.Run(addr)
}
