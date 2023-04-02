package server

import (
	"github.com/Wuchieh/IntelligentAirQualityMonitoringSystem/Line"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/ping", ping)
	r.GET("/test", test)

	api := r.Group("/api")

	users := api.Group("/users")
	users.GET("/login", userLineLogin)
	users.POST("/register", userRegister)
	users.POST("/login", userLogin)
	users.POST("/logout", userLogout)
	users.POST("/verify", verifyUser)

	announcements := api.Group("/announcements")
	announcements.GET("", getAnnouncements)
	announcements.GET("/content", getAnnouncementContent)

	aqi := api.Group("/aqi")
	aqi.GET("/search", aqiSearch)
	aqi.POST("/create", aqiCreateNew)

	location := api.Group("/location")
	location.POST("/create", createLocation)
	location.POST("/edit", editLocation)
	location.POST("/delete", deleteLocation)

	setting := api.Group("/setting")
	setting.GET("/info", getSettingInfo)
	setting.POST("/edit", editNotificationRange)

	line := api.Group("/line")
	line.POST("/callback", Line.CallBack)

	test1 := api.Group("/test")
	test1.GET("/", test)

}

func test(c *gin.Context) {
	c.Param("")
	c.JSON(200, gin.H{})
}

func ping(c *gin.Context) {
	c.String(200, "Pong!")
}
