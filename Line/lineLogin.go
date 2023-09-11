package Line

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/i18n"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/redis"
	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"time"
)

func lineLogin(event *linebot.Event) {
	targetID := event.Source.GroupID
	if targetID == "" {
		targetID = event.Source.UserID
	}

	loginUUID, err := getUserLoginUUID(targetID)

	//fmt.Println("targetID:" + targetID)
	//fmt.Println("UserID:" + event.Source.UserID)
	//fmt.Println("GroupID:" + event.Source.GroupID)
	//fmt.Println("RoomID:" + event.Source.RoomID)
	if err != nil {
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do(); err != nil {
			log.Println(err)
		}
		return
	}

	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		return
	}

	var ButtonsTemplateTitle, ButtonsTemplateText, ButtonsTemplateLabel string

	ButtonsTemplateTitle = i18n.Get(profile.Language, "ButtonsTemplateTitle")
	ButtonsTemplateText = i18n.Get(profile.Language, "ButtonsTemplateText")
	ButtonsTemplateLabel = i18n.Get(profile.Language, "ButtonsTemplateLabel")

	uriAction := linebot.NewURIAction(ButtonsTemplateLabel, s.ServerAddr+"/api/users/login?userId="+loginUUID)

	buttonTemplate := linebot.NewButtonsTemplate(
		s.ServerAddr+"/Logo/Ico/LogoWhile.ico",
		ButtonsTemplateTitle,
		ButtonsTemplateText,
		uriAction,
	)

	templateMessage := linebot.NewTemplateMessage("My Menu", buttonTemplate)

	if _, err := bot.ReplyMessage(event.ReplyToken, templateMessage).Do(); err != nil {
		log.Print(err)
		return
	}
}

func getUserLoginUUID(userid string) (string, error) {
	result, err := redis.Redis.Get(redis.CTX, userid).Result()
	if err != nil {
		log.Println(err)
	} else {
		return result, nil
	}

	uid := uuid.New()
	set := redis.Redis.Set(redis.CTX, userid, uid.String(), 5*time.Minute)
	set2 := redis.Redis.Set(redis.CTX, uid.String(), userid, 5*time.Minute)
	if set.Err() != nil || set2.Err() != nil {
		redis.Redis.Del(redis.CTX, userid, uid.String())
		return "", set.Err()
	}
	return uid.String(), nil
}
