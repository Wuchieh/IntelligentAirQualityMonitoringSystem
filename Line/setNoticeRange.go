package Line

import (
	"fmt"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/i18n"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"strconv"
	"strings"
)

func setNoticeRange(event *linebot.Event, arg string) {
	num, err := strconv.Atoi(arg)
	if err != nil {
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(err)
		}
		return
	}

	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	if !isInNotificationRange(num) {
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(i18n.Get(profile.Language, "inputDataError"))).Do()
		if err != nil {
			log.Println(err)
		}
		return
	}

	var u = database.User{LineId: event.Source.UserID}
	if err = u.GetUserIdFromLineId(); err != nil {
		msg := strings.ReplaceAll(i18n.Get(profile.Language, "notGetUseridFromLine"), "{commandPrefix}", s.LineCommandPrefix)
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
			log.Println(err)
		}
		return
	}

	err = u.SetNoticeRange(num)
	if err != nil {
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	msg := fmt.Sprintf("%s %s ", i18n.Get(profile.Language, "setNoticeRange"),
		func() string {
			if num == 0 {
				return i18n.Get(profile.Language, "noNotification")
			} else if num == 1 {
				return fmt.Sprintf("%d %s", num, i18n.Get(profile.Language, "hour"))
			} else if num > 1 {
				return fmt.Sprintf("%d %s", num, i18n.Get(profile.Language, "hours"))
			}
			return ""
		}())
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do()
	if err != nil {
		log.Println(err)
	}
}
