package routes

import (
	"github.com/gin-gonic/gin"
	"farm/datastruct"
	"farm/event"
)

func getTest(r *gin.Engine,eventHandler *event.EventHandler) {
	data:=new(datastruct.TestData)
	data.UserName = "user1"
	data.Avatar="avatar1"
	r.GET("/test", func(c *gin.Context) {
	   c.JSON(200, gin.H{
		 "data": data,
	   })
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

func Register(r *gin.Engine,eventHandler *event.EventHandler){
	 getTest(r,eventHandler)
	 getShopData(r,eventHandler)
	 login(r,eventHandler)
	 test1(r,eventHandler)
	 test2(r,eventHandler)
}