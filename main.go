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
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
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

