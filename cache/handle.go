package cache

import (
	"farm/datastruct"
	"farm/log"
	"farm/tools"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func (handle *CACHEHandler) GetPlayerData(conn redis.Conn, code string) (*datastruct.PlayerData, bool) {
	var rs *datastruct.PlayerData
	isExist := handle.IsExistUserWithConn(conn, code)
	if isExist {
		rs = handle.ReadPlayerData(conn, code)
	}
	return rs, isExist
}

func (handle *CACHEHandler) IsRemoveGuest(key string) (bool, int) {
	conn := handle.GetConn()
	defer conn.Close()
	value, err := redis.Values(conn.Do("hmget", key, datastruct.PermissionIdField, datastruct.CreatedAtField, datastruct.IdField))
	length := len(value)
	if err != nil || length == 0 {
		log.Debug("CACHEHandler IsRemoveGuest err:%s ,player:%s", err.Error(), key)
		return false, -1
	}
	var permissionId int
	var createdAt int64
	var userId int
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			permissionId = tools.StringToInt(str)
		case 1:
			createdAt = tools.StringToInt64(str)
		case 2:
			userId = tools.StringToInt(str)
		}
	}
	if permissionId == int(datastruct.Guest) {
		tmp := time.Now().Unix() - createdAt
		if tmp >= datastruct.ExpiredTime {
			return true, userId
		}

	}
	return false, -1
}

func (handle *CACHEHandler) GetConn() redis.Conn {
	conn := handle.redisClient.Get()
	return conn
}

func (handle *CACHEHandler) GetUserId(key string) (int, datastruct.CodeType) {
	conn := handle.GetConn()
	defer conn.Close()
	value, err := redis.String(conn.Do("hget", key, datastruct.IdField))
	if err != nil {
		return -1, datastruct.GetDataFailed
	}
	return tools.StringToInt(value), datastruct.NULLError
}

func (handle *CACHEHandler) SetPlayerID(conn redis.Conn, key string, p_id int) {
	_, err := conn.Do("hset", key, datastruct.IdField, p_id)
	if err != nil {
		log.Debug("CACHEHandler SetPlayerID err:%s", err.Error())
	}
}

func (handle *CACHEHandler) SetPlayerSomeData(conn redis.Conn, p_data *datastruct.PlayerData) {
	key := p_data.Token
	//add
	_, err := conn.Do("hmset", key, datastruct.PermissionIdField, p_data.PermissionId, datastruct.UpdateTimeField, p_data.UpdateTime)
	if err != nil {
		log.Debug("CACHEHandler SetPlayerData err:%s", err.Error())
	}
}

func (handle *CACHEHandler) SetPlayerAllData(conn redis.Conn, p_data *datastruct.PlayerData) {
	key := p_data.Token
	//add

	speedup_str := ""
	if p_data.SpeedUp != nil {
		var isError bool
		speedup_str, isError = tools.SpeedUpToString(p_data.SpeedUp)
		if isError {
			log.Debug("CACHEHandler SetPlayerData SpeedUpToString err player:%s", key)
			return
		}
	}
	args := make([]interface{}, 0, 40)
	args = append(args, key)

	args = append(args, datastruct.IdField)
	args = append(args, p_data.Id)

	args = append(args, datastruct.GoldField)
	args = append(args, p_data.GoldCount)

	args = append(args, datastruct.HoneyField)
	args = append(args, p_data.HoneyCount)

	args = append(args, datastruct.PermissionIdField)
	args = append(args, p_data.PermissionId)

	args = append(args, datastruct.CreatedAtField)
	args = append(args, p_data.CreatedAt)

	args = append(args, datastruct.UpdateTimeField)
	args = append(args, p_data.UpdateTime)

	args = append(args, datastruct.NickNameField)
	args = append(args, p_data.NickName)

	args = append(args, datastruct.AvatarField)
	args = append(args, p_data.Avatar)

	args = append(args, datastruct.SoilLevelField)
	args = append(args, p_data.SoilLevel)

	args = append(args, datastruct.SpeedUpField)
	args = append(args, speedup_str)

	args = append(args, datastruct.StaminaField)
	args = append(args, p_data.Stamina)

	args = append(args, datastruct.ShieldField)
	args = append(args, p_data.Shield)

	args = append(args, datastruct.ReferrerField)
	args = append(args, p_data.Referrer)

	args = append(args, datastruct.InviteSpeedFactorField)
	args = append(args, p_data.InviteSpeedFactor)

	for k, v := range p_data.Soil {
		soiltableName := fmt.Sprintf("soil%d", k)
		value, isError := tools.PlayerSoilToString(v)
		if isError {
			log.Debug("CACHEHandler SetPlayerData PlayerSoilToString err:%s player:%s", soiltableName, key)
			return
		}
		args = append(args, soiltableName)
		args = append(args, value)
	}

	for k, v := range p_data.PetBar {
		petbartableName := fmt.Sprintf("petbar%d", int(k))
		value, isError := tools.PlayerPetbarToString(v)
		if isError {
			log.Debug("CACHEHandler SetPlayerData PlayerPetbarToString err:%s player:%s", petbartableName, key)
			return
		}
		args = append(args, petbartableName)
		args = append(args, value)
	}

	_, err := conn.Do("hmset", args...)

	if err != nil {
		log.Debug("CACHEHandler SetPlayerData err:%s", err.Error())
	}
}

func (handle *CACHEHandler) DeletedKeys(keys []interface{}, petbars map[datastruct.AnimalType]datastruct.PetbarData, soils map[int]datastruct.SoilData) {
	conn := handle.GetConn()
	defer conn.Close()
	log.Debug("DeletedKeys:%v", keys)
	_, err := conn.Do("del", keys...)
	if err != nil {
		log.Debug("CACHEHandler DeletedKeys err:%s", err.Error())
	}
}

func (handle *CACHEHandler) ReadPlayerData(conn redis.Conn, key string) *datastruct.PlayerData {
	rs := new(datastruct.PlayerData)
	//add
	value, err := redis.Values(conn.Do("hmget", key,
		datastruct.IdField, datastruct.GoldField, datastruct.HoneyField,
		datastruct.PermissionIdField, datastruct.CreatedAtField, datastruct.UpdateTimeField,
		datastruct.NickNameField, datastruct.AvatarField, datastruct.SpeedUpField,
		datastruct.StaminaField, datastruct.ShieldField, datastruct.ReferrerField,
		datastruct.InviteSpeedFactorField))
	length := len(value)
	if err != nil {
		log.Debug("CACHEHandler ReadPlayerData err:%s ,player:%s", err.Error(), key)
		return rs
	}
	for i := 0; i < length; i++ {
		tmp := value[i].([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			rs.Id = tools.StringToInt(str)
		case 1:
			rs.GoldCount = tools.StringToInt64(str)
		case 2:
			rs.HoneyCount = tools.StringToInt64(str)
		case 3:
			rs.PermissionId = tools.StringToInt(str)
		case 4:
			rs.CreatedAt = tools.StringToInt64(str)
		case 5:
			rs.UpdateTime = tools.StringToInt64(str)
		case 6:
			rs.NickName = str
		case 7:
			rs.Avatar = str
		case 8:
			rs.SpeedUp = nil
			if str != "" {
				rs.SpeedUp, _ = tools.BytesToSpeedUp(tmp)
			}
		case 9:
			rs.Stamina = tools.StringToInt(str)
		case 10:
			rs.Shield = tools.StringToInt(str)
		case 11:
			rs.Referrer = tools.StringToInt(str)
		case 12:
			rs.InviteSpeedFactor = tools.StringToInt(str)
		}
	}

	len_soil := 5
	rs.Soil = make(map[int]*datastruct.PlayerSoil)
	for i := 1; i <= len_soil; i++ {
		soiltableName := fmt.Sprintf("soil%d", i)
		value, err := redis.String(conn.Do("hget", key, soiltableName))
		if err == nil {
			tmp, _ := tools.BytesToPlayerSoil([]byte(value))
			rs.Soil[i] = tmp
		}
	}

	len_petbar := int(datastruct.Deity)
	rs.PetBar = make(map[datastruct.AnimalType]*datastruct.PlayerPetbar)
	for i := int(datastruct.Sea); i <= len_petbar; i++ {
		petbartableName := fmt.Sprintf("petbar%d", i)
		value, err := redis.String(conn.Do("hget", key, petbartableName))
		if err == nil {
			tmp, _ := tools.BytesToPlayerPetbar([]byte(value))
			rs.PetBar[datastruct.AnimalType(i)] = tmp
		}
	}
	rs.Token = key
	return rs
}

func (handle *CACHEHandler) UpdatePermisson(key string, permissionId int, body *datastruct.UserAuthBody) (datastruct.CodeType, int, int) {
	conn := handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("hmset", key, datastruct.PermissionIdField, permissionId, datastruct.NickNameField, body.NickName, datastruct.AvatarField, body.Avatar)

	if err != nil {
		log.Debug("CACHEHandler UpdatePermisson hmset err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1, -1
	}
	value, err := redis.Values(conn.Do("hmget", key, datastruct.IdField, datastruct.ReferrerField))
	if err != nil {
		log.Debug("CACHEHandler UpdatePermisson hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, -1, -1
	}

	var userId int
	var referrer int
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		switch i {
		case 0:
			userId = tools.StringToInt(string(tmp[:]))
		case 1:
			referrer = tools.StringToInt(string(tmp[:]))
		}
	}

	return datastruct.NULLError, userId, referrer
}

func (handle *CACHEHandler) UpgradeSoil(key string, upgradeSoil *datastruct.UpgradeSoil, soils map[int]datastruct.SoilData, plants []datastruct.Plant, animals map[datastruct.AnimalType]map[int]datastruct.Animal) (datastruct.CodeType, *datastruct.ResponseUpgradeSoil) {
	conn := handle.GetConn()
	defer conn.Close()
	var resp_tmp *datastruct.ResponseUpgradeSoil
	resp_tmp = nil

	code, gold := handle.ComputeCurrentGold(conn, key, plants, animals)

	if code != datastruct.NULLError {
		return datastruct.UpdateDataFailed, resp_tmp
	}

	soiltableName := fmt.Sprintf("soil%d", upgradeSoil.SoilId)
	value, err := redis.String(conn.Do("hget", key, soiltableName))
	if err != nil {
		return datastruct.GetDataFailed, resp_tmp
	}
	tmp, _ := tools.BytesToPlayerSoil([]byte(value))
	if tmp.State != datastruct.Bought || gold < tmp.UpgradeLevelPrice {
		resp_tmp := new(datastruct.ResponseUpgradeSoil)
		resp_tmp.Level = tmp.Level
		resp_tmp.UpgradePrice = tmp.UpgradeLevelPrice
		resp_tmp.Factor = tmp.Factor
		resp_tmp.GoldCount = gold
		return datastruct.NULLError, resp_tmp
	}
	gold, resp_tmp = tools.ComputeSoilLevelPrice(gold, upgradeSoil.Level, tmp)
	value, _ = tools.PlayerSoilToString(tmp)
	_, err = conn.Do("hmset", key, datastruct.GoldField, gold, soiltableName, value)

	if err != nil {
		log.Debug("CACHEHandler UpgradeSoil err:%s", err.Error())
		return datastruct.UpdateDataFailed, nil
	}

	return datastruct.NULLError, resp_tmp
}

func (handle *CACHEHandler) PlantInSoil(key string, plantInSoil *datastruct.PlantInSoil, soils map[int]datastruct.SoilData, plants []datastruct.Plant, animals map[datastruct.AnimalType]map[int]datastruct.Animal) (datastruct.CodeType, int64, string) {
	conn := handle.GetConn()
	defer conn.Close()
	soiltableName := fmt.Sprintf("soil%d", plantInSoil.SoilId)

	value, err := redis.String(conn.Do("hget", key, soiltableName))
	if err != nil {
		log.Debug("CACHEHandler PlantInSoil hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, -1, ""
	}

	player_soil, _ := tools.BytesToPlayerSoil([]byte(value))
	if player_soil.State != datastruct.Bought {
		return datastruct.SoilIsNotOwned, -1, ""
	}

	if player_soil.PlantId == plantInSoil.PlantId {
		return datastruct.RepeatPlant, -1, ""
	}

	code, gold := handle.ComputeCurrentGold(conn, key, plants, animals)
	if code != datastruct.NULLError {
		return datastruct.UpdateDataFailed, -1, ""
	}

	args := make([]interface{}, 0, 5)
	args = append(args, key)

	plant := plants[plantInSoil.PlantId-1]
	plantLevel := player_soil.PlantLevel
	if plantLevel >= plant.Level {
		return datastruct.UpdateDataFailed, -1, ""
	}

	if gold < plant.Price {
		return datastruct.GoldIsNotEnoughForPlant, gold, plant.CName
	}

	if plantLevel+1 == plant.Level {
		gold = gold - plant.Price
		plantLevel = plant.Level
		player_soil.PlantLevel = plantLevel
	} else {
		last_plant := plants[plant.Level-2]
		return datastruct.PlantRequireUnlock, gold, last_plant.CName
	}

	player_soil.PlantId = plantInSoil.PlantId

	player_soil_str, isError := tools.PlayerSoilToString(player_soil)
	if isError {
		log.Debug("CACHEHandler PlantInSoil PlayerSoilToString err:%s player:%s", soiltableName, key)
		return datastruct.UpdateDataFailed, -1, ""
	}
	args = append(args, datastruct.GoldField)
	args = append(args, gold)
	args = append(args, soiltableName)
	args = append(args, player_soil_str)

	_, err = conn.Do("hmset", args...)

	if err != nil {
		log.Debug("CACHEHandler PlantInSoil MULTI set data err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1, ""
	}

	return datastruct.NULLError, gold, ""
}

func (handle *CACHEHandler) BuySoil(key string, soil_id int, soils map[int]datastruct.SoilData, plants []datastruct.Plant, animals map[datastruct.AnimalType]map[int]datastruct.Animal) (datastruct.CodeType, int64, int) {
	var gold int64
	conn := handle.GetConn()
	defer conn.Close()
	soiltableName := fmt.Sprintf("soil%d", soil_id)
	value, err := redis.Values(conn.Do("hmget", key, soiltableName, datastruct.SoilLevelField))
	if err != nil {
		log.Debug("CACHEHandler BuySoil hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, -1, -1
	}

	var player_soil *datastruct.PlayerSoil
	var soilLevel int
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		switch i {
		case 0:
			player_soil, _ = tools.BytesToPlayerSoil(tmp)
		case 1:
			soilLevel = tools.StringToInt(string(tmp[:]))
		}
	}

	if player_soil.State == datastruct.Bought {
		return datastruct.RepeatBuy, -1, -1
	}

	code, gold := handle.ComputeCurrentGold(conn, key, plants, animals)
	if code != datastruct.NULLError {
		return datastruct.UpdateDataFailed, -1, -1
	}

	soil := soils[soil_id]
	if gold < soil.Price {
		return datastruct.GoldIsNotEnoughForSoil, gold, -1
	}
	if soilLevel+1 == soil.Require {
		gold = gold - soil.Price
		soilLevel = soil.Require
		player_soil.State = datastruct.Bought
	} else {
		return datastruct.SoilRequireUnlock, gold, soil.LastId
	}

	player_soil_str, isError := tools.PlayerSoilToString(player_soil)
	if isError {
		log.Debug("CACHEHandler BuySoil PlayerSoilToString err:%s player:%s", soiltableName, key)
		return datastruct.UpdateDataFailed, -1, -1
	}
	args := make([]interface{}, 0, 7)
	args = append(args, key)
	args = append(args, datastruct.SoilLevelField)
	args = append(args, soilLevel)
	args = append(args, datastruct.GoldField)
	args = append(args, gold)
	args = append(args, soiltableName)
	args = append(args, player_soil_str)

	_, err = conn.Do("hmset", args...)
	if err != nil {
		log.Debug("CACHEHandler BuySoil MULTI set data err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1, -1
	}
	return datastruct.NULLError, gold, -1
}

func (handle *CACHEHandler) BuyPetbar(key string, soid_id int, petbars map[datastruct.AnimalType]datastruct.PetbarData, plants []datastruct.Plant, animals map[datastruct.AnimalType]map[int]datastruct.Animal) (datastruct.CodeType, int64, *datastruct.ResponseAnimal, int) {
	var animal *datastruct.ResponseAnimal
	animal = nil
	var tmp *datastruct.PetbarData
	tmp = nil
	var soil_id int
	var petbar_type datastruct.AnimalType
	for k, v := range petbars {
		if v.Id == soid_id {
			tmp = &v
			petbar_type = k
			break
		}
	}
	if tmp == nil {
		return datastruct.UpdateDataFailed, -1, animal, soil_id
	}

	conn := handle.GetConn()
	defer conn.Close()

	petbartableName := fmt.Sprintf("petbar%d", int(petbar_type))

	value, err := redis.Values(conn.Do("hmget", key, petbartableName, datastruct.SoilLevelField))
	if err != nil {
		log.Debug("CACHEHandler BuyPetbar hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, -1, nil, -1
	}

	var rs_tmp *datastruct.PlayerPetbar
	var soilLevel int
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		switch i {
		case 0:
			rs_tmp, _ = tools.BytesToPlayerPetbar(tmp)
		case 1:
			soilLevel = tools.StringToInt(string(tmp[:]))
		}
	}

	if rs_tmp.State == datastruct.Bought {
		return datastruct.RepeatBuy, -1, animal, soil_id
	}

	code, gold := handle.ComputeCurrentGold(conn, key, plants, animals)
	if code != datastruct.NULLError {
		return datastruct.UpdateDataFailed, -1, animal, soil_id
	}

	if gold < tmp.Price {
		return datastruct.GoldIsNotEnoughForSoil, gold, animal, soil_id
	}

	if soilLevel+1 != tmp.Require {
		soil_id = tmp.LastId
		return datastruct.SoilRequireUnlock, gold, animal, soil_id
	}
	gold = gold - tmp.Price
	soilLevel = tmp.Require
	rs_tmp.State = datastruct.Bought
	animalNumber := 1
	rs_tmp.AnimalNumber = animalNumber
	rs_tmp.CurrentExp = 0

	playerPetbar_str, _ := tools.PlayerPetbarToString(rs_tmp)
	_, err = conn.Do("hmset", key, datastruct.GoldField, gold, datastruct.SoilLevelField, soilLevel, petbartableName, playerPetbar_str)

	if err != nil {
		log.Debug("CACHEHandler UpgradeSoil err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1, animal, soil_id
	}

	rs_ani := animals[petbar_type][animalNumber]
	animal = new(datastruct.ResponseAnimal)
	animal.CurrentExp = 0
	animal.Exp = rs_ani.Exp
	animal.InCome = rs_ani.InCome
	animal.Name = rs_ani.Name
	animal.HoneyCount = rs_ani.HoneyCount
	animal.IsLast = 0
	return datastruct.NULLError, gold, animal, soil_id
}

func (handle *CACHEHandler) AnimalUpgrade(key string, perbarId int, petbars map[datastruct.AnimalType]datastruct.PetbarData, animals map[datastruct.AnimalType]map[int]datastruct.Animal) (datastruct.CodeType, *datastruct.ResponseAnimalUpgrade) {
	var resp_data *datastruct.ResponseAnimalUpgrade
	resp_data = nil

	var tmp *datastruct.PetbarData
	tmp = nil
	var petbar_type datastruct.AnimalType
	for k, v := range petbars {
		if v.Id == perbarId {
			tmp = &v
			petbar_type = k
			break
		}
	}
	if tmp == nil {
		return datastruct.UpdateDataFailed, resp_data
	}

	conn := handle.GetConn()
	defer conn.Close()

	petbartableName := fmt.Sprintf("petbar%d", int(petbar_type))

	value, err := redis.Values(conn.Do("hmget", key, petbartableName, datastruct.HoneyField))
	if err != nil {
		log.Debug("CACHEHandler AnimalUpgrade hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, nil
	}

	var rs_tmp *datastruct.PlayerPetbar
	var honeyCount int64
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		switch i {
		case 0:
			rs_tmp, _ = tools.BytesToPlayerPetbar(tmp)
		case 1:
			honeyCount = tools.StringToInt64(string(tmp[:]))
		}
	}

	if rs_tmp.AnimalNumber == 0 || rs_tmp.State != datastruct.Bought {
		return datastruct.UpdateDataFailed, resp_data
	}

	animal := animals[petbar_type][rs_tmp.AnimalNumber]
	num_animal := len(animals)
	//最后一个的动物无法升级
	if animal.Number == num_animal {
		return datastruct.UpdateDataFailed, resp_data
	}
	resp_data = new(datastruct.ResponseAnimalUpgrade)
	if rs_tmp.CurrentExp < animal.Exp {
		resp_data.RightExp = rs_tmp.CurrentExp
		return datastruct.ExpIsNotFullForUpgradeAnimal, resp_data
	}
	if honeyCount < animal.HoneyCount {
		resp_data.HoneyCount = honeyCount
		return datastruct.HoneyCountIsNotEnoughForUpgradeAnimal, resp_data
	}
	new_number := rs_tmp.AnimalNumber + 1
	rs_tmp.AnimalNumber = new_number
	rs_tmp.CurrentExp = 0
	last_honey := honeyCount - animal.HoneyCount
	new_animal := animals[petbar_type][new_number]
	isLast := 0
	if new_animal.Number == num_animal {
		isLast = 1
	}

	playerPetbar_str, _ := tools.PlayerPetbarToString(rs_tmp)

	_, err = conn.Do("hmset", key, datastruct.HoneyField, last_honey, petbartableName, playerPetbar_str)

	if err != nil {
		log.Debug("CACHEHandler AnimalUpgrade err:%s", err.Error())
		return datastruct.UpdateDataFailed, nil
	}

	resp_data.HoneyCount = last_honey
	resp_animal := new(datastruct.ResponseAnimal)
	resp_animal.CurrentExp = 0
	resp_animal.Exp = new_animal.Exp
	resp_animal.InCome = new_animal.InCome
	resp_animal.Name = new_animal.Name
	resp_animal.HoneyCount = new_animal.HoneyCount
	resp_animal.IsLast = isLast
	resp_data.Animal = resp_animal
	return datastruct.NULLError, resp_data
}

func (handle *CACHEHandler) ComputeCurrentGold(conn redis.Conn, key string, plants []datastruct.Plant, animals map[datastruct.AnimalType]map[int]datastruct.Animal) (datastruct.CodeType, int64) {

	value, err := redis.Values(conn.Do("hmget", key, datastruct.GoldField, datastruct.UpdateTimeField, datastruct.SpeedUpField, datastruct.InviteSpeedFactorField))
	length := len(value)
	var currentGold int64
	var playerUpdateTime int64
	var currentSpeedUp *datastruct.SpeedUpData
	var inviteSpeedFactor int
	for i := 0; i < length; i++ {
		tmp := value[i].([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			currentGold = tools.StringToInt64(str)
		case 1:
			playerUpdateTime = tools.StringToInt64(str)
		case 2:
			currentSpeedUp = nil
			if str != "" {
				currentSpeedUp, _ = tools.BytesToSpeedUp(tmp)
			}
		case 3:
			inviteSpeedFactor = tools.StringToInt(str)
		}
	}

	len_soil := 5
	soils := make(map[int]*datastruct.PlayerSoil)
	for i := 1; i <= len_soil; i++ {
		soiltableName := fmt.Sprintf("soil%d", i)
		value, err := redis.String(conn.Do("hget", key, soiltableName))
		if err == nil {
			tmp, _ := tools.BytesToPlayerSoil([]byte(value))
			soils[i] = tmp
		}
	}

	len_petbar := int(datastruct.Deity)
	petBars := make(map[datastruct.AnimalType]*datastruct.PlayerPetbar)
	for i := int(datastruct.Sea); i <= len_petbar; i++ {
		petbartableName := fmt.Sprintf("petbar%d", i)
		value, err := redis.String(conn.Do("hget", key, petbartableName))
		if err == nil {
			tmp, _ := tools.BytesToPlayerPetbar([]byte(value))
			petBars[datastruct.AnimalType(i)] = tmp
		}
	}

	last_UpdateTime := playerUpdateTime
	current_UpdateTime := time.Now().Unix()
	var addGold int64
	addGold = 0

	var speedFactor int
	speedFactor = inviteSpeedFactor
	if currentSpeedUp != nil {
		sec := currentSpeedUp.Ending - current_UpdateTime
		if sec > 0 {
			beforeSpeed_Sec := currentSpeedUp.Starting - last_UpdateTime
			if beforeSpeed_Sec > 0 {
				//normal 无加速计算 秒数为beforeSpeed_Sec
				addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, datastruct.DefaultSpeedUpFactor, beforeSpeed_Sec, plants, animals)
				//speed 加速计算 秒数为current_UpdateTime-p_data.SpeedUp.Starting
				addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, currentSpeedUp.Factor, current_UpdateTime-currentSpeedUp.Starting, plants, animals)
			} else {
				//speed 加速计算 秒数为current_UpdateTime-last_UpdateTime
				addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, currentSpeedUp.Factor, current_UpdateTime-last_UpdateTime, plants, animals)
			}
		} else {
			if last_UpdateTime >= currentSpeedUp.Ending {
				//normal 无加速计算 秒数为current_UpdateTime-last_UpdateTime
				addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, datastruct.DefaultSpeedUpFactor, current_UpdateTime-last_UpdateTime, plants, animals)
			} else {
				afterSpeed_Sec := current_UpdateTime - currentSpeedUp.Ending //afterSpeed_Sec 为加速完成后还剩多少时间

				beforeSpeed_Sec := currentSpeedUp.Starting - last_UpdateTime //beforeSpeed_Sec 没有加速前的正常时间
				if beforeSpeed_Sec > 0 {
					//normal 无加速计算 秒数为beforeSpeed_Sec
					addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, datastruct.DefaultSpeedUpFactor, beforeSpeed_Sec, plants, animals)
					//speed 加速计算 秒数为p_data.SpeedUp.Ending - p_data.SpeedUp.Starting
					addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, currentSpeedUp.Factor, currentSpeedUp.Ending-currentSpeedUp.Starting, plants, animals)
					//normal 无加速计算 秒数为afterSpeed_Sec
					addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, datastruct.DefaultSpeedUpFactor, afterSpeed_Sec, plants, animals)
				} else {
					//speed 加速计算  p_data.SpeedUp.Ending - last_UpdateTime
					addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, currentSpeedUp.Factor, currentSpeedUp.Ending-last_UpdateTime, plants, animals)
					//normal 无加速计算 秒数为afterSpeed_Sec
					addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, datastruct.DefaultSpeedUpFactor, afterSpeed_Sec, plants, animals)
				}
			}
			currentSpeedUp = nil
		}
	} else {
		//normal 无加速计算 秒数为current_UpdateTime-last_UpdateTime
		addGold += tools.ComputeCurrentGold(speedFactor, soils, petBars, datastruct.DefaultSpeedUpFactor, current_UpdateTime-last_UpdateTime, plants, animals)
	}
	currentGold += addGold
	args := make([]interface{}, 0, 5)
	args = append(args, key)
	if currentSpeedUp == nil {
		args = append(args, datastruct.SpeedUpField)
		args = append(args, "")
	}
	args = append(args, datastruct.GoldField)
	args = append(args, currentGold)
	args = append(args, datastruct.UpdateTimeField)
	args = append(args, current_UpdateTime)
	_, err = conn.Do("hmset", args...)
	if err != nil {
		log.Debug("CACHEHandler ComputeCurrentGold err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1
	}
	return datastruct.NULLError, currentGold
}

func (handle *CACHEHandler) clearData() {
	conn := handle.GetConn()
	defer conn.Close()
	conn.Do("flushdb")
}

func (handle *CACHEHandler) GetPlantLevel(conn redis.Conn, key string, soil_id int) (datastruct.CodeType, int) {

	soiltableName := fmt.Sprintf("soil%d", soil_id)
	value, err := redis.String(conn.Do("hget", key, soiltableName))
	var player_soil *datastruct.PlayerSoil
	if err == nil {
		player_soil, _ = tools.BytesToPlayerSoil([]byte(value))
	} else {
		return datastruct.GetDataFailed, -1
	}
	return datastruct.NULLError, player_soil.PlantLevel
}

func (handle *CACHEHandler) IsExistUser(key string) bool {
	conn := handle.GetConn()
	defer conn.Close()
	isExist := false
	ilen, err := conn.Do("hlen", key)
	if err == nil && (ilen.(int64)) > 0 {
		isExist = true
	}
	return isExist
}

func (handle *CACHEHandler) IsExistUserWithConn(conn redis.Conn, key string) bool {
	isExist := false
	ilen, err := conn.Do("hlen", key)
	if err == nil && (ilen.(int64)) > 0 {
		isExist = true
	}
	return isExist
}

func (handle *CACHEHandler) AddExpForAnimal(key string, body *datastruct.AddExpForAnimal, petbars map[datastruct.AnimalType]datastruct.PetbarData, plants []datastruct.Plant) (datastruct.CodeType, int64) {
	var currentExp int64
	var tmp *datastruct.PetbarData
	tmp = nil
	var petbar_type datastruct.AnimalType
	for k, v := range petbars {
		if v.Id == body.PetbarId {
			tmp = &v
			petbar_type = k
			break
		}
	}
	if tmp == nil {
		return datastruct.UpdateDataFailed, -1
	}
	conn := handle.GetConn()
	defer conn.Close()

	petbartableName := fmt.Sprintf("petbar%d", int(petbar_type))
	soiltableName := fmt.Sprintf("soil%d", body.SoilId)

	value, err := redis.Values(conn.Do("hmget", key, petbartableName, soiltableName))
	if err != nil {
		log.Debug("CACHEHandler AddExpForAnimal hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, -1
	}

	var playerPetbar *datastruct.PlayerPetbar
	var player_soil *datastruct.PlayerSoil
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		switch i {
		case 0:
			playerPetbar, _ = tools.BytesToPlayerPetbar(tmp)
		case 1:
			player_soil, _ = tools.BytesToPlayerSoil(tmp)
		}
	}

	//没有购买宠物栏
	if playerPetbar.State != datastruct.Bought {
		return datastruct.UpdateDataFailed, -1
	}

	//没有植物可提供经验
	if player_soil.PlantId == 0 || player_soil.State != datastruct.Bought {
		return datastruct.UpdateDataFailed, -1
	}

	plant := plants[player_soil.PlantId-1]
	player_soil.PlantLevel = 0
	player_soil.PlantId = 0
	playerPetbar.CurrentExp += plant.ExpForAnimal

	soil_value, _ := tools.PlayerSoilToString(player_soil)
	petbar_value, _ := tools.PlayerPetbarToString(playerPetbar)
	_, err = conn.Do("hmset", key, soiltableName, soil_value, petbartableName, petbar_value)
	if err != nil {
		log.Debug("CACHEHandler AddExpForAnimal err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1
	}
	currentExp = playerPetbar.CurrentExp
	return datastruct.NULLError, currentExp
}

func (handle *CACHEHandler) AddHoneyCount(key string) (datastruct.CodeType, *datastruct.ResponseAddHoney) {
	conn := handle.GetConn()
	defer conn.Close()

	value, err := redis.Values(conn.Do("hmget", key, datastruct.SpeedUpField, datastruct.HoneyField))
	if err != nil {
		log.Debug("CACHEHandler AddHoneyCount hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, nil
	}

	var speedup_str string
	var honeyCount int64
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			speedup_str = str
		case 1:
			honeyCount = tools.StringToInt64(str)
		}
	}

	var rs_tmp *datastruct.SpeedUpData
	resp_data := new(datastruct.ResponseAddHoney)
	now_Time := time.Now()
	if speedup_str == "" {
		rs_tmp = new(datastruct.SpeedUpData)
		rs_tmp.Factor = 2
		rs_tmp.Starting = now_Time.Unix()
		hh, _ := time.ParseDuration("4h")
		rs_tmp.Ending = now_Time.Add(hh).Unix()
	} else {
		rs_tmp, _ = tools.BytesToSpeedUp([]byte(speedup_str))
		CD := tools.EnableSpeedUp(rs_tmp.Ending, now_Time.Unix())
		if CD > 0 {
			resp_data.CD = CD
			return datastruct.AddHoneyCD, resp_data
		}
		rs_tmp.Factor += 2
		rs_tmp.Ending += 4 * 3600
	}

	nextspeedcd := tools.EnableSpeedUp(rs_tmp.Ending, now_Time.Unix())
	resp_data.CD = nextspeedcd

	//compute honeyCount
	honeyCount += 100

	resp_data.HoneyCount = honeyCount
	speedup_str, _ = tools.SpeedUpToString(rs_tmp)
	_, err = conn.Do("hmset", key, datastruct.HoneyField, honeyCount, datastruct.SpeedUpField, speedup_str)
	if err != nil {
		log.Debug("CACHEHandler AddHoneyCount err:%s", err.Error())
		return datastruct.GetDataFailed, nil
	}

	resp_data.SpeedUp = new(datastruct.ResponesSpeedUpData)
	resp_data.SpeedUp.Factor = rs_tmp.Factor
	resp_data.SpeedUp.Ending = rs_tmp.Ending - now_Time.Unix()
	return datastruct.NULLError, resp_data
}

func (handle *CACHEHandler) EnableCollectHoney(key string) (datastruct.CodeType, int64) {
	conn := handle.GetConn()
	defer conn.Close()
	value, err := redis.String(conn.Do("hget", key, datastruct.SpeedUpField))
	if err != nil {
		return datastruct.GetDataFailed, -1
	}
	if value == "" {
		return datastruct.NULLError, 0
	}
	rs_tmp, _ := tools.BytesToSpeedUp([]byte(value))
	nowtime := time.Now().Unix()
	CD := tools.EnableSpeedUp(rs_tmp.Ending, nowtime)
	return datastruct.NULLError, CD
}

func (handle *CACHEHandler) GetStamina(key string, conn redis.Conn) (datastruct.CodeType, int, int) {
	value, err := redis.Values(conn.Do("hmget", key, datastruct.IdField, datastruct.StaminaField))
	length := len(value)
	if err != nil {
		log.Debug("CACHEHandler GetStamina err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, -1, -1
	}
	var player_id int
	var stamina int
	for i := 0; i < length; i++ {
		tmp := value[i].([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			player_id = tools.StringToInt(str)
		case 1:
			stamina = tools.StringToInt(str)
		}
	}
	return datastruct.NULLError, player_id, stamina
}

func (handle *CACHEHandler) SetStamina(key string, stamina int, conn redis.Conn) {
	_, err := conn.Do("hset", key, datastruct.StaminaField, stamina)
	if err != nil {
		log.Debug("CACHEHandler SetStamina err:%s", err.Error())
	}
}

func (handle *CACHEHandler) LotteryNomal(key string, rewardType datastruct.RewardType, expend int, stamina int, conn redis.Conn, currentGold int64) (datastruct.CodeType, *datastruct.ResponesLotteryData, datastruct.RewardType) {

	resp_data := new(datastruct.ResponesLotteryData)

	if rewardType == datastruct.Gold_10k || rewardType == datastruct.Gold_103k || rewardType == datastruct.Gold_48k || rewardType == datastruct.Gold_16k {
		var addGold int64
		switch rewardType {
		case datastruct.Gold_10k:
			addGold = 10 * 1000
		case datastruct.Gold_103k:
			addGold = 103 * 1000
		case datastruct.Gold_48k:
			addGold = 48 * 1000
		case datastruct.Gold_16k:
			addGold = 16 * 1000
		}
		currentGold += addGold * int64(expend)
		resp_data.GoldCount = currentGold
		resp_data.Stamina = stamina

		_, err := conn.Do("hmset", key, datastruct.GoldField, currentGold, datastruct.StaminaField, stamina)
		if err != nil {
			log.Debug("CACHEHandler Lottery hmset_0 err:%s", err.Error())
			return datastruct.UpdateDataFailed, nil, -1
		}
	} else if rewardType == datastruct.Energy_UI1 || rewardType == datastruct.Energy_UI2 {
		resp_data.GoldCount = currentGold
		stamina += expend * 1
		resp_data.Stamina = stamina
		handle.SetStamina(key, stamina, conn)
	} else if rewardType == datastruct.Gog {
		resp_data.GoldCount = currentGold
		resp_data.Stamina = stamina
		value, _ := redis.String(conn.Do("hget", key, datastruct.ShieldField))
		shields := tools.StringToInt(value)
		shields += expend * 1
		if shields >= datastruct.MaxShield {
			shields = datastruct.MaxShield
		}
		resp_data.Shield = shields
		_, err := conn.Do("hmset", key, datastruct.ShieldField, shields, datastruct.StaminaField, stamina)
		if err != nil {
			log.Debug("CACHEHandler Lottery hmset_1 err:%s", err.Error())
			return datastruct.UpdateDataFailed, nil, -1
		}
	}
	return datastruct.NULLError, resp_data, rewardType
}

func (handle *CACHEHandler) LotterySteal(key string, addGold int64, addHoney int64, stamina int, resp_data *datastruct.ResponesLotteryData, conn redis.Conn) (datastruct.CodeType, *datastruct.ResponesLotteryData) {

	value, err := redis.Values(conn.Do("hmget", key, datastruct.GoldField, datastruct.HoneyField))
	if err != nil {
		log.Debug("CACHEHandler LotterySteal hmget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, nil
	}

	var rs_goldCount int64
	var rs_honeyCount int64
	for i := 0; i < len(value); i++ {
		tmp := value[i].([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			rs_goldCount = tools.StringToInt64(str)
		case 1:
			rs_honeyCount = tools.StringToInt64(str)
		}
	}

	rs_goldCount += addGold
	rs_honeyCount += addHoney

	_, err = conn.Do("hmset", key, datastruct.GoldField, rs_goldCount, datastruct.StaminaField, stamina, datastruct.HoneyField, rs_honeyCount)
	if err != nil {
		log.Debug("CACHEHandler LotterySteal  err:%s", err.Error())
		return datastruct.UpdateDataFailed, nil
	}
	resp_data.GoldCount = rs_goldCount
	resp_data.Stamina = stamina
	return datastruct.NULLError, resp_data
}

func (handle *CACHEHandler) TestMoney(key string) {
	conn := handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("hmset", key, datastruct.GoldField, 10000, datastruct.HoneyField, 2000)
	if err != nil {
		log.Debug("CACHEHandler TestMoney err:%s", err.Error())
	}
}
