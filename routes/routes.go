package routes

import (
	"github.com/gin-gonic/gin"
	"farm/datastruct"
	"farm/event"
)

func getTest(r *gin.Engine,eventHandler *event.EventHandler) {
	r.GET("/test", func(c *gin.Context) {
    if checkVersion(c){
			data:=new(datastruct.TestData)
			data.UserName = "user1"
			data.Avatar="avatar1"
			c.JSON(200, gin.H{
			"data": data,
			})
		}
  })
}

func login(r *gin.Engine,eventHandler *event.EventHandler) {
 r.POST("/login", func(c *gin.Context) {
   eventHandler.Login(c)
 })
}

func getShopData(r *gin.Engine,eventHandler *event.EventHandler) {
	r.GET("/shop", func(c *gin.Context) {
		eventHandler.GetShopData(c)
   })
}


func test1(r *gin.Engine,eventHandler *event.EventHandler) {
 r.POST("/Test1", func(c *gin.Context) {
   eventHandler.Test1(c)
 })
}

func test2(r *gin.Engine,eventHandler *event.EventHandler) {
	r.POST("/Test2", func(c *gin.Context) {
	  eventHandler.Test2(c)
	})
}

func checkVersion(c *gin.Context) bool{
	  //map[string][]string
		version,isExist:=c.Request.Header["Version"]
		tf:=false
		if isExist && version[0] == "1.0"{
			 tf = true
		} else {
			c.JSON(200, gin.H{
				"code": datastruct.VersionError,
			})
		}
		return tf
}

func Register(r *gin.Engine,eventHandler *event.EventHandler){
	 getTest(r,eventHandler)
	 getShopData(r,eventHandler)
	 login(r,eventHandler)
	 test1(r,eventHandler)
	 test2(r,eventHandler)
}