package server

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func getAnnouncements(c *gin.Context) {
	// SELECT * FROM "public"."announcements" ORDER BY "id" DESC LIMIT 5;
	if a, err := database.GetAnnouncements(5); err != nil {
		log.Println(err)
		c.JSON(500, err)
	} else {
		c.JSON(200, a)
	}
}

func getAnnouncementContent(c *gin.Context) {
	atoi, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Println(err)
		return
	}
	var a database.Announcement
	err = a.GetContent(atoi)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, a.Content)
}
