package event

import(
	"farm/datastruct"
	"github.com/gin-gonic/gin"
	"farm/log"
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
		   //var isExistRedis bool  //test
		   var isExistMysql bool
		   var p_data *datastruct.PlayerData
		   //p_data,isExistRedis = handle.cacheHandler.GetPlayerData(body.Code) //find in redis test
		   // if !isExistRedis{ test
			 p_data,isExistMysql = handle.dbHandler.GetPlayerData(body.Code) //find in mysql
			 log.Debug("%v",p_data) //test
			 if !isExistMysql{
				p_data = datastruct.CreateUser(body.Code,false)
			 }
			 handle.cacheHandler.SetPlayerData(p_data)
		   //} test
		   handle.fromRedisToMysql(body.Code) //test
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

func (handle *EventHandler)fromRedisToMysql(token string){
	conn:=handle.cacheHandler.GetConn()
	defer conn.Close()
	p_data:=handle.cacheHandler.ReadPlayerData(conn,token)
	handle.dbHandler.SetPlayerData(p_data)
}