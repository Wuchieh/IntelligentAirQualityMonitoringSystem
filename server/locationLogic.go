package server

import (
	"errors"
	"fmt"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func createLocationLogic(rl rqLocation, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	location := database.Location{UserID: uid}
	if ls, err := location.GetLocationList(); err != nil {
		return err
	} else if len(ls) >= 3 {
		return errors.New("locationsMaximum")
	}

	newLocation := database.Location{
		NiceName: rl.Name,
		UserID:   uuid.MustParse(id),
	}

	if err := newLocation.Location.UnmarshalJSON([]byte(fmt.Sprintf("(%.50e,%.50e)", rl.Lat, rl.Lng))); err != nil {
		return err
	}

	if err := newLocation.Create(); err != nil {
		return err
	}

	return nil
}

func getLocations(id string) ([]database.Location, error) {
	l := database.Location{UserID: uuid.MustParse(id)}
	return l.GetLocationList()
}

func editLocationLogic(c *gin.Context, id string) error {
	var el editLocationReq
	if err := c.Bind(&el); err != nil {
		return errors.New("inputDataError")
	}
	l := database.Location{ID: el.ID}
	if userId, err := l.GetUserId(); err != nil {
		return err
	} else if userId.String() != id {
		return errors.New("idVerifyError")
	}
	if err := l.EditRange(el.Range); err != nil {
		return err
	}
	return nil
}

func deleteLocationLogic(id string, l editLocationReq) error {
	la := database.Location{ID: l.ID}
	if userId, err := la.GetUserId(); err != nil {
		return err
	} else {
		if userId.String() != id {
			return errors.New("idVerifyError")
		}
	}
	return la.Delete()
}
