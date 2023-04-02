package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	db   *sql.DB
	Sign = make(chan int, 1)

	host     string
	user     string
	password string
	dbname   string
)

type setting struct {
	DatabaseHost     string `json:"databaseHost"`
	DatabaseUser     string `json:"databaseUser"`
	DatabasePassword string `json:"databasePassword"`
	DatabaseDbname   string `json:"databaseDbname"`
}

func init() {
	var s setting

	if file, err := os.ReadFile("setting.json"); err != nil {
		panic(err)
	} else {
		err = json.Unmarshal(file, &s)
		if err != nil {
			panic(err)
		}
	}

	host = s.DatabaseHost
	user = s.DatabaseUser
	password = s.DatabasePassword
	dbname = s.DatabaseDbname
}

func DatabaseInit() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	//fmt.Println(dsn)
	if db, err = sql.Open("postgres", dsn); err != nil {
		log.Println(err)
		Sign <- 1
		return
	}
	if err = db.Ping(); err != nil {
		log.Println(err)
		Sign <- 1
		return
	}
	Sign <- 0
}
