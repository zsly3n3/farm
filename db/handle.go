package db

import (
	"farm/datastruct"
	"farm/log"
	"fmt"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) GetPlayerData(code string) (*datastruct.PlayerData, bool) {
	isExist := false
	var rs *datastruct.PlayerData
	user := new(datastruct.UserInfo)
	engine := handle.mysqlEngine
	has, _ := engine.Where("identity_id=?", code).Get(user)
	if has {
		isExist = true
		//add
		rs = new(datastruct.PlayerData)
		rs.Avatar = user.Avatar
		rs.CreatedAt = user.CreatedAt
		rs.Id = user.Id
		rs.NickName = user.NickName
		rs.PermissionId = user.PermissionId
		rs.Token = user.IdentityId
		rs.UpdateTime = user.UpdateTime
		var playerInfo datastruct.PlayerInfo
		engine.Id(rs.Id).Get(&playerInfo)
		rs.GoldCount = playerInfo.GoldCount
		rs.HoneyCount = playerInfo.HoneyCount
	}
	return rs, isExist
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
	sql := fmt.Sprintf("REPLACE INTO player_info (id,honey_count,gold_count,soil_level)VALUES(%d,%d,%d,%d)", userinfo.Id, p_data.HoneyCount, p_data.GoldCount, p_data.SoilLevel)
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
