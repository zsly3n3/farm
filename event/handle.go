package event

import(
	"time"
	"farm/datastruct"
	"github.com/gin-gonic/gin"
)

func (handle *EventHandler) Login(c *gin.Context){
	var body datastruct.UserLogin
	err:=c.BindJSON(&body)
	code:=datastruct.NULLError
	if err == nil {
		 switch body.PlatformId{
		 case datastruct.PC_Platform:
			body.Code = "test2"
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
				p_data = createUser(body.Code,false)
			 }
			 handle.cacheHandler.SetPlayerData(p_data)
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

func createUser(code string,isAuth bool)*datastruct.PlayerData{
	 player:=new(datastruct.PlayerData)
	 timestamp:=time.Now().Unix()
	 player.IsAuth = isAuth
	 player.CreatedAt = timestamp
	 player.UpdateTime = timestamp
	 player.IdentityId = code
	 player.GoldCount = 0
	 player.HoneyCount = 0
	 return player
}