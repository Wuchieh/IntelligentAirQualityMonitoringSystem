package Line

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

func setRichMenuLanguage(event *linebot.Event, lang string) {
	var err error
	switch lang {
	case "tw":
		_, err = bot.LinkUserRichMenu(event.Source.UserID, s.RichMenuIdTw).Do()
	default:
		_, err = bot.LinkUserRichMenu(event.Source.UserID, s.RichMenuIdEn).Do()
	}
	if err != nil {
		log.Println(err)
	}
}
