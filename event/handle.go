package event

import(
	"time"
	"farm/datastruct"
	"github.com/gin-gonic/gin"
	//"farm/log"
)

func (handle *EventHandler)Login(c *gin.Context){
	var body datastruct.UserLogin
	err:=c.BindJSON(&body)
	code:=datastruct.NULLError
	if err == nil {
		switch body.PlatformId{
		 case datastruct.PC_Platform:
			  fallthrough
		 case datastruct.WX_Platform:
			if body.Code == ""{
			  code=datastruct.JsonParseFailedFromPostBody
			}
		 default:
			code=datastruct.JsonParseFailedFromPostBody
		 }
		 if code == datastruct.NULLError{
		   var isExistRedis bool
		   var isExistMysql bool
		   var p_data *datastruct.PlayerData
		   conn:=handle.cacheHandler.GetConn()
		   defer conn.Close()
	
		   openid:=getOpenId(body.Code)
		   p_data,isExistRedis = handle.cacheHandler.GetPlayerData(conn,openid) //find in redis
		   if !isExistRedis{
			 p_data,isExistMysql = handle.dbHandler.GetPlayerData(openid) //find in mysql
			 if !isExistMysql{
				p_data= datastruct.CreateUser(openid,getPermissionId(body.IsAuth))
			 } else {
				refreshPlayerData(p_data,body.IsAuth)
			 }
			 handle.cacheHandler.SetPlayerAllData(conn,p_data)
		   } else {
			 refreshPlayerData(p_data,body.IsAuth)
			 handle.cacheHandler.SetPlayerSomeData(conn,p_data)
		   }
		   
		   c.JSON(200, gin.H{
			"code":code,
			"data":p_data,
		   })
		 } else {
			c.JSON(200, gin.H{
				"code":code,
			})
		 }
	}else{
	   code=datastruct.JsonParseFailedFromPostBody
	   c.JSON(200, gin.H{
		"code":code,
	   })
	}
}

func getOpenId(code string) string{
	 return code
}


func refreshPlayerData(p_data *datastruct.PlayerData,isauth int){
	if isauth == 1 && p_data.PermissionId == int(datastruct.Guest){
	  p_data.PermissionId = int(datastruct.Player)
	}
	p_data.UpdateTime = time.Now().Unix()
}

func (handle *EventHandler)fromRedisToMysql(token string){
	conn:=handle.cacheHandler.GetConn()
	defer conn.Close()
	p_data:=handle.cacheHandler.ReadPlayerData(conn,token)
	user_id:=handle.dbHandler.SetPlayerData(p_data)
	if p_data.Id<=0 && user_id > 0{
	  handle.cacheHandler.SetPlayerID(conn,token,user_id)
	}
}

func getPermissionId(isauth int) int{
	rs:= 1
	if isauth == 1{
		rs = 2
	}
	return rs
}

func (handle *EventHandler)UpdatePermisson(key string,permissionId int) datastruct.CodeType{
	code:=handle.cacheHandler.UpdatePermisson(key,permissionId)
	return code
}


func (handle *EventHandler)GetShopData(c *gin.Context){
	 var data datastruct.ShopData
	 code,plants:= handle.dbHandler.GetPlantsData()
	 data.Plants = plants 
	 c.JSON(200, gin.H{
		"code":int(code),
		"data":data,
	})
}






func (handle *EventHandler)Test1(c *gin.Context){
	var body datastruct.UserLogin
	c.BindJSON(&body)
	handle.fromRedisToMysql(body.Code)
	c.JSON(200, gin.H{
		"code": 0,
	})
}

func (handle *EventHandler)Test2(c *gin.Context){
	var body datastruct.UserLogin
	c.BindJSON(&body)
    handle.cacheHandler.TestMoney(body.Code)
	handle.fromRedisToMysql(body.Code)
	c.JSON(200, gin.H{
		"code": 0,
	})
}