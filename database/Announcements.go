package database

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	//AnnouncementCache = make(map[int][]string)

	AnnouncementList []Announcement
)

type Announcement struct {
	Id         int
	Title      []string
	Content    []string
	CreateTime time.Time
	Hidden     bool
}

type tAnnouncement struct {
	Id         int
	Title      *[]uint8
	Content    *[]string
	CreateTime *time.Time
	Hidden     *bool
}

/*
CREATE TABLE "public"."announcements" (
  "id" int4 NOT NULL DEFAULT nextval('"Announcements_id_seq"'::regclass),
  "title" varchar(50)[] COLLATE "pg_catalog"."default",
  "content" text[] COLLATE "pg_catalog"."default",
  "createTime" timestamp(6),
  "hidden" bool DEFAULT false,
  CONSTRAINT "Announcements_pkey" PRIMARY KEY ("id")
)
*/

func GetAnnouncements(count int) ([]Announcement, error) {
	if len(AnnouncementList) > 1 {
		return AnnouncementList, nil
	}
	query, err := db.Query("SELECT id,title,\"createTime\" FROM announcements WHERE hidden=false ORDER BY \"createTime\" DESC LIMIT $1", count)
	if err != nil {
		return nil, err
	}
	var as []Announcement
	for query.Next() {
		var ta tAnnouncement
		if err = query.Scan(
			&ta.Id,
			&ta.Title,
			&ta.CreateTime); err != nil {
			break
		}
		var title []string
		if ta.Title != nil {
			a := string(*ta.Title)
			title = strings.Split(a[1:len(a)-1], ",")
		}
		a := Announcement{
			Id:         ta.Id,
			Title:      title,
			CreateTime: *ta.CreateTime,
		}
		as = append(as, a)
	}
	AnnouncementList = as
	return as, err
}

func (a *Announcement) GetContent(id int) error {
	result, err := redis.Redis.LRange(redis.CTX, "AnnouncementCache_"+strconv.Itoa(id), 0, -1).Result()
	if err != nil {
		log.Println(err)
	} else if len(result) > 0 {
		a.Content = result
		return nil
	}

	//if c, ok := AnnouncementCache[id]; ok {
	//	a.Content = c
	//	return nil
	//}

	query, err := db.Query("SELECT unnest(content) FROM announcements WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return err
	}

	var contentUint8 []uint8
	a.Content = []string{}
	for query.Next() {
		err := query.Scan(&contentUint8)
		if err != nil {
			return err
		}
		a.Content = append(a.Content, string(contentUint8))
	}

	redis.Redis.RPush(redis.CTX, "AnnouncementCache_"+strconv.Itoa(id), a.Content)
	//AnnouncementCache[id] = a.Content
	return err
}
