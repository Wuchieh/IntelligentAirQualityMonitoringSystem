package main

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	_ "github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/redis"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	database.DatabaseInit()
	sc := make(chan os.Signal, 1)

	go func() {
		defer func() { sc <- syscall.SIGINT }()
		err := server.Run()
		log.Println(err)
	}()

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
