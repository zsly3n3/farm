package event

import(
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
			body.Code = "test1"
		 case datastruct.WX_Platform:
			if body.Code == ""{
			 code=datastruct.JsonParseFailedFromPostBody
			}
		 default:
			code=datastruct.JsonParseFailedFromPostBody
		 }
		 if code == datastruct.NULLError{
		   //var isExist bool
		   var p_data *datastruct.PlayerData
		   //data,isExist=find in redis 
		   //data,tf=find in mysql
		   //create in redis
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