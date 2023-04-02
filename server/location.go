package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

type rqLocation struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Name string  `json:"name"`
}

type editLocationReq struct {
	ID    uuid.UUID `json:"ID"`
	Range int       `json:"Range"`
}

func createLocation(c *gin.Context) {
	id, _, ok := verifyUserLogic(c)
	if !ok {
		return
	}

	var rl rqLocation
	if err := c.Bind(&rl); err != nil || rl.Name == "" {
		c.JSON(400, gin.H{"status": false, "msg": "inputDataError"})
		return
	}

	if err := createLocationLogic(rl, id); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"status": true, "msg": "locationCreateSuccess"})
	}
}

func editLocation(c *gin.Context) {
	id, _, ok := verifyUserLogic(c)
	if !ok {
		return
	}
	if err := editLocationLogic(c, id); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"status": true, "msg": "editLocationSuccess"})
	}
}

func deleteLocation(c *gin.Context) {
	id, _, ok := verifyUserLogic(c)
	if !ok {
		return
	}
	var l editLocationReq
	if err := c.Bind(&l); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": "inputDataError"})
		return
	}
	err := deleteLocationLogic(id, l)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"status": true, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "deleteLocationSuccess"})
}
