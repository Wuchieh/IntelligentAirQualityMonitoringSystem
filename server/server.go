package server

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func Run(ip, port, mode string) error {
	setGinMode(mode)
	r := gin.Default()
	Router(r)
	return r.Run(ip + ":" + port)
}

func setGinMode(mode string) {
	if len(mode) < 1 {
		return
	}
	m := strings.ToLower(mode[:1])
	switch m {
	case "r":
		gin.SetMode(gin.ReleaseMode)
	case "t":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
