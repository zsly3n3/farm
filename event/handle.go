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
		  conn:=handle.cacheHandler.GetConn()
		  defer conn.Close()
		   var isExistRedis bool
		   var isExistMysql bool
		   var p_data *datastruct.PlayerData
		   p_data,isExistRedis = handle.cacheHandler.GetPlayerData(conn,body.Code) //find in redis
		   if !isExistRedis{
			 p_data,isExistMysql = handle.dbHandler.GetPlayerData(body.Code) //find in mysql
			 if !isExistMysql{
				p_data = datastruct.CreateUser(body.Code,getPermissionId(body.IsAuth))
			 } else {
				refreshPlayerData(p_data,body.IsAuth)
			 }
		   } else {
			 refreshPlayerData(p_data,body.IsAuth)
		   }
		   handle.cacheHandler.SetPlayerData(conn,p_data)
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