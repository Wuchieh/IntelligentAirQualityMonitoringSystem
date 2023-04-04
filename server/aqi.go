package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
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

func aqiCreate(c *gin.Context) {
	var req requestApi
	err := c.Bind(&req)
	if err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		return
	} else if req.Aqi == 0 {
		c.JSON(400, gin.H{"status": false, "msg": errors.New("input data error").Error()})
		return
	}

	if req.Aqi >= 150.4 {
		go sendAlert(req)
		return
	}

	aqi := database.Aqi{
		Aqi: req.Aqi,
	}

	if err := aqi.Location.UnmarshalJSON([]byte(fmt.Sprintf("(%.50e,%.50e)", req.Lat, req.Lng))); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		return
	}

	if err := database.InsertAqi(aqi); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"status": false, "msg": "insertDataError"})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "success"})
}

func sendAlert(aqiReq requestApi) {
	marshal, err := json.Marshal(aqiReq)
	if err != nil {
		log.Println(err)
		return
	}
	req, err := http.NewRequest("POST", setting.AlertServerAddr+"/v1/danger", bytes.NewReader(marshal))
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Wuchieh/1.0")
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	all, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(all))
}
