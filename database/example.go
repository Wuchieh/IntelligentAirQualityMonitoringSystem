package database

import (
	"database/sql"
	"fmt"
	"time"
)

func e() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	db, err := sql.Open("postgres", "host=127.0.0.1 user=wuchieh password=Fa26701993 dbname=wuchieh sslmode=disable")
	checkErr(err)

	//插入資料
	stmt, err := db.Prepare("INSERT INTO userinfo(username,department,created) VALUES($1,$2,$3) RETURNING uid")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研發部門", "2012-12-09")
	checkErr(err)

	//pg 不支援這個函式，因為他沒有類似 MySQL 的自增 ID
	// id, err := res.LastInsertId()
	// checkErr(err)
	// fmt.Println(id)

	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研發部門", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("最後插入 id =", lastInsertId)

	//更新資料
	stmt, err = db.Prepare("update userinfo set username=$1 where uid=$2")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", 1)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查詢資料
	sqlStatement := "SELECT * FROM \"public\".\"users\" LIMIT 1000 OFFSET 0;"
	rows, err := db.Query(sqlStatement)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		var (
			id         int
			username   string
			password   string
			email      string
			createDate time.Time
			updateDate *time.Time
			deleteAt   bool
		)

		switch err := rows.Scan(&id,
			&username,
			&password,
			&email,
			&createDate,
			&updateDate,
			&deleteAt); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
		case nil:
			fmt.Println(id,
				username,
				password,
				email,
				createDate,
				updateDate,
				deleteAt)
		default:
			checkErr(err)
		}
	}

	//刪除資料
	stmt, err = db.Prepare("delete from userinfo where uid=$1")
	checkErr(err)

	res, err = stmt.Exec(1)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()
}
