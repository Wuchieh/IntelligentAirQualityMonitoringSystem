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

var (
	NotificationRange = [7]int{0, 1, 2, 4, 8, 12, 24}
)

func getSetNoticeRangeList(event *linebot.Event) {
	var u = database.User{LineId: event.Source.UserID}

	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		resp, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(resp, err)
			return
		}
		return
	}

	if err = u.GetUserIdFromLineId(); err != nil {
		msg := strings.ReplaceAll(i18n.Get(profile.Language, "notGetUseridFromLine"), "{commandPrefix}", s.LineCommandPrefix)
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
			log.Println(err)
		}
		return
	}

	menu := genSetNoticeRangeListMenu(profile.Language)
	flex := linebot.NewFlexMessage("LocationsList", menu)
	resp, err := bot.ReplyMessage(event.ReplyToken, flex).Do()

	if err != nil {
		log.Println(resp, err)
		return
	}
}

func genSetNoticeRangeListMenu(language string) *linebot.BubbleContainer {
	var menu = createFlexTemplate()
	menu.Body.Contents = genSetNoticeRangeListMenuBodyContent(language)
	menu.Footer.Contents = genSetNoticeRangeListMenuFooterContent(language)
	return menu
}

func genSetNoticeRangeListMenuFooterContent(language string) []linebot.FlexComponent {
	var menu []linebot.FlexComponent
	for i, v := range NotificationRange {
		button := &linebot.ButtonComponent{
			Type:   linebot.FlexComponentTypeButton,
			Margin: linebot.FlexComponentMarginTypeMd,
			Style: func() linebot.FlexButtonStyleType {
				if i%2 == 0 {
					return "primary"
				} else {
					return "secondary"
				}
			}(),
			Action: &linebot.PostbackAction{
				Data: "setNoticeRange." + strconv.Itoa(v),
				Label: func() string {
					if v == 0 {
						return i18n.Get(language, "noNotification")
					}
					return fmt.Sprintf("%d %s", v, hourFormat(language, v))
				}(),
			},
		}
		menu = append(menu, button)
	}
	return menu
}

func hourFormat(language string, i int) string {
	if i == 1 {
		return i18n.Get(language, "hour")
	} else if i > 1 {
		return i18n.Get(language, "hours")
	}
	return ""
}

func genSetNoticeRangeListMenuBodyContent(language string) []linebot.FlexComponent {
	return []linebot.FlexComponent{
		&linebot.TextComponent{
			Type:  linebot.FlexComponentTypeText,
			Align: linebot.FlexComponentAlignTypeCenter,
			Size:  linebot.FlexTextSizeTypeXxl,
			Text:  i18n.Get(language, "NoticeRange"),
		}}
}

func isInNotificationRange(i int) bool {
	for _, v := range NotificationRange {
		if i == v {
			return true
		}
	}
	return false
}
