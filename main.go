package main

import (
	"encoding/json"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	_ "github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/redis"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	setting Setting
)

type Setting struct {
	IP      string `json:"IP"`
	PORT    string `json:"PORT"`
	RUNMODE string `json:"RUN MODE"`
}

func init() {
	file, err := os.ReadFile("setting.json")
	if err != nil {
		log.Panicln(err)
	}
	if err = json.Unmarshal(file, &setting); err != nil {
		log.Panicln(err)
	}
}

func main() {
	database.DatabaseInit()
	sc := make(chan os.Signal, 1)

	go func() {
		defer func() { sc <- syscall.SIGINT }()
		_ = server.Run(setting.IP, setting.PORT, setting.RUNMODE)
	}()

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
