package main

import (
	"IntelligentAirQualityMonitoringSystem/server"
	"encoding/json"
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
	go func() {
		err := server.Run(setting.IP, setting.PORT, setting.RUNMODE)
		if err != nil {
			log.Panicln(err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
