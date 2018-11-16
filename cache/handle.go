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
	isExist := isExistUser(conn, code)
	if isExist {
		rs = handle.ReadPlayerData(conn, code)
	}
	return rs, isExist
}

func (handle *CACHEHandler) GetConn() redis.Conn {
	conn := handle.redisClient.Get()
	return conn
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
	_, err := conn.Do("hmset", key,
		datastruct.PermissionIdField, p_data.PermissionId,
		datastruct.UpdateTimeField, p_data.UpdateTime)
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

	conn.Send("MULTI")
	conn.Send("hmset", key,
		datastruct.IdField, p_data.Id,
		datastruct.GoldField, p_data.GoldCount,
		datastruct.HoneyField, p_data.HoneyCount,
		datastruct.PermissionIdField, p_data.PermissionId,
		datastruct.CreatedAtField, p_data.CreatedAt,
		datastruct.UpdateTimeField, p_data.UpdateTime,
		datastruct.NickNameField, p_data.NickName,
		datastruct.AvatarField, p_data.Avatar,
		datastruct.SoilLevelField, p_data.SoilLevel,
		datastruct.SpeedUpField, speedup_str,
		datastruct.StaminaField, p_data.Stamina)

	for k, v := range p_data.Soil {
		soiltableName := fmt.Sprintf("soil%d", k)
		value, isError := tools.PlayerSoilToString(v)
		if isError {
			log.Debug("CACHEHandler SetPlayerData PlayerSoilToString err:%s player:%s", soiltableName, key)
			return
		}
		conn.Send("hset", soiltableName, key, value)
	}

	for k, v := range p_data.PetBar {
		petbartableName := fmt.Sprintf("petbar%d", int(k))
		value, isError := tools.PlayerPetbarToString(v)
		if isError {
			log.Debug("CACHEHandler SetPlayerData PlayerPetbarToString err:%s player:%s", petbartableName, key)
			return
		}
		conn.Send("hset", petbartableName, key, value)
	}

	_, err := conn.Do("EXEC")

	if err != nil {
		log.Debug("CACHEHandler SetPlayerData err:%s", err.Error())
	}
}

func (handle *CACHEHandler) ReadPlayerData(conn redis.Conn, key string) *datastruct.PlayerData {
	rs := new(datastruct.PlayerData)
	//add
	value, err := redis.Values(conn.Do("hmget", key,
		datastruct.IdField, datastruct.GoldField, datastruct.HoneyField,
		datastruct.PermissionIdField, datastruct.CreatedAtField, datastruct.UpdateTimeField,
		datastruct.NickNameField, datastruct.AvatarField, datastruct.SpeedUpField,
		datastruct.StaminaField))
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
		}
	}

	len_soil := 5
	rs.Soil = make(map[int]*datastruct.PlayerSoil)
	for i := 1; i <= len_soil; i++ {
		soiltableName := fmt.Sprintf("soil%d", i)
		value, err := redis.String(conn.Do("hget", soiltableName, key))
		if err == nil {
			tmp, _ := tools.BytesToPlayerSoil([]byte(value))
			rs.Soil[i] = tmp
		}
	}

	len_petbar := int(datastruct.Deity)
	rs.PetBar = make(map[datastruct.AnimalType]*datastruct.PlayerPetbar)
	for i := int(datastruct.Sea); i <= len_petbar; i++ {
		petbartableName := fmt.Sprintf("petbar%d", i)
		value, err := redis.String(conn.Do("hget", petbartableName, key))
		if err == nil {
			tmp, _ := tools.BytesToPlayerPetbar([]byte(value))
			rs.PetBar[datastruct.AnimalType(i)] = tmp
		}
	}
	rs.Token = key
	return rs
}

func (handle *CACHEHandler) UpdatePermisson(key string, permissionId int) datastruct.CodeType {
	conn := handle.GetConn()
	defer conn.Close()
	if !isExistUser(conn, key) {
		return datastruct.TokenError
	}
	rep, err := conn.Do("hset", key, datastruct.PermissionIdField, permissionId)
	log.Debug("rep:%v", rep)
	code := datastruct.NULLError
	if err != nil {
		code = datastruct.UpdateDataFailed
		log.Debug("CACHEHandler UpdatePermisson err:%s", err.Error())
	}
	return code
}

func (handle *CACHEHandler) UpgradeSoil(key string, upgradeSoil *datastruct.UpgradeSoil, soils map[int]datastruct.SoilData) (datastruct.CodeType, *datastruct.ResponseUpgradeSoil) {
	conn := handle.GetConn()
	defer conn.Close()
	var resp_tmp *datastruct.ResponseUpgradeSoil
	resp_tmp = nil
	if !isExistUser(conn, key) {
		return datastruct.UpdateDataFailed, resp_tmp
	}

	code, gold := handle.ComputeCurrentGold(conn, key)
	if code != datastruct.NULLError {
		return datastruct.UpdateDataFailed, resp_tmp
	}

	soiltableName := fmt.Sprintf("soil%d", upgradeSoil.SoilId)
	value, err := redis.String(conn.Do("hget", soiltableName, key))
	if err != nil {
		return datastruct.GetDataFailed, resp_tmp
	}
	tmp, _ := tools.BytesToPlayerSoil([]byte(value))
	if tmp.State != datastruct.Owned || gold < tmp.UpgradeLevelPrice {
		resp_tmp := new(datastruct.ResponseUpgradeSoil)
		resp_tmp.Level = tmp.Level
		resp_tmp.UpgradePrice = tmp.UpgradeLevelPrice
		resp_tmp.Factor = tmp.Factor
		resp_tmp.GoldCount = gold
		return datastruct.NULLError, resp_tmp
	}
	gold, resp_tmp = tools.ComputeSoilLevelPrice(gold, upgradeSoil.Level, tmp)
	value, _ = tools.PlayerSoilToString(tmp)
	conn.Send("MULTI")
	conn.Send("hset", key, datastruct.GoldField, gold)
	conn.Send("hset", soiltableName, key, value)
	_, err = conn.Do("EXEC")

	if err != nil {
		log.Debug("CACHEHandler UpgradeSoil err:%s", err.Error())
		return datastruct.UpdateDataFailed, nil
	}

	return datastruct.NULLError, resp_tmp
}

func (handle *CACHEHandler) PlantInSoil(key string, plantInSoil *datastruct.PlantInSoil, plants []datastruct.Plant, soils map[int]datastruct.SoilData) (datastruct.CodeType, int64, string, int) {
	conn := handle.GetConn()
	defer conn.Close()
	if !isExistUser(conn, key) {
		return datastruct.UpdateDataFailed, -1, "", -1
	}

	soiltableName := fmt.Sprintf("soil%d", plantInSoil.SoilId)
	value, err := redis.String(conn.Do("hget", soiltableName, key))
	var player_soil *datastruct.PlayerSoil
	if err == nil {
		player_soil, _ = tools.BytesToPlayerSoil([]byte(value))
	} else {
		return datastruct.GetDataFailed, -1, "", -1
	}

	plant := plants[plantInSoil.PlantId-1]
	plantLevel := player_soil.PlantLevel
	if plantLevel >= plant.Level {
		return datastruct.UpdateDataFailed, -1, "", -1
	}

	code, gold := handle.ComputeCurrentGold(conn, key)
	if code != datastruct.NULLError {
		return datastruct.UpdateDataFailed, -1, "", -1
	}

	if gold < plant.Price {
		return datastruct.GoldIsNotEnoughForPlant, gold, plant.CName, -1
	}

	if plantLevel+1 == plant.Level {
		gold = gold - plant.Price
		plantLevel = plant.Level
		player_soil.PlantLevel = plantLevel
	} else {
		last_plant := plants[plant.Level-2]
		return datastruct.PlantRequireUnlock, gold, last_plant.CName, -1
	}

	player_soil.PlantId = plantInSoil.PlantId
	if player_soil.State != datastruct.Owned {
		soil := soils[plantInSoil.SoilId]
		if gold < soil.Price {
			return datastruct.GoldIsNotEnoughForSoil, gold, "", -1
		}
		value, err = redis.String(conn.Do("hget", key, datastruct.SoilLevelField))
		if err != nil {
			log.Debug("CACHEHandler PlantInSoil hget err:%s ,player:%s", err.Error(), key)
			return datastruct.GetDataFailed, -1, "", -1
		}
		soilLevel := tools.StringToInt(value)
		if soilLevel+1 == soil.Require {
			gold = gold - soil.Price
			soilLevel = soil.Require
			player_soil.State = datastruct.Owned
		} else {
			return datastruct.SoilRequireUnlock, gold, "", soil.LastId
		}
		conn.Send("MULTI")
		conn.Send("hmset", key,
			datastruct.GoldField, gold,
			datastruct.SoilLevelField, soilLevel)
	} else {
		conn.Send("MULTI")
		conn.Send("hset", key, datastruct.GoldField, gold)
	}

	value, isError := tools.PlayerSoilToString(player_soil)
	if isError {
		log.Debug("CACHEHandler PlantInSoil PlayerSoilToString err:%s player:%s", soiltableName, key)
		return datastruct.UpdateDataFailed, -1, "", -1
	}
	conn.Send("hset", soiltableName, key, value)

	_, err = conn.Do("EXEC")

	if err != nil {
		log.Debug("CACHEHandler PlantInSoil MULTI set data err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1, "", -1
	}

	return datastruct.NULLError, gold, "", -1
}

func (handle *CACHEHandler) BuyPetbar(key string, soid_id int, petbars map[datastruct.AnimalType]datastruct.PetbarData, animals map[datastruct.AnimalType]map[int]datastruct.Animal) (datastruct.CodeType, int64, *datastruct.ResponseAnimal, int) {
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
	if !isExistUser(conn, key) {
		return datastruct.UpdateDataFailed, -1, animal, soil_id
	}

	petbartableName := fmt.Sprintf("petbar%d", int(petbar_type))
	value, err := redis.String(conn.Do("hget", petbartableName, key))
	if err != nil {
		return datastruct.UpdateDataFailed, -1, animal, soil_id
	}
	rs_tmp, _ := tools.BytesToPlayerPetbar([]byte(value))
	if rs_tmp.State == datastruct.Owned {
		return datastruct.UpdateDataFailed, -1, animal, soil_id
	}

	code, gold := handle.ComputeCurrentGold(conn, key)
	if code != datastruct.NULLError {
		return datastruct.UpdateDataFailed, -1, animal, soil_id
	}

	if gold < tmp.Price {
		return datastruct.GoldIsNotEnoughForSoil, gold, animal, soil_id
	}

	value, err = redis.String(conn.Do("hget", key, datastruct.SoilLevelField))
	if err != nil {
		log.Debug("CACHEHandler BuyPetbar hget err:%s ,player:%s", err.Error(), key)
		return datastruct.GetDataFailed, -1, animal, soil_id
	}

	soilLevel := tools.StringToInt(value)
	if soilLevel+1 != tmp.Require {
		soil_id = tmp.LastId
		return datastruct.SoilRequireUnlock, gold, animal, soil_id
	}
	gold = gold - tmp.Price
	soilLevel = tmp.Require
	rs_tmp.State = datastruct.Owned
	animalNumber := 1
	rs_tmp.AnimalNumber = animalNumber
	rs_tmp.CurrentExp = 0

	value, _ = tools.PlayerPetbarToString(rs_tmp)
	conn.Send("MULTI")
	conn.Send("hset", key, datastruct.GoldField, gold)
	conn.Send("hset", petbartableName, key, value)
	_, err = conn.Do("EXEC")

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
	if !isExistUser(conn, key) {
		return datastruct.GetDataFailed, resp_data
	}

	petbartableName := fmt.Sprintf("petbar%d", int(petbar_type))
	value, err := redis.String(conn.Do("hget", petbartableName, key))
	if err != nil {
		return datastruct.GetDataFailed, resp_data
	}
	rs_tmp, _ := tools.BytesToPlayerPetbar([]byte(value))
	if rs_tmp.AnimalNumber == 0 || rs_tmp.State != datastruct.Owned {
		return datastruct.UpdateDataFailed, resp_data
	}
	value, err = redis.String(conn.Do("hget", key, datastruct.HoneyField))
	if err != nil {
		return datastruct.GetDataFailed, resp_data
	}
	honeyCount := tools.StringToInt64(value)
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

	value, _ = tools.PlayerPetbarToString(rs_tmp)
	conn.Send("MULTI")
	conn.Send("hset", key, datastruct.HoneyField, last_honey)
	conn.Send("hset", petbartableName, key, value)
	_, err = conn.Do("EXEC")

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

func (handle *CACHEHandler) ComputeCurrentGold(conn redis.Conn, key string) (datastruct.CodeType, int64) {
	value, err := redis.String(conn.Do("hget", key, datastruct.GoldField))
	code := datastruct.NULLError
	if err != nil {
		log.Debug("CACHEHandler ComputeCurrentGold err:%s", err.Error())
		return datastruct.UpdateDataFailed, -1
	}
	return code, tools.StringToInt64(value)
}

func (handle *CACHEHandler) clearData() {
	conn := handle.GetConn()
	defer conn.Close()
	conn.Do("flushdb")
}

func (handle *CACHEHandler) GetPlantLevel(key string, soil_id int) (datastruct.CodeType, int) {
	conn := handle.GetConn()
	defer conn.Close()
	if !isExistUser(conn, key) {
		return datastruct.GetDataFailed, -1
	}
	soiltableName := fmt.Sprintf("soil%d", soil_id)
	value, err := redis.String(conn.Do("hget", soiltableName, key))
	var player_soil *datastruct.PlayerSoil
	if err == nil {
		player_soil, _ = tools.BytesToPlayerSoil([]byte(value))
	} else {
		return datastruct.GetDataFailed, -1
	}
	return datastruct.NULLError, player_soil.PlantLevel
}

func isExistUser(conn redis.Conn, key string) bool {
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
	if !isExistUser(conn, key) {
		return datastruct.GetDataFailed, -1
	}

	petbartableName := fmt.Sprintf("petbar%d", int(petbar_type))
	value, err := redis.String(conn.Do("hget", petbartableName, key))
	if err != nil {
		return datastruct.GetDataFailed, -1
	}
	playerPetbar, _ := tools.BytesToPlayerPetbar([]byte(value))
	//没有购买宠物栏
	if playerPetbar.State != datastruct.Owned {
		return datastruct.UpdateDataFailed, -1
	}

	soiltableName := fmt.Sprintf("soil%d", body.SoilId)
	value, err = redis.String(conn.Do("hget", soiltableName, key))
	if err != nil {
		return datastruct.GetDataFailed, -1
	}
	player_soil, _ := tools.BytesToPlayerSoil([]byte(value))

	//没有植物可提供经验
	if player_soil.PlantId == 0 || player_soil.State != datastruct.Owned {
		return datastruct.UpdateDataFailed, -1
	}

	plant := plants[player_soil.PlantId-1]
	player_soil.PlantLevel = 0
	player_soil.PlantId = 0
	playerPetbar.CurrentExp += plant.ExpForAnimal

	soil_value, _ := tools.PlayerSoilToString(player_soil)
	petbar_value, _ := tools.PlayerPetbarToString(playerPetbar)
	conn.Send("MULTI")
	conn.Send("hset", soiltableName, key, soil_value)
	conn.Send("hset", petbartableName, key, petbar_value)
	_, err = conn.Do("EXEC")
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
	if !isExistUser(conn, key) {
		return datastruct.GetDataFailed, nil
	}
	value, err := redis.String(conn.Do("hget", key, datastruct.SpeedUpField))
	if err != nil {
		return datastruct.GetDataFailed, nil
	}
	var rs_tmp *datastruct.SpeedUpData
	resp_data := new(datastruct.ResponseAddHoney)
	now_Time := time.Now()
	if value == "" {
		rs_tmp = new(datastruct.SpeedUpData)
		rs_tmp.Factor = 2
		rs_tmp.Starting = now_Time.Unix()
		hh, _ := time.ParseDuration("4h")
		rs_tmp.Ending = now_Time.Add(hh).Unix()
	} else {
		rs_tmp, _ = tools.BytesToSpeedUp([]byte(value))
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

	value, err = redis.String(conn.Do("hget", key, datastruct.HoneyField))
	if err != nil {
		return datastruct.GetDataFailed, nil
	}
	honeyCount := tools.StringToInt64(value)

	//compute honeyCount
	honeyCount += 100

	resp_data.HoneyCount = honeyCount
	value, _ = tools.SpeedUpToString(rs_tmp)
	_, err = conn.Do("hmset", key, datastruct.HoneyField, honeyCount, datastruct.SpeedUpField, value)
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
	if !isExistUser(conn, key) {
		return datastruct.GetDataFailed, -1
	}
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

func (handle *CACHEHandler) GetStamina(key string) (datastruct.CodeType, int, int) {
	conn := handle.GetConn()
	defer conn.Close()
	if !isExistUser(conn, key) {
		return datastruct.GetDataFailed, -1, -1
	}
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

func (handle *CACHEHandler) TestMoney(key string) {
	conn := handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("hmset", key, datastruct.GoldField, 10000, datastruct.HoneyField, 2000)
	if err != nil {
		log.Debug("CACHEHandler TestMoney err:%s", err.Error())
	}
}
