package server

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

var (
	notificationRange = [7]int{0, 1, 2, 4, 8, 12, 24}
)

type NotificationRangeReq struct {
	NotificationRange int `json:"notificationRange"`
}

func (r NotificationRangeReq) Check() bool {
	for _, i := range notificationRange {
		if r.NotificationRange == i {
			return true
		}
	}
	return false
}

func getSettingInfo(c *gin.Context) {
	id, lineID, ok := verifyUserLogic(c)
	if !ok {
		return
	}
	returnData := gin.H{"status": true}

	locations, err := getLocations(id)
	if err != nil {
		log.Println(err)
	} else {
		returnData["data"] = locations
	}

	NoticeInterval, err := getUserNoticeRange(id)
	if err != nil {
		log.Println(err)
		returnData["NoticeInterval"] = 0
	} else {
		returnData["NoticeInterval"] = NoticeInterval
	}

	if lineID != "" {
		returnData["LinkingLine"] = true
	}

	c.JSON(200, returnData)
}

func editNotificationRange(c *gin.Context) {
	id, _, ok := verifyUserLogic(c)
	if !ok {
		return
	}
	var nr NotificationRangeReq
	if err := c.Bind(&nr); err != nil || !nr.Check() {
		c.JSON(400, gin.H{"status": false, "msg": "inputDataError"})
		return
	}

	if uid, err := uuid.Parse(id); err != nil || !IsInNotificationRange(nr.NotificationRange) {
		c.JSON(400, gin.H{"status": false, "msg": "inputDataError"})
	} else {
		u := database.User{Id: uid}
		err = u.SetNoticeRange(nr.NotificationRange)
		if err != nil {
			c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": true, "msg": "editNotificationRangeSuccess"})
		}

	}
}

func IsInNotificationRange(i int) bool {
	for _, v := range notificationRange {
		if i == v {
			return true
		}
	}
	return false
}
