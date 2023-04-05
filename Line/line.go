package Line

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"os"
)

var (
	bot *linebot.Client
	s   setting
)

type setting struct {
	LineChannelSecret      string `json:"lineChannelSecret"`
	LineChannelAccessToken string `json:"lineChannelAccessToken"`
	LineCommandPrefix      string `json:"lineCommandPrefix"`
	RichMenuIdTw           string `json:"richMenuIdTw"`
	RichMenuIdEn           string `json:"richMenuIdEn"`
	ServerAddr             string `json:"ServerAddr"`
}

func init() {

	if file, err := os.ReadFile("setting.json"); err != nil {
		panic(err)
	} else {
		err = json.Unmarshal(file, &s)
		if err != nil {
			panic(err)
		}
	}

	if err := func() (err error) {
		bot, err = linebot.New(s.LineChannelSecret, s.LineChannelAccessToken)
		return
	}(); err != nil {
		panic(err)
	}
}

func createFlexTemplate() *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Header: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{},
		},
		Hero: &linebot.ImageComponent{
			Type:       linebot.FlexComponentTypeImage,
			URL:        s.ServerAddr + "/Logo/Logo.png",
			Gravity:    linebot.FlexComponentGravityTypeCenter,
			Size:       linebot.FlexImageSizeTypeMd,
			AspectMode: linebot.FlexImageAspectModeTypeFit,
		},
		Body: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeBaseline,
			Position: linebot.FlexComponentPositionTypeRelative,
			Spacing:  linebot.FlexComponentSpacingTypeMd,
			Margin:   linebot.FlexComponentMarginTypeMd,
			Contents: nil,
		},
		Footer: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: nil,
		},
		Styles: &linebot.BubbleStyle{
			Footer: &linebot.BlockStyle{
				Separator: true,
			},
		},
	}
}
