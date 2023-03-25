package server

import (
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/ping", ping)
}

func ping(c *gin.Context) {
	c.String(200, "Pong!")
}
