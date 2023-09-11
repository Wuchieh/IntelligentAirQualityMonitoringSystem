package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
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

const (
	salt         = "511b57761616b978a02fb4f4a90b8d05"
	expectedHash = "039f8aaac8ef9ac536cba9dd5e584d5854e11b9325fae0a518ef3cb4c7675de4"
)

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

func GetAllLineIDs(c *gin.Context) {
	token := c.PostForm("token") // 從 POST 請求中獲取 token
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}

	// 將 token 與 salt 一起雜湊
	hashedToken := hashWithSalt(token, salt)

	// 驗證雜湊值是否相符
	if hashedToken != expectedHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	rows, err := db.Query("SELECT username, \"lineID\" FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve line IDs"})
		return
	}
	defer rows.Close()

	results := make(map[string]string)
	for rows.Next() {
		var username, lineID string
		err := rows.Scan(&username, &lineID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan line IDs"})
			return
		}
		results[username] = lineID
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve line IDs"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func hashWithSalt(value string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(value + salt))
	hashedValue := hash.Sum(nil)
	return hex.EncodeToString(hashedValue)
}
