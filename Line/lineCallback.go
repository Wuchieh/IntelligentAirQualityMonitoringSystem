package Line

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"reflect"
)

func CallBack(c *gin.Context) {
	defer c.AbortWithStatus(200)
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.AbortWithStatus(400)
		} else {
			c.AbortWithStatus(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch msg := event.Message.(type) {
			case *linebot.TextMessage:
				TextMessageEvent(event, msg)
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is %s", msg.StickerID, msg.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			default:
				fmt.Println(reflect.TypeOf(msg))
			}
		} else if event.Type == linebot.EventTypePostback {
			postbackEven(event)
		}
	}
}
