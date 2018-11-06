package event

import(
	"time"
	"farm/datastruct"
	"github.com/gin-gonic/gin"
	//"farm/log"
)

func (handle *EventHandler) Login(c *gin.Context){
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
		   p_data,isExistRedis = handle.cacheHandler.GetPlayerData(body.Code) //find in redis
		   if !isExistRedis{
			 p_data,isExistMysql = handle.dbHandler.GetPlayerData(body.Code) //find in mysql
			 if !isExistMysql{
				p_data = datastruct.CreateUser(body.Code,false)
			 } else {
				refreshPlayerData(p_data)
			 }
			 handle.cacheHandler.SetPlayerData(p_data)
		   } else {
			 refreshPlayerData(p_data)
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

func refreshPlayerData(p_data *datastruct.PlayerData){
	p_data.UpdateTime = time.Now().Unix()   
}

func (handle *EventHandler)fromRedisToMysql(token string){
	conn:=handle.cacheHandler.GetConn()
	defer conn.Close()
	p_data:=handle.cacheHandler.ReadPlayerData(conn,token)
	handle.dbHandler.SetPlayerData(p_data)
}

func (handle *EventHandler)Test1(c *gin.Context){
	var body datastruct.UserLogin
	c.BindJSON(&body)
	handle.fromRedisToMysql(body.Code)
	c.JSON(200, gin.H{
		"code": 0,
	})
}