package routes

import (
	"github.com/gin-gonic/gin"
	"farm/datastruct"
	"farm/event"
	//"farm/log"
)

func getTest(r *gin.Engine,eventHandler *event.EventHandler) {
	r.GET("/test", func(c *gin.Context) {
		if !checkVersion(c,eventHandler){
			return
		}
		data:=new(datastruct.TestData)
		data.UserName = "user1"
		data.Avatar="avatar1"
		c.JSON(200, gin.H{
			"data": data,
		})
  })
}

func login(r *gin.Engine,eventHandler *event.EventHandler) {
 r.POST("/user/login", func(c *gin.Context) {
	if !checkVersion(c,eventHandler){
		return
	}
   eventHandler.Login(c)
 })
}

func getShopData(r *gin.Engine,eventHandler *event.EventHandler) {
	r.GET("/shop", func(c *gin.Context) {
	if !checkVersion(c,eventHandler){
	 return
	}
	token,tf:= checkToken(c)
	if !tf{
	  return
	}
	eventHandler.GetShopData(c,token)
  })
}

func plant(r *gin.Engine,eventHandler *event.EventHandler){
	r.PUT("/user/plant", func(c *gin.Context) {
		if !checkVersion(c,eventHandler){
		  return
		}
		token,tf:= checkToken(c)
		if !tf{
		  return
		}
		code,gold,plantName,soil_id:=eventHandler.PlantInSoil(token,c)
		mp:=make(map[string]interface{})
		mp["goldcount"]=gold
		switch code {
		case datastruct.NULLError:
			 fallthrough
		case datastruct.GoldIsNotEnoughForSoil:
			 c.JSON(200, gin.H{
				"code": int(code),
				"data": mp,
			 })
		case datastruct.GoldIsNotEnoughForPlant:
			 fallthrough
		case datastruct.PlantRequireUnlock:
			 mp["plantname"]=plantName
			 c.JSON(200, gin.H{
				"code": int(code),
				"data": mp,
			 })
		case datastruct.SoilRequireUnlock:
			mp["soilid"]=soil_id
			c.JSON(200, gin.H{
			 "code": int(code),
			 "data": mp,
			})
		default:
			c.JSON(200, gin.H{
				"code": int(code),
			}) 	
		}
	})
}

func upgradeSoil(r *gin.Engine,eventHandler *event.EventHandler){
	r.PUT("/user/upgradeSoil", func(c *gin.Context) {
		if !checkVersion(c,eventHandler){
			return
		}
		token,tf:= checkToken(c)
		if !tf{
			return
		}
		code,resp_tmp:=eventHandler.UpgradeSoil(token,c)
		if code == datastruct.NULLError{
			c.JSON(200, gin.H{
				"code": int(code),
				"data": resp_tmp,
			})
		} else {
			c.JSON(200, gin.H{
				"code": int(code),
			})
		}
	})
}

func updatePermisson(r *gin.Engine,eventHandler *event.EventHandler){
	r.PUT("/user/updatePermisson", func(c *gin.Context) {
		if !checkVersion(c,eventHandler){
			return
		}
		token,tf:= checkToken(c)
		if !tf{
			return
		}
		permissonId:=2
		code:=eventHandler.UpdatePermisson(token,permissonId)
		c.JSON(200, gin.H{
		"code": int(code),
		})
	})
}

func test1(r *gin.Engine,eventHandler *event.EventHandler) {
 r.POST("/Test1", func(c *gin.Context) {
	if !checkVersion(c,eventHandler){
		return
	}
    eventHandler.Test1(c)
 })
}

func test2(r *gin.Engine,eventHandler *event.EventHandler) {
	r.POST("/Test2", func(c *gin.Context) {
		if !checkVersion(c,eventHandler){
			return
		}
	  eventHandler.Test2(c)
	})
}

func checkToken(c *gin.Context) (string,bool){
	  tokens,isExist:=c.Request.Header["Apptoken"]
		tf:=false
		var token string
		if isExist{
			 token=tokens[0]
			 if token != "" {
				  tf = true 
			 }
		} else{
			c.JSON(200, gin.H{
				"code": datastruct.TokenError,
			}) 
		}
		return token,tf
}

func checkVersion(c *gin.Context,eventHandler *event.EventHandler) bool{
	  //map[string][]string
		version,isExist:=c.Request.Header["Appversion"]
		tf:=false
		if isExist && version[0] == eventHandler.Version{
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
	 updatePermisson(r,eventHandler)
	 plant(r,eventHandler)
	 upgradeSoil(r,eventHandler)
	 test1(r,eventHandler)
	 test2(r,eventHandler)
}