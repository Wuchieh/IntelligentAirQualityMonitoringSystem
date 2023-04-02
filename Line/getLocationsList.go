package Line

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/i18n"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"strings"
)

func getLocationsList(event *linebot.Event) {
	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error()))
		return
	}

	u := database.User{LineId: event.Source.UserID}
	if err := u.GetUserIdFromLineId(); err != nil {
		_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(notGetUseridFromLine(profile.Language))).Do()
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	l := database.Location{UserID: u.Id}
	locationList, err := l.GetLocationList()

	if err != nil {
		log.Println(err)
		return
	} else if len(locationList) == 0 {
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(locationsCount0(profile.Language))).Do()
		if err != nil {
			log.Println(err)
		}
		return
	}

	menu := createLocationListFlexContainer(profile.Language, locationList)
	flex := linebot.NewFlexMessage("LocationsList", menu)
	resp, err := bot.ReplyMessage(event.ReplyToken, flex).Do()
	if err != nil {
		log.Println(resp, err)
		return
	}
}

func createLocationListFlexContainer(language string, list []database.Location) *linebot.BubbleContainer {
	menu := createFlexTemplate()
	menu.Body.Contents = genLocationListBodyContents(language)
	menu.Footer.Contents = genLocationListFooterContents(list)
	return menu
}

func genLocationListFooterContents(locations []database.Location) []linebot.FlexComponent {
	var lFC []linebot.FlexComponent
	for i, location := range locations {
		locationInfo := &linebot.ButtonComponent{
			Type: linebot.FlexComponentTypeButton,
			Style: func(i int) linebot.FlexButtonStyleType {
				if i%2 == 0 {
					return linebot.FlexButtonStyleTypePrimary
				} else {
					return linebot.FlexButtonStyleTypeSecondary
				}
			}(i),
			Margin:       linebot.FlexComponentMarginTypeMd,
			OffsetBottom: linebot.FlexComponentOffsetTypeMd,
			Action: &linebot.PostbackAction{
				Data:  "location." + location.ID.String(),
				Label: location.NiceName,
			}}
		lFC = append(lFC, locationInfo)
	}
	return lFC
}

func genLocationListBodyContents(language string) []linebot.FlexComponent {
	var lFC []linebot.FlexComponent
	title := &linebot.TextComponent{
		Type:  linebot.FlexComponentTypeText,
		Align: linebot.FlexComponentAlignTypeCenter,
		Size:  linebot.FlexTextSizeTypeXxl,
	}
	title.Text = i18n.Get(language, "Locations")
	lFC = append(lFC, title)
	return lFC
}

func locationsCount0(language string) string {
	return i18n.Get(language, "locationsCount0")
}

func notGetUseridFromLine(language string) string {
	return strings.ReplaceAll(i18n.Get(language, "notGetUseridFromLine"), "{commandPrefix}", s.LineCommandPrefix)
}
