package event

import (
	"farm/datastruct"
	"time"
	"github.com/gin-gonic/gin"
	"farm/tools"
	"farm/log"
)

func (handle *EventHandler) Login(c *gin.Context) {
	var body datastruct.UserLogin
	err := c.BindJSON(&body)
	code := datastruct.NULLError
	if err == nil {
		switch body.PlatformId {
		case datastruct.PC_Platform:
			fallthrough
		case datastruct.WX_Platform:
			if body.Code == "" {
				code = datastruct.JsonParseFailedFromPostBody
			}
		default:
			code = datastruct.JsonParseFailedFromPostBody
		}
		if code == datastruct.NULLError {
			var isExistRedis bool
			var isExistMysql bool
			var p_data *datastruct.PlayerData
			conn := handle.cacheHandler.GetConn()
			defer conn.Close()

			openid := getOpenId(body.Code)
			var tmpLoginData *datastruct.TmpLoginData
			p_data, isExistRedis = handle.cacheHandler.GetPlayerData(conn, openid) //find in redis
			if !isExistRedis {
				p_data, isExistMysql = handle.dbHandler.GetPlayerData(openid) //find in mysql
				if !isExistMysql {
					p_data = handle.createUser(openid, getPermissionId(body.IsAuth), "test", "avatar")
				} else {
					tmpLoginData=handle.refreshPlayerData(p_data, body.IsAuth)
				}
				handle.cacheHandler.SetPlayerAllData(conn, p_data)
			} else {
				tmpLoginData=handle.refreshPlayerData(p_data, body.IsAuth)
				handle.cacheHandler.SetPlayerSomeData(conn, p_data)
			}
			c.JSON(200, gin.H{
				"code": code,
				"data": datastruct.ResponseLoginData(tmpLoginData,p_data, handle.plants, handle.petbars, handle.animals),
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	} else {
		code = datastruct.JsonParseFailedFromPostBody
		c.JSON(200, gin.H{
			"code": code,
		})
	}
}

func getOpenId(code string) string {
	return code
}

func (handle *EventHandler) refreshPlayerData(p_data *datastruct.PlayerData, isauth int)*datastruct.TmpLoginData{
	if isauth == 1 && p_data.PermissionId == int(datastruct.Guest) {
		p_data.PermissionId = int(datastruct.Player)
	}
	last_UpdateTime := p_data.UpdateTime
	current_UpdateTime:=time.Now().Unix()
	

	gold := p_data.GoldCount
	for k, v := range handle.soils {
		if gold >= v.Price {
			p_data.Soil[k].State = datastruct.Unlocked
		}
	}
	for k, v := range handle.petbars {
		if gold >= v.Price {
			p_data.PetBar[k].State = datastruct.Unlocked
		}
	}
    
	tmpLoginData:=new(datastruct.TmpLoginData)
    if p_data.SpeedUp != nil{
		sec:=p_data.SpeedUp.Ending-current_UpdateTime
		if sec > 0{
			beforeSpeed_Sec:= p_data.SpeedUp.Starting - last_UpdateTime
			if beforeSpeed_Sec > 0{
			   //normal 无加速计算 秒数为beforeSpeed_Sec
			   //speed 加速计算 秒数为current_UpdateTime-p_data.SpeedUp.Starting
			} else {
			   //speed 加速计算 秒数为current_UpdateTime-last_UpdateTime
			}
			tmpLoginData.CD = tools.EnableSpeedUp(p_data.SpeedUp.Ending,current_UpdateTime)
			tmpLoginData.Sec_EndingSpeedUp = p_data.SpeedUp.Ending - current_UpdateTime
		} else{
		    if last_UpdateTime >= p_data.SpeedUp.Ending{
			   //normal 无加速计算 秒数为current_UpdateTime-last_UpdateTime
			} else {    
			 afterSpeed_Sec:=current_UpdateTime-p_data.SpeedUp.Ending //afterSpeed_Sec 为加速完成后还剩多少时间
			 log.Debug("afterSpeed_Sec:%d",afterSpeed_Sec)
		
			 beforeSpeed_Sec:= p_data.SpeedUp.Starting - last_UpdateTime //beforeSpeed_Sec 没有加速前的正常时间
			 if beforeSpeed_Sec > 0{
				//normal 无加速计算 秒数为beforeSpeed_Sec
				//speed 加速计算 秒数为p_data.SpeedUp.Ending - p_data.SpeedUp.Starting
				//normal 无加速计算 秒数为afterSpeed_Sec
			 }else{
				//speed 加速计算  p_data.SpeedUp.Ending - last_UpdateTime
				//normal 无加速计算 秒数为afterSpeed_Sec
			 }
			}
		    p_data.SpeedUp = nil
		}
	} else {
        //normal 无加速计算 秒数为current_UpdateTime-last_UpdateTime
	}
	p_data.UpdateTime = current_UpdateTime
	//compute 计算 金币，蜂蜜，体力(阈值30)，狗(盾牌)

	return tmpLoginData
}

func (handle *EventHandler) fromRedisToMysql(token string) {
	conn := handle.cacheHandler.GetConn()
	defer conn.Close()
	p_data := handle.cacheHandler.ReadPlayerData(conn, token)
	user_id := handle.dbHandler.SetPlayerData(p_data)
	if p_data.Id <= 0 && user_id > 0 {
		handle.cacheHandler.SetPlayerID(conn, token, user_id)
	}
}

func getPermissionId(isauth int) int {
	rs := 1
	if isauth == 1 {
		rs = 2
	}
	return rs
}

func (handle *EventHandler) UpdatePermisson(key string, permissionId int) datastruct.CodeType {
	code := handle.cacheHandler.UpdatePermisson(key, permissionId)
	return code
}

func (handle *EventHandler) UpgradeSoil(key string, c *gin.Context) (datastruct.CodeType, *datastruct.ResponseUpgradeSoil) {
	var body datastruct.UpgradeSoil
	err := c.BindJSON(&body)
	code := datastruct.NULLError
	var resp_tmp *datastruct.ResponseUpgradeSoil
	resp_tmp = nil
	if err == nil {
		_, tf := handle.soils[body.SoilId]
		if tf {
			code, resp_tmp = handle.cacheHandler.UpgradeSoil(key, &body, handle.soils)
		} else {
			code = datastruct.UpdateDataFailed
		}
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	return code, resp_tmp
}

func (handle *EventHandler) PlantInSoil(key string, c *gin.Context) (datastruct.CodeType, int64, string, int) {
	var body datastruct.PlantInSoil
	err := c.BindJSON(&body)
	code := datastruct.NULLError
	var gold int64
	var plantName string
	var soil_id int
	if err == nil {
		_, tf := handle.soils[body.SoilId]
		index := body.PlantId - 1
		if tf && index >= 0 && index < len(handle.plants) {
			code, gold, plantName, soil_id = handle.cacheHandler.PlantInSoil(key, &body, handle.plants, handle.soils)
		} else {
			code = datastruct.UpdateDataFailed
		}
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	return code, gold, plantName, soil_id
}

func (handle *EventHandler) BuyPetbar(key string, c *gin.Context) (datastruct.CodeType, int64, *datastruct.ResponseAnimal, int) {
	var body datastruct.BuyPetbar
	err := c.BindJSON(&body)
	code := datastruct.NULLError
	var gold int64
	var animal *datastruct.ResponseAnimal
	var soil_id int
	animal = nil
	if err == nil {
		code, gold, animal, soil_id = handle.cacheHandler.BuyPetbar(key, body.PetbarId, handle.petbars, handle.animals)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	return code, gold, animal, soil_id
}

func (handle *EventHandler) GetShopData(c *gin.Context, token string,soil_id int) {
	_, tf := handle.soils[soil_id]
    if !tf{
		c.JSON(200, gin.H{
			"code": datastruct.GetDataFailed,
		})
		return
	}
	code,plantlevel:= handle.cacheHandler.GetPlantLevel(token,soil_id)
	if code != datastruct.NULLError {
		c.JSON(200, gin.H{
			"code": code,
		})
		return
	}
	len := len(handle.plants)
	index := 0
	num := 40
	plants := make([]*datastruct.ResponsePlant, 0, num)
	for i := 0; i < len; i++ {
		plant := new(datastruct.ResponsePlant)
		plant.Plant = handle.plants[i]
		if plantlevel >= plant.Level {
			plant.State = datastruct.Owned
			plants = append(plants, plant)
		} else if plantlevel+1 == plant.Level {
			index = i + 1
			plant.State = datastruct.Unlocked
			plants = append(plants, plant)
			break
		}
	}

	for i := index; i < len; i++ {
		plant := new(datastruct.ResponsePlant)
		plant.Plant = handle.plants[i]
		plant.State = datastruct.Locked
		plants = append(plants, plant)
	}

	shopData := new(datastruct.ShopData)
	shopData.Plants = plants

	c.JSON(200, gin.H{
		"code": code,
		"data": shopData,
	})
}

func (handle *EventHandler) AddExpForAnimal(key string,c *gin.Context)(datastruct.CodeType,int64){
	var body datastruct.AddExpForAnimal
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	var currentExp int64
	if err != nil{
	   return datastruct.JsonParseFailedFromPostBody,-1
	}
	_, tf := handle.soils[body.SoilId]
	if !tf {
	   return datastruct.UpdateDataFailed,-1
	}
    code,currentExp = handle.cacheHandler.AddExpForAnimal(key,&body,handle.petbars,handle.plants)
	return code,currentExp
}

func (handle *EventHandler)AnimalUpgrade(key string,c *gin.Context)(datastruct.CodeType,*datastruct.ResponseAnimalUpgrade){
	var resp_data *datastruct.ResponseAnimalUpgrade
	resp_data = nil 
	var body datastruct.BuyPetbar
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err != nil{
	 return datastruct.JsonParseFailedFromPostBody,resp_data
	}
    code,resp_data = handle.cacheHandler.AnimalUpgrade(key,body.PetbarId,handle.petbars,handle.animals)
	return code,resp_data
}


func (handle *EventHandler)AddHoneyCount(key string)(datastruct.CodeType,*datastruct.ResponseAddHoney){
	var resp_data *datastruct.ResponseAddHoney
	resp_data = nil
	var code datastruct.CodeType
	code,resp_data = handle.cacheHandler.AddHoneyCount(key)
	return code,resp_data
}

func (handle *EventHandler)EnableCollectHoney(key string)(datastruct.CodeType,int64){
    return handle.cacheHandler.EnableCollectHoney(key)
}

func (handle *EventHandler) Test1(c *gin.Context) {
	var body datastruct.UserLogin
	c.BindJSON(&body)
	handle.fromRedisToMysql(body.Code)
	c.JSON(200, gin.H{
		"code": 0,
	})
}


func (handle *EventHandler) Test2(c *gin.Context) {
	var body datastruct.UserLogin
	c.BindJSON(&body)
	handle.cacheHandler.TestMoney(body.Code)
	handle.fromRedisToMysql(body.Code)
	c.JSON(200, gin.H{
		"code": 0,
	})
}
