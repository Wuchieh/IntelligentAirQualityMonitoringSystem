package Line

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"strings"
)

func postbackEven(event *linebot.Event) {
	switch event.Postback.Data {
	case "Locations":
		getLocationsList(event)
	case "SetNoticeRange":
		getSetNoticeRangeList(event)
	default:
		switch {
		case strings.HasPrefix(event.Postback.Data, "location."):
			getSelectLocationInfoMenu(event, event.Postback.Data[len("location."):])
		case strings.HasPrefix(event.Postback.Data, "setNoticeRange."):
			setNoticeRange(event, event.Postback.Data[len("setNoticeRange."):])
		default:
			log.Println("Warring:", event.Postback)
		}
	}
}
