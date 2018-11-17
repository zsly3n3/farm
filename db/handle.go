package db

import (
	"farm/datastruct"
	"farm/log"
	"farm/tools"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) GetPlayerData(code string, soils map[int]datastruct.SoilData) (*datastruct.PlayerData, bool) {
	isExist := false
	var p_data *datastruct.PlayerData
	user := new(datastruct.UserInfo)
	engine := handle.mysqlEngine
	has, _ := engine.Where("identity_id=?", code).Get(user)
	if has {
		isExist = true
		p_data = handle.getPlayerData(user, soils)
	}
	return p_data, isExist
}

func (handle *DBHandler) SetPlayerData(p_data *datastruct.PlayerData) int {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	//add
	var userinfo datastruct.UserInfo
	userinfo.IdentityId = p_data.Token
	userinfo.CreatedAt = p_data.CreatedAt
	userinfo.PermissionId = p_data.PermissionId
	userinfo.UpdateTime = p_data.UpdateTime
	userinfo.Avatar = p_data.Avatar
	userinfo.NickName = p_data.NickName

	var err error
	if p_data.Id <= 0 {
		_, err = session.Insert(&userinfo)
	} else {
		var tmp datastruct.UserInfo
		var has bool
		has, err = session.Where("id=?", p_data.Id).Get(&tmp)
		if has {
			_, err = session.Where("id=?", p_data.Id).Update(&userinfo)
			userinfo.Id = p_data.Id
		} else {
			_, err = session.Insert(&userinfo)
		}
	}
	if err != nil {
		str := fmt.Sprintf("DBHandler->SetPlayerData InsertOrUpdate UserInfo :%s", err.Error())
		rollback(str, session)
		return userinfo.Id
	}
	sql := fmt.Sprintf("REPLACE INTO player_info (id,honey_count,gold_count,soil_level,stamina,shield)VALUES(%d,%d,%d,%d,%d,%d)", userinfo.Id, p_data.HoneyCount, p_data.GoldCount, p_data.SoilLevel, p_data.Stamina, p_data.Shield)
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->SetPlayerData REPLACE PlayerInfo :%s", err.Error())
		rollback(str, session)
		return userinfo.Id
	}

	if p_data.SpeedUp != nil {
		sql := fmt.Sprintf("REPLACE INTO player_speed_up (id,factor,starting,ending)VALUES(%d,%d,%d,%d)", userinfo.Id, p_data.SpeedUp.Factor, p_data.SpeedUp.Starting, p_data.SpeedUp.Ending)
		_, err = session.Exec(sql)
		if err != nil {
			str := fmt.Sprintf("DBHandler->SetPlayerData REPLACE player_speed_up :%s", err.Error())
			rollback(str, session)
			return userinfo.Id
		}
	}

	for k, v := range p_data.PetBar {
		sql = fmt.Sprintf("REPLACE INTO petbar%d (p_id,animal_number,current_exp,state)VALUES(%d,%d,%d,%d)", int(k), userinfo.Id, v.AnimalNumber, v.CurrentExp, int(v.State))
		_, err = session.Exec(sql)
		if err != nil {
			str := fmt.Sprintf("DBHandler->SetPlayerData REPLACE PetBar :%s", err.Error())
			rollback(str, session)
			return userinfo.Id
		}
	}

	for k, v := range p_data.Soil {
		sql = fmt.Sprintf("REPLACE INTO soil%d (p_id,level,plant_id,upgrade_level_price,factor,state,plant_level)VALUES(%d,%d,%d,%d,%d,%d,%d)", k, userinfo.Id, v.Level, v.PlantId, v.UpgradeLevelPrice, v.Factor, int(v.State), v.PlantLevel)
		_, err = session.Exec(sql)
		if err != nil {
			str := fmt.Sprintf("DBHandler->SetPlayerData REPLACE PetBar :%s", err.Error())
			rollback(str, session)
			return userinfo.Id
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->SetPlayerData Commit :%s", err.Error())
		rollback(str, session)
	}
	return userinfo.Id
}

func rollback(err_str string, session *xorm.Session) {
	log.Debug("will rollback,err_str:%v", err_str)
	session.Rollback()
}

func (handle *DBHandler) InsertRewardStamina(player_id int) {
	engine := handle.mysqlEngine
	var rewardStamina datastruct.RewardStamina
	rewardStamina.Id = player_id
	rewardStamina.GetTime = time.Now().Unix()
	_, err := engine.Insert(&rewardStamina)
	if err != nil {
		log.Debug("InsertRewardStamina error:%v", err.Error())
	}
	return
}

func (handle *DBHandler) IsGetStamina(player_id int) bool {
	engine := handle.mysqlEngine
	var rewardStamina datastruct.RewardStamina
	has, _ := engine.Where("id=?", player_id).Get(&rewardStamina)
	return has
}

func (handle *DBHandler) LotterySteal(player_id int, soils map[int]datastruct.SoilData) *datastruct.PlayerData {
	//compute
	engine := handle.mysqlEngine
	users := make([]*datastruct.UserInfo, 0)
	engine.Where("id <> ?", player_id).Find(&users)
	var user datastruct.UserInfo
	length := len(users)
	if length <= 0 {
		engine.Where("id = ?", player_id).Get(&user)
	} else {
		randIndex := tools.RandInt(0, length)
		user = *(users[randIndex])
	}
	p_data := handle.getPlayerData(&user, soils)
	return p_data
}

func (handle *DBHandler) getPlayerData(user *datastruct.UserInfo, soils map[int]datastruct.SoilData) *datastruct.PlayerData {
	engine := handle.mysqlEngine
	p_data := new(datastruct.PlayerData)
	p_data.Avatar = user.Avatar
	p_data.CreatedAt = user.CreatedAt
	p_data.Id = user.Id
	p_data.NickName = user.NickName
	p_data.PermissionId = user.PermissionId
	p_data.Token = user.IdentityId
	p_data.UpdateTime = user.UpdateTime

	var playerInfo datastruct.PlayerInfo
	engine.Id(p_data.Id).Get(&playerInfo)
	p_data.GoldCount = playerInfo.GoldCount
	p_data.HoneyCount = playerInfo.HoneyCount
	p_data.Stamina = playerInfo.Stamina
	p_data.Shield = playerInfo.Shield
	p_data.SoilLevel = playerInfo.SoilLevel

	soil_mp := make(map[int]*datastruct.PlayerSoil)
	petBar_mp := make(map[datastruct.AnimalType]*datastruct.PlayerPetbar)

	var soil_1 datastruct.Soil1
	engine.Where("p_id = ?", p_data.Id).Find(&soil_1)
	soild_data := soils[1]
	soil_mp[1] = tools.CreatePlayerSoil1(&soil_1, &soild_data)

	var soil_2 datastruct.Soil2
	engine.Where("p_id = ?", p_data.Id).Find(&soil_2)
	soild_data = soils[2]
	soil_mp[2] = tools.CreatePlayerSoil2(&soil_2, &soild_data)

	var soil_3 datastruct.Soil3
	engine.Where("p_id = ?", p_data.Id).Find(&soil_3)
	soild_data = soils[3]
	soil_mp[3] = tools.CreatePlayerSoil3(&soil_3, &soild_data)

	var soil_4 datastruct.Soil4
	engine.Where("p_id = ?", p_data.Id).Find(&soil_4)
	soild_data = soils[4]
	soil_mp[4] = tools.CreatePlayerSoil4(&soil_4, &soild_data)

	var soil_5 datastruct.Soil5
	engine.Where("p_id = ?", p_data.Id).Find(&soil_5)
	soild_data = soils[5]
	soil_mp[5] = tools.CreatePlayerSoil5(&soil_5, &soild_data)

	// PetBar       map[AnimalType]*PlayerPetbar //宠物栏信息
	// SpeedUp      *SpeedUpData                 //全局加速数据
	p_data.Soil = soil_mp
	p_data.PetBar = petBar_mp
	return p_data
}

func (handle *DBHandler) GetPlantsSlice() []datastruct.Plant {
	engine := handle.mysqlEngine
	plants := make([]datastruct.Plant, 0)
	err := engine.Asc("level").Find(&plants)
	if err != nil {
		log.Debug("GetPlantsSlice error:%v", err.Error())
	}
	return plants
}

func (handle *DBHandler) GetAnimalsMap() map[datastruct.AnimalType]map[int]datastruct.Animal {
	engine := handle.mysqlEngine
	mp := make(map[datastruct.AnimalType]map[int]datastruct.Animal)
	start := int(datastruct.Sea)
	end := int(datastruct.Deity)
	for i := start; i <= end; i++ {
		arr := make([]datastruct.Animal, 0)
		err := engine.Asc("number").Where("class_id = ?", i).Find(&arr)
		if err != nil {
			log.Debug("GetAnimalsMap error:%v", err.Error())
			return nil
		}
		tmp_mp := make(map[int]datastruct.Animal)
		for _, v := range arr {
			tmp_mp[v.Number] = v
		}
		mp[datastruct.AnimalType(i)] = tmp_mp
	}
	return mp
}
