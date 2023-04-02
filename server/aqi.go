package server

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/gin-gonic/gin"
	"log"
)

func aqiSearch(c *gin.Context) {
	s := database.AqiSearch{Lat: -1}
	aqis, err := database.GetAqi(s)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"status": false, "msg": "getDataError"})
	} else {
		c.JSON(200, aqis)
	}
}
