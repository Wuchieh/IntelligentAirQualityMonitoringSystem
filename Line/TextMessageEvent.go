package Line

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"strings"
)

func TextMessageEvent(event *linebot.Event, msg *linebot.TextMessage) {

	if strings.HasPrefix(msg.Text, s.LineCommandPrefix) {
		command := strings.SplitN(msg.Text[len(s.LineCommandPrefix):], " ", 2)
		if len(command) >= 1 {
			switch strings.ToLower(command[0]) {
			case "login":
				lineLogin(event)
			case "setting":
				showSetting(event)
			}
		}
	} else if _, err := bot.ReplyMessage(event.ReplyToken,
		linebot.NewTextMessage(msg.Text)).Do(); err != nil {
		log.Println(err)
	}

}
