package database

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	"time"
)

type User struct {
	Id          uuid.UUID
	Username    string
	Password    string
	Email       string
	CreateDate  time.Time
	UpdateDate  time.Time
	LineId      string
	DeleteAt    bool
	Admin       bool
	NoticeRange int
}

func (u *User) CreateNew() error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO users (username, password, createdate) VALUES($1, $2, now())")
	if err != nil {
		return pgError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Username, u.Password)
	if err != nil {
		return pgError(err)
	}

	err = tx.Commit()
	if err != nil {
		return pgError(err)
	}

	return nil
}

func (u *User) Verify(passwordHash string) (*uuid.UUID, bool, error) {
	stmt, err := db.Prepare("SELECT id,password FROM users WHERE username = $1 AND deleteat = false")
	if err != nil {
		return nil, false, pgError(err)
	}
	defer stmt.Close()

	var p string
	var id *uuid.UUID
	err = stmt.QueryRow(u.Username).Scan(&id, &p)
	if err != nil {
		return nil, false, pgError(err)
	}
	u.Id = *id
	return id, passwordHash == p, nil
}

func (u *User) GetLineID() string {
	stmt, err := db.Prepare("SELECT \"lineID\" FROM users WHERE id = $1")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer stmt.Close()

	var id sql.NullString
	err = stmt.QueryRow(u.Id).Scan(&id)
	if err != nil {
		log.Println(err)
		return ""
	}

	return id.String
}

func (u *User) GetNoticeRange() (int, error) {
	if u.NoticeRange > 0 {
		return u.NoticeRange, nil
	}
	stmt, err := db.Prepare("SELECT \"noticeRange\" FROM users WHERE id = $1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(u.Id).Scan(&u.NoticeRange)
	if err != nil {
		return 0, err
	}

	return u.NoticeRange, nil
}

func (u *User) SetNoticeRange(notificationRange int) error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE users SET \"noticeRange\" = $1, updatedate = now() WHERE id = $2")
	if err != nil {
		return pgError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(notificationRange, u.Id)
	if err != nil {
		return pgError(err)
	}

	err = tx.Commit()
	if err != nil {
		return pgError(err)
	}

	return nil
}

func (u *User) GetUserIdFromLineId() error {
	q, err := db.Prepare("SELECT id FROM users WHERE \"lineID\" = $1")
	if err != nil {
		return err
	}
	return pgError(q.QueryRow(u.LineId).Scan(&u.Id))
}

func (u *User) SetLineId(lineId string) error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()
	stmt, err := db.Prepare("UPDATE users SET \"lineID\" = $1, updatedate = now() WHERE id = $2;")
	if err != nil {
		return pgError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(lineId, u.Id)
	if err != nil {
		return pgError(err)
	}

	err = tx.Commit()
	if err != nil {
		return pgError(err)
	}
	return nil
}
