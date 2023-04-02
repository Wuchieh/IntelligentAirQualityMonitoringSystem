package server

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/database"
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func passwordHash(u, p string) string {
	h := sha256.New()
	h.Write([]byte(u + p))
	return hex.EncodeToString(h.Sum(nil))
}

func userRegisterLogic(username, password string) error {
	var u database.User
	u.Username = username
	u.Password = passwordHash(username, password)
	err := u.CreateNew()
	if err != nil {
		return err
	}
	return nil
}

func userLoginLogic(c *gin.Context, u resUser) error {
	user := database.User{Username: u.Username, Password: u.Password}
	if id, ok, err := user.Verify(passwordHash(u.Username, u.Password)); err != nil || !ok {
		log.Println(err)
		return errors.New("verifyError")
	} else {
		t := userRemember(u.RememberMe)
		if err := linkLine(c, user); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return err
		}
		token, err := generateToken(*id, t)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return err
		}

		c.SetCookie("token", token, 60*60*t, "", "", false, true)
		return nil
	}
}

func linkLine(c *gin.Context, user database.User) error {
	userUUID := c.Query("userId")
	result, err := redis.Redis.Get(redis.CTX, userUUID).Result()
	if err != nil {
		return nil
	}
	if result != "" {
		if err := user.SetLineId(result); err != nil {
			return err
		} else {
			redis.Redis.Del(redis.CTX, userUUID, result)
		}
	}
	return nil
}

func userRemember(b bool) int {
	if b {
		return 24 * 7
	} else {
		return 1
	}
}

func verifyUserLogic(c *gin.Context) (string, string, bool) {
	token, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "idVerifyError",
			})
			return "", "", false
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Bad request",
		})
		return "", "", false
	}

	id, lineID, err := verifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": err.Error(),
		})
		return "", "", false
	}
	return id, lineID, true
}

func getUserNoticeRange(id string) (*int, error) {
	if uid, err := uuid.Parse(id); err != nil {
		return nil, err
	} else {
		u := database.User{Id: uid}
		noticeRange, err := u.GetNoticeRange()
		if err != nil {
			return nil, err
		} else {
			return &noticeRange, nil
		}
	}
}

func userLineLoginLogic(c *gin.Context, id string) {
	result, err := redis.Redis.Get(redis.CTX, id).Result()
	if err != nil {
		c.String(401, "Invalid userId")
		return
	}
	//c.String(200, result)

	u := database.User{LineId: result}
	if err := u.GetUserIdFromLineId(); err != nil { // 找不到已綁定lineID的user 跳轉至login.html進行綁定
		c.Redirect(302, "/login.html?userId="+id)
		return
	}

	token, err := generateToken(u.Id, 24*7)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("token", token, 60*60*24*7, "", "", false, true)
	c.Redirect(302, "/")
	redis.Redis.Del(redis.CTX, id, result)
}
