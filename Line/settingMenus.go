package Line

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/i18n"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

func showSetting(event *linebot.Event) {
	profile, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		log.Println(err)
		_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do()
		log.Println(err)
		return
	}
	menu := createFlexTemplate()
	menu.Body.Contents = getSettingBodyContents(profile.Language)
	menu.Footer.Contents = getSettingFooterContents(profile.Language)

	settingMenu := linebot.NewFlexMessage("setting", menu)
	_, err = bot.ReplyMessage(event.ReplyToken, settingMenu).Do()
	if err != nil {
		log.Println(err)
		return
	}
}

func getSettingFooterContents(language string) []linebot.FlexComponent {
	var lFC []linebot.FlexComponent
	var locationsButton = &linebot.ButtonComponent{
		Type:         linebot.FlexComponentTypeButton,
		Style:        linebot.FlexButtonStyleTypeSecondary,
		Margin:       linebot.FlexComponentMarginTypeMd,
		OffsetBottom: linebot.FlexComponentOffsetTypeMd,
	}
	var SetNoticeRangeButton = &linebot.ButtonComponent{
		Type:         linebot.FlexComponentTypeButton,
		Style:        linebot.FlexButtonStyleTypePrimary,
		Margin:       linebot.FlexComponentMarginTypeMd,
		OffsetBottom: linebot.FlexComponentOffsetTypeMd,
	}

	locationsButton.Action = &linebot.PostbackAction{
		Data:  "Locations",
		Label: i18n.Get(language, "Locations"),
	}
	SetNoticeRangeButton.Action = &linebot.PostbackAction{
		Data:  "SetNoticeRange",
		Label: i18n.Get(language, "NoticeRange"),
	}

	lFC = append(append(lFC, locationsButton), SetNoticeRangeButton)
	return lFC
}

func getSettingBodyContents(language string) []linebot.FlexComponent {
	var lFC []linebot.FlexComponent
	var bodyContent = &linebot.TextComponent{
		Type:  linebot.FlexComponentTypeText,
		Align: linebot.FlexComponentAlignTypeCenter,
		Size:  linebot.FlexTextSizeType3xl,
	}

	bodyContent.Text = i18n.Get(language, "Setting")

	lFC = append(lFC, bodyContent)
	return lFC
}
