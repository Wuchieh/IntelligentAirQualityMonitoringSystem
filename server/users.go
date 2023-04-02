package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type resUser struct {
	Username   string `json:"Username"`
	Password   string `json:"Password"`
	RememberMe bool   `json:"RememberMe"`
}

func userRegister(c *gin.Context) {
	var u resUser

	if err := c.Bind(&u); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": "inputDataError"})
		return
	}

	switch {
	case len(u.Username) < 8 || len(u.Password) < 8:
		c.JSON(400, gin.H{"status": false, "msg": "inputLengthLt8"})
	default:
		err := userRegisterLogic(u.Username, u.Password)
		if err != nil {
			c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": true, "msg": "registerSuccess"})
		}
	}
}

func userLogin(c *gin.Context) {
	var u resUser

	if err := c.Bind(&u); err != nil {
		c.JSON(400, gin.H{"status": false, "msg": "inputDataError"})
		return
	}

	switch {
	case len(u.Username) < 8 || len(u.Password) < 8:
		c.JSON(400, gin.H{"status": false, "msg": "inputLengthLt8"})
	default:
		var err = userLoginLogic(c, u)
		if err != nil {
			c.JSON(400, gin.H{"status": false, "msg": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": true, "msg": "loginSuccess"})
		}
	}
}

func verifyUser(c *gin.Context) {
	if _, _, ok := verifyUserLogic(c); ok {
		c.JSON(200, nil)
	}
}

func userLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"msg": "Logout success",
	})
}

func userLineLogin(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.Redirect(302, "/")
		return
	}
	userLineLoginLogic(c, userId)
}
