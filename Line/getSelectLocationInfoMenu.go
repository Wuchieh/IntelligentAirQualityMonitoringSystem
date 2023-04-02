package Line

import (
	"fmt"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/i18n"
	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"strconv"
	"strings"
)

var (
	rangeRanges = []int{100, 200, 500}
)

func getSelectLocationInfoMenu(event *linebot.Event, locationId string) {
	lid, arg := func() (string, string) {
		a := strings.SplitN(locationId, ".", 2)
		if len(a) == 1 {
			return a[0], ""
		}
		return a[0], a[1]
	}()

	if arg != "" {
		setLocation(event, lid, arg)
		return
	}

	uid, err := uuid.Parse(lid)
	if err != nil {
		return
	}
	var l = database.Location{
		ID: uid,
	}

	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		resp, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(resp, err)
			return
		}
		return
	}

	if deleteAt, err := l.IsDelete(); err != nil {
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(err)
		}
		return
	} else if deleteAt {
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(i18n.Get(profile.Language, "editError_isDeleted"))).Do()
		if err != nil {
			log.Println(err)
		}
		return
	}

	if name, err := l.GetNickName(); err != nil {
		resp, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(err, resp)
		}
		return
	} else {
		menu := genSelectLocationInfoMenu(uid.String(), name, profile.Language)
		flex := linebot.NewFlexMessage("LocationsList", menu)
		resp, err := bot.ReplyMessage(event.ReplyToken, flex).Do()
		if err != nil {
			log.Println(resp, err)
			return
		}
	}
}

func setLocation(event *linebot.Event, lid string, arg string) {
	var l database.Location
	if uid, err := uuid.Parse(lid); err != nil {
		log.Println(err)
		return
	} else {
		l.ID = uid
		if userid, err := l.GetUserId(); err != nil {
			_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
			if err != nil {
				log.Println(err)
			}
			return
		} else {
			u := database.User{Id: *userid}
			if u.GetLineID() != event.Source.UserID {
				log.Println("身分認證錯誤")
				return
			}
		}
	}

	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		resp, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(resp, err)
			return
		}
		return
	}

	if deleteAt, err := l.IsDelete(); err != nil {
		log.Println(err)
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		if err != nil {
			log.Println(err)
		}
		return
	} else if deleteAt {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(i18n.Get(profile.Language, "editError_isDeleted"))).Do(); err != nil {
			log.Println(err)
		}
		return
	}

	if num, err := strconv.Atoi(arg); err != nil { //刪除
		if arg == "delete" {
			err := l.Delete()
			if err != nil {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do(); err != nil {
					log.Println(err)
				}
			} else {
				nickname, _ := l.GetNickName()
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(nickname+" "+i18n.Get(profile.Language, "isDeleted"))).Do(); err != nil {
					log.Println(err)
				}
			}
			return
		}
	} else { //設置Range
		if isInRange(num) {
			err := l.EditRange(num)
			if err != nil {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do(); err != nil {
					log.Println(err)
				}
				return
			} else {
				nickname, _ := l.GetNickName()
				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage(
						fmt.Sprintf("%s %s %d %s", nickname, i18n.Get(profile.Language, "isSelectRange"), num, i18n.Get(profile.Language, "meters")))).Do(); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func genSelectLocationInfoMenu(lid, locationNickName, language string) *linebot.BubbleContainer {
	menu := createFlexTemplate()
	menu.Body.Contents = genSelectLocationInfoBodyMenu(lid, locationNickName, language)
	menu.Body.Layout = linebot.FlexBoxLayoutTypeVertical
	menu.Footer.Contents = genSelectLocationInfoFooterMenu(lid, language)
	return menu
}

func genSelectLocationInfoFooterMenu(lid, language string) []linebot.FlexComponent {
	var menu []linebot.FlexComponent
	for _, i := range rangeRanges {
		button := &linebot.ButtonComponent{
			Type: linebot.FlexComponentTypeButton,
			Action: &linebot.PostbackAction{
				Data:  "location." + lid + "." + strconv.Itoa(i),
				Label: genMeter(language, i),
			},
			Margin: linebot.FlexComponentMarginTypeMd,
			Style: func(i int) linebot.FlexButtonStyleType {
				switch i {
				case 100:
					return "secondary"
				case 200:
					return "link"
				case 500:
					return "primary"
				}
				return "secondary"
			}(i),
		}
		menu = append(menu, button)
	}
	return menu
}

func genMeter(language string, i int) string {
	return strconv.Itoa(i) + " " + i18n.Get(language, "meters")
}

func genSelectLocationInfoBodyMenu(lid, locationNickName, language string) []linebot.FlexComponent {
	return []linebot.FlexComponent{
		&linebot.TextComponent{
			Type:  linebot.FlexComponentTypeText,
			Align: linebot.FlexComponentAlignTypeCenter,
			Size:  linebot.FlexTextSizeTypeXxl,
			Text:  locationNickName},
		&linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Margin: linebot.FlexComponentMarginTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.PostbackAction{
						Data:  "location." + lid + ".delete",
						Label: getDeleteLabel(language),
					},
				},
			},
		},
	}
}

func getDeleteLabel(language string) string {
	return i18n.Get(language, "delete")
}

func isInRange(num int) bool {
	for _, v := range rangeRanges {
		if v == num {
			return true
		}
	}
	return false
}
