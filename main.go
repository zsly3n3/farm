package main

import (
	"github.com/gin-gonic/gin"
	"farm/event"
	"farm/routes"
)


var eventHandler *event.EventHandler

//跨域
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		// c.Next()

		method := c.Request.Method
 
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
 
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
		   c.AbortWithStatus(204)
		}
		// 处理请求
		c.Next()
	}
}

func main() {
	eventHandler=event.CreateEventHandler()
	r := gin.Default()
	r.Use(cors())
	routes.Register(r,eventHandler)
	r.Run("192.168.0.161:8080")//listen and serve on 0.0.0.0:8080
}

