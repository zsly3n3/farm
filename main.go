package main

import (
	"github.com/gin-gonic/gin"
	"farm/datastruct"
	"farm/event"
)


var eventHandler *event.EventHandler

func getTest(r *gin.Engine) {
	 data:=new(datastruct.TestData)
	 data.UserName = "user1"
	 data.Avatar="avatar1"
	 r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
		  "data": data,
		})
	})
}

func login(r *gin.Engine) {
  r.POST("/login", func(c *gin.Context) {
	eventHandler.Login(c)
  })
}


//跨域
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	eventHandler=event.CreateEventHandler()
	r := gin.Default()
	r.Use(cors())
	getTest(r)
	login(r)
	r.Run("192.168.0.161:8080")//listen and serve on 0.0.0.0:8080
}

