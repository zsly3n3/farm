package event

import (
	"farm/datastruct"
	"farm/thirdParty"
	"farm/tools"
	"time"

	"github.com/gin-gonic/gin"
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

			openid := getOpenId(body.Code, body.PlatformId)
			handle.PlayerIsOnline(openid)
			var tmpLoginData *datastruct.TmpLoginData
			p_data, isExistRedis = handle.cacheHandler.GetPlayerData(conn, openid) //find in redis
			if !isExistRedis {
				p_data, isExistMysql = handle.dbHandler.GetPlayerData(openid) //find in mysql
				if !isExistMysql {
					p_data = handle.createUser(openid, getPermissionId(body.IsAuth), body.NickName, body.Avatar)
					p_data.Referrer = body.Referrer
					p_data.Id = handle.dbHandler.SetPlayerData(p_data) //入库
				} else {
					tmpLoginData = handle.refreshPlayerData(p_data, body.IsAuth)
				}
				handle.cacheHandler.SetPlayerAllData(conn, p_data)
			} else {
				tmpLoginData = handle.refreshPlayerData(p_data, body.IsAuth)
				handle.cacheHandler.SetPlayerSomeData(conn, p_data)
			}
			c.JSON(200, gin.H{
				"code": code,
				"data": datastruct.ResponseLoginData(tmpLoginData, p_data, handle.soils, handle.plants, handle.petbars, handle.animals),
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

func getOpenId(code string, platform datastruct.Platform) string {
	if platform != datastruct.WX_Platform {
		return code
	}
	//select token from userinfo where code = openid
	return thirdParty.GetWXOpenID(code)
}

func (handle *EventHandler) refreshPlayerData(p_data *datastruct.PlayerData, isauth int) *datastruct.TmpLoginData {
	if isauth == 1 && p_data.PermissionId == int(datastruct.Guest) {
		p_data.PermissionId = int(datastruct.Player)
	}
	last_UpdateTime := p_data.UpdateTime
	current_UpdateTime := time.Now().Unix()

	tmpLoginData := new(datastruct.TmpLoginData)
	tmpLoginData.CD = 0
	var addGold int64
	addGold = 0
	animals := handle.animals
	plants := handle.plants
	speedFactor := p_data.InviteSpeedFactor
	if p_data.SpeedUp != nil {
		sec := p_data.SpeedUp.Ending - current_UpdateTime
		if sec > 0 {
			beforeSpeed_Sec := p_data.SpeedUp.Starting - last_UpdateTime
			if beforeSpeed_Sec > 0 {
				//normal 无加速计算 秒数为beforeSpeed_Sec
				addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, datastruct.DefaultSpeedUpFactor, beforeSpeed_Sec, plants, animals)
				//speed 加速计算 秒数为current_UpdateTime-p_data.SpeedUp.Starting
				addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, p_data.SpeedUp.Factor, current_UpdateTime-p_data.SpeedUp.Starting, plants, animals)
			} else {
				//speed 加速计算 秒数为current_UpdateTime-last_UpdateTime
				addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, p_data.SpeedUp.Factor, current_UpdateTime-last_UpdateTime, plants, animals)
			}
			tmpLoginData.CD = tools.EnableSpeedUp(p_data.SpeedUp.Ending, current_UpdateTime)
			tmpLoginData.Sec_EndingSpeedUp = p_data.SpeedUp.Ending - current_UpdateTime
		} else {
			if last_UpdateTime >= p_data.SpeedUp.Ending {
				//normal 无加速计算 秒数为current_UpdateTime-last_UpdateTime
				addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, datastruct.DefaultSpeedUpFactor, current_UpdateTime-last_UpdateTime, plants, animals)
			} else {
				afterSpeed_Sec := current_UpdateTime - p_data.SpeedUp.Ending //afterSpeed_Sec 为加速完成后还剩多少时间

				beforeSpeed_Sec := p_data.SpeedUp.Starting - last_UpdateTime //beforeSpeed_Sec 没有加速前的正常时间
				if beforeSpeed_Sec > 0 {
					//normal 无加速计算 秒数为beforeSpeed_Sec
					addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, datastruct.DefaultSpeedUpFactor, beforeSpeed_Sec, plants, animals)
					//speed 加速计算 秒数为p_data.SpeedUp.Ending - p_data.SpeedUp.Starting
					addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, p_data.SpeedUp.Factor, p_data.SpeedUp.Ending-p_data.SpeedUp.Starting, plants, animals)
					//normal 无加速计算 秒数为afterSpeed_Sec
					addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, datastruct.DefaultSpeedUpFactor, afterSpeed_Sec, plants, animals)
				} else {
					//speed 加速计算  p_data.SpeedUp.Ending - last_UpdateTime
					addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, p_data.SpeedUp.Factor, p_data.SpeedUp.Ending-last_UpdateTime, plants, animals)
					//normal 无加速计算 秒数为afterSpeed_Sec
					addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, datastruct.DefaultSpeedUpFactor, afterSpeed_Sec, plants, animals)
				}
			}
			p_data.SpeedUp = nil
		}
	} else {
		//normal 无加速计算 秒数为current_UpdateTime-last_UpdateTime
		addGold += tools.ComputeCurrentGold(speedFactor, p_data.Soil, p_data.PetBar, datastruct.DefaultSpeedUpFactor, current_UpdateTime-last_UpdateTime, plants, animals)
	}

	isGetedStamina := handle.dbHandler.IsGetStamina(p_data.Id)
	if !isGetedStamina && p_data.Stamina < datastruct.MaxStamina {
		p_data.Stamina = datastruct.MaxStamina
	}
	p_data.UpdateTime = current_UpdateTime
	p_data.GoldCount += addGold

	return tmpLoginData
}

func (handle *EventHandler) fromRedisToMysql(token string) {
	conn := handle.cacheHandler.GetConn()
	defer conn.Close()
	p_data := handle.cacheHandler.ReadPlayerData(conn, token)
	handle.dbHandler.SetPlayerData(p_data)
	// if p_data.Id <= 0 && user_id > 0 {
	// 	handle.cacheHandler.SetPlayerID(conn, token, user_id)
	// }
}

func getPermissionId(isauth int) int {
	rs := 1
	if isauth == 1 {
		rs = 2
	}
	return rs
}

func (handle *EventHandler) UpdatePermisson(key string, permissionId int, c *gin.Context) datastruct.CodeType {
	var body datastruct.UserAuthBody
	err := c.BindJSON(&body)
	if err != nil {
		return datastruct.JsonParseFailedFromPostBody
	}
	code, userId, referrer := handle.cacheHandler.UpdatePermisson(key, permissionId, &body)
	if referrer > 0 && referrer < userId {
		handle.dbHandler.InsertInviteInfo(userId, referrer)
	}
	return code
}

func (handle *EventHandler) GetInvitecount(key string) ([]*datastruct.ResponseInviteCount, datastruct.CodeType) {
	var userId int
	var code datastruct.CodeType
	var arr []*datastruct.ResponseInviteCount
	userId, code = handle.cacheHandler.GetUserId(key)
	if code != datastruct.NULLError {
		return nil, code
	}
	arr, code = handle.dbHandler.GetInvitecount(userId, 1)
	return arr, code
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
			code, resp_tmp = handle.cacheHandler.UpgradeSoil(key, &body, handle.soils, handle.plants, handle.animals)
		} else {
			code = datastruct.UpdateDataFailed
		}
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	return code, resp_tmp
}

func (handle *EventHandler) PlantInSoil(key string, c *gin.Context) (datastruct.CodeType, int64, string) {
	var body datastruct.PlantInSoil
	err := c.BindJSON(&body)
	code := datastruct.NULLError
	var gold int64
	var plantName string
	if err == nil {
		_, tf := handle.soils[body.SoilId]
		index := body.PlantId - 1
		if tf && index >= 0 && index < len(handle.plants) {
			code, gold, plantName = handle.cacheHandler.PlantInSoil(key, &body, handle.soils, handle.plants, handle.animals)
		} else {
			code = datastruct.UpdateDataFailed
		}
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	return code, gold, plantName
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
		code, gold, animal, soil_id = handle.cacheHandler.BuyPetbar(key, body.PetbarId, handle.petbars, handle.plants, handle.animals)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	return code, gold, animal, soil_id
}
func (handle *EventHandler) BuySoil(key string, c *gin.Context) (datastruct.CodeType, int64, int) {
	var body datastruct.BuySoil
	err := c.BindJSON(&body)
	if err != nil {
		return datastruct.JsonParseFailedFromPostBody, -1, -1
	}
	return handle.cacheHandler.BuySoil(key, body.SoilId, handle.soils, handle.plants, handle.animals)
}

func (handle *EventHandler) GetShopData(c *gin.Context, token string, soil_id int) {
	_, tf := handle.soils[soil_id]
	if !tf {
		c.JSON(200, gin.H{
			"code": datastruct.GetDataFailed,
		})
		return
	}
	conn := handle.cacheHandler.GetConn()
	defer conn.Close()
	code, plantlevel := handle.cacheHandler.GetPlantLevel(conn, token, soil_id)
	if code != datastruct.NULLError {
		c.JSON(200, gin.H{
			"code": code,
		})
		return
	}
	len := len(handle.plants)
	index := 0
	num := 40

	_, currentGold := handle.cacheHandler.ComputeCurrentGold(conn, token, handle.plants, handle.animals)
	plants := make([]*datastruct.ResponsePlant, 0, num)
	for i := 0; i < len; i++ {
		plant := new(datastruct.ResponsePlant)
		plant.Plant = handle.plants[i]
		if plantlevel >= plant.Level {
			plant.State = datastruct.Owned
			plants = append(plants, plant)
		} else if plantlevel+1 == plant.Level {
			if currentGold >= plant.Price {
				plant.State = datastruct.Unlocked
			} else {
				plant.State = datastruct.Locked
			}
			index = i + 1
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
	shopData.GoldCount = currentGold
	c.JSON(200, gin.H{
		"code": code,
		"data": shopData,
	})
}

func (handle *EventHandler) AddExpForAnimal(key string, c *gin.Context) (datastruct.CodeType, int64) {
	var body datastruct.AddExpForAnimal
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	var currentExp int64
	if err != nil {
		return datastruct.JsonParseFailedFromPostBody, -1
	}
	_, tf := handle.soils[body.SoilId]
	if !tf {
		return datastruct.UpdateDataFailed, -1
	}
	code, currentExp = handle.cacheHandler.AddExpForAnimal(key, &body, handle.petbars, handle.plants)
	return code, currentExp
}

func (handle *EventHandler) AnimalUpgrade(key string, c *gin.Context) (datastruct.CodeType, *datastruct.ResponseAnimalUpgrade) {
	var resp_data *datastruct.ResponseAnimalUpgrade
	resp_data = nil
	var body datastruct.BuyPetbar
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err != nil {
		return datastruct.JsonParseFailedFromPostBody, resp_data
	}
	code, resp_data = handle.cacheHandler.AnimalUpgrade(key, body.PetbarId, handle.petbars, handle.animals)
	return code, resp_data
}

func (handle *EventHandler) AddHoneyCount(key string) (datastruct.CodeType, *datastruct.ResponseAddHoney) {
	return handle.cacheHandler.AddHoneyCount(key)
}

func (handle *EventHandler) EnableCollectHoney(key string) (datastruct.CodeType, int64) {
	return handle.cacheHandler.EnableCollectHoney(key)
}

func (handle *EventHandler) GetStamina(key string) (datastruct.CodeType, *datastruct.ResponesStaminaData) {
	conn := handle.cacheHandler.GetConn()
	defer conn.Close()
	code, player_id, stamina := handle.cacheHandler.GetStamina(key, conn)
	if code != datastruct.NULLError {
		return code, nil
	}
	isGetedStamina := handle.dbHandler.IsGetStamina(player_id)
	if !isGetedStamina && stamina < datastruct.MaxStamina {
		stamina = datastruct.MaxStamina
		handle.cacheHandler.SetStamina(key, stamina, conn)
	}
	resp_data := new(datastruct.ResponesStaminaData)
	resp_data.Stamina = stamina

	now_Time := time.Now()
	tomorrow := now_Time.Add(24 * time.Hour)

	year, month, day := tomorrow.Date()
	tomorrow_Time := time.Date(year, month, day, 0, 0, 1, 0, time.Local)
	resp_data.NextRequest = tomorrow_Time.Unix() - now_Time.Unix()
	return datastruct.NULLError, resp_data
}

func (handle *EventHandler) Lottery(key string, c *gin.Context) (datastruct.CodeType, *datastruct.ResponesLotteryData, datastruct.RewardType) {
	var body datastruct.LotteryBody
	err := c.BindJSON(&body)
	if err != nil {
		return datastruct.JsonParseFailedFromPostBody, nil, -1
	}
	if body.RewardType < int(datastruct.Gold_10k) || body.RewardType > int(datastruct.Energy_UI2) {
		return datastruct.UpdateDataFailed, nil, -1
	}
	rewardType := datastruct.RewardType(body.RewardType)
	conn := handle.cacheHandler.GetConn()
	defer conn.Close()
	code, player_id, stamina := handle.cacheHandler.GetStamina(key, conn)
	if code != datastruct.NULLError {
		return code, nil, -1
	}
	isGetedStamina := handle.dbHandler.IsGetStamina(player_id)
	if !isGetedStamina && stamina < datastruct.MaxStamina {
		stamina = datastruct.MaxStamina
		handle.cacheHandler.SetStamina(key, stamina, conn)
	}
	if stamina < body.Expend {
		return datastruct.UpdateDataFailed, nil, -1
	}
	stamina -= body.Expend
	if rewardType != datastruct.Steal {
		_, goldCount := handle.cacheHandler.ComputeCurrentGold(conn, key, handle.plants, handle.animals)
		return handle.cacheHandler.LotteryNomal(key, rewardType, body.Expend, stamina, conn, goldCount)
	}
	//compute
	users := handle.dbHandler.LotterySteal(player_id)
	length := len(users)
	var player_data *datastruct.PlayerData
	if length <= 0 {
		player_data = handle.cacheHandler.ReadPlayerData(conn, key)
	} else {
		var user *datastruct.UserInfo
		randIndex := tools.RandInt(0, length)
		user = users[randIndex]
		if handle.cacheHandler.IsExistUserWithConn(conn, key) {
			player_data = handle.cacheHandler.ReadPlayerData(conn, key)
		} else {
			player_data = handle.dbHandler.GetPlayerDataFromDataBase(user)
		}
	}

	resp_data, addGold, addHoney := handle.computeSteal(player_data, body.Expend)
	handle.cacheHandler.ComputeCurrentGold(conn, key, handle.plants, handle.animals)
	code, resp_data = handle.cacheHandler.LotterySteal(key, addGold, addHoney, stamina, resp_data, conn)
	return code, resp_data, rewardType
}

func (handle *EventHandler) computeSteal(p_data *datastruct.PlayerData, expend int) (*datastruct.ResponesLotteryData, int64, int64) {
	//compute
	resp_data := new(datastruct.ResponesLotteryData)
	resp_data.Stolen = new(datastruct.ResponseStolen)
	resp_data.Stolen.Succeed = 1
	if p_data.Shield > 0 {
		resp_data.Stolen.Succeed = 0
	}
	var addGold int64
	var addHoney int64
	if resp_data.Stolen.Succeed == 1 {
		addGold = int64(1000 * expend)
		addHoney = int64(1000 * expend)
	} else {
		addGold = int64(100 * expend)
		addHoney = int64(100 * expend)
	}
	resp_data.Stolen.StolenGold = addGold
	resp_data.Stolen.StolenHoney = addHoney
	player_mp := make(map[string]interface{})
	farm_mp := make(map[string]interface{})
	farm_mp["goldcount"] = &(p_data.GoldCount)
	farm_mp["honeycount"] = &(p_data.HoneyCount)
	farm_mp["dogs"] = &(p_data.Shield)
	farm_mp["soil"] = datastruct.GetResponsePlayerSoil(p_data, handle.plants, handle.soils)
	farm_mp["petbar"] = datastruct.GetResponsePetbarData(p_data, handle.petbars, handle.animals)

	player_mp["nickname"] = &p_data.NickName
	player_mp["avatar"] = &p_data.Avatar
	player_mp["farm"] = farm_mp
	resp_data.Stolen.PlayerData = player_mp
	return resp_data, addGold, addHoney
}

func (handle *EventHandler) RefreshOnlineState(key string) datastruct.CodeType {
	return datastruct.NULLError
}

func (handle *EventHandler) PlayerIsOnline(token string) {

	now := time.Now()
	mm, _ := time.ParseDuration(datastruct.AddMinute)
	onlineTime := now.Add(mm)

	onlinePlayerData := new(datastruct.OnlinePlayerData)
	onlinePlayerData.OnlineTime = onlineTime.Unix()
	onlinePlayerData.WillDelete = false
	handle.onlinePlayers.Set(token, onlinePlayerData)
}

func (handle *EventHandler) IsExistUser(token string) bool {
	return handle.cacheHandler.IsExistUser(token)
}

func (handle *EventHandler) GoldDesc(token string) ([]*datastruct.ResponseGoldDesc, datastruct.CodeType) {
	var userId int
	var code datastruct.CodeType
	userId, code = handle.cacheHandler.GetUserId(token)
	if code != datastruct.NULLError {
		return nil, code
	}
	arr, code := handle.dbHandler.GoldDesc(userId)
	return arr, code
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
