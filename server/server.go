package server

import (
	"errors"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/gin-gonic/gin"
	"strings"
)

func Run(ip, port, mode string) error {
	if s := <-database.Sign; s != 0 {
		return errors.New("\033[31m 資料庫初始化錯誤 \033[0m")
	}

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
