package main

import (
	"github.com/gin-gonic/gin"
	"farm/datastruct"
	"farm/db"
)

var dbHandler *db.DBHandler

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

//添加版本号
func version()gin.HandlerFunc{
	return func(c *gin.Context) {
		c.Writer.Header().Add("Version", "1.0")
		c.Next()
	}
}

//跨域
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	dbHandler = db.CreateDBHandler()
	r := gin.Default()
	r.Use(cors())
	r.Use(version())
	getTest(r)
	r.Run("192.168.0.176:8080") // listen and serve on 0.0.0.0:8080
}

