package server

import (
	"encoding/json"
	"errors"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	setting Setting
)

type Setting struct {
	IP           string `json:"IP"`
	PORT         string `json:"PORT"`
	RUNMODE      string `json:"RUN MODE"`
	JWTsecretKey string `json:"JWTsecretKey"`
}

func init() {
	if file, err := os.ReadFile("setting.json"); err != nil {
		log.Panicln("server os.ReadFile Error", err)
	} else {
		err = json.Unmarshal(file, &setting)
		if err != nil {
			log.Panicln("server json.Unmarshal Error", err)
		}
	}
}

func Run() error {
	if s := <-database.Sign; s != 0 {
		return errors.New("\033[31m 資料庫初始化錯誤 \033[0m")
	}

	setGinMode(setting.RUNMODE)
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	Router(r)
	return r.Run(setting.IP + ":" + setting.PORT)
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
