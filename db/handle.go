package db

import (
	"farm/datastruct"
	"farm/log"
	"farm/tools"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) GetPlayerData(code string) (*datastruct.PlayerData, bool) {
	isExist := false
	var p_data *datastruct.PlayerData
	user := new(datastruct.UserInfo)
	engine := handle.mysqlEngine
	has, _ := engine.Where("identity_id=?", code).Get(user)
	if has {
		isExist = true
		p_data = handle.GetPlayerDataFromDataBase(user)
	}
	return p_data, isExist
}

func (handle *DBHandler) DeleteUser(userId int, soils map[int]datastruct.SoilData, petbars map[datastruct.AnimalType]datastruct.PetbarData) {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	user := new(datastruct.UserInfo)
	_, err := session.Id(userId).Delete(user)
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeleteUser UserInfo :%s", err.Error())
		rollback(str, session)
	}

	player := new(datastruct.PlayerInfo)
	_, err = session.Id(userId).Delete(player)
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeleteUser playerInfo :%s", err.Error())
		rollback(str, session)
	}

	speedup := new(datastruct.PlayerSpeedUp)
	_, err = session.Id(userId).Delete(speedup)
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeleteUser playerSpeedUp :%s", err.Error())
		rollback(str, session)
	}

	for k, _ := range petbars {
		sql := fmt.Sprintf("DELETE FROM petbar%d WHERE p_id = %d", int(k), userId)
		_, err = session.Exec(sql)
		if err != nil {
			str := fmt.Sprintf("DBHandler->SetPlayerData DELETE FROM PetBar%d :%s", int(k), err.Error())
			rollback(str, session)
		}
	}
	for k, _ := range soils {
		sql := fmt.Sprintf("DELETE FROM soil%d WHERE p_id = %d", k, userId)
		_, err = session.Exec(sql)
		if err != nil {
			str := fmt.Sprintf("DBHandler->SetPlayerData DELETE FROM soil%d :%s", k, err.Error())
			rollback(str, session)
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeleteUser Commit :%s", err.Error())
		rollback(str, session)
	}
}

func (handle *DBHandler) InsertInviteInfo(userId int, referrer int) {
	engine := handle.mysqlEngine
	var invite datastruct.InviteInfo
	var user datastruct.UserInfo
	var has bool
	has, _ = engine.Id(referrer).Get(&user)
	if has && user.PermissionId == int(datastruct.Player) {
		has, _ = engine.Where("received = ?", userId).Get(&invite)
		if !has {
			invite.Received = userId
			invite.Sended = referrer
			invite.CreatedAt = time.Now().Unix()
			_, err := engine.Insert(&invite)
			if err != nil {
				log.Debug("InsertInviteInfo insert error:%v", err.Error())
			}
		}
	}
}

func (handle *DBHandler) GetInvitecount(userId int, inviteSpeedFactor int) ([]*datastruct.ResponseInviteCount, datastruct.CodeType) {

	engine := handle.mysqlEngine
	users := make([]datastruct.UserInfo, 0)
	arr := make([]*datastruct.ResponseInviteCount, 0)
	engine.Join("INNER", "invite_info", "invite_info.sended = user_info.id").Find(&users)
	// for _, v := range users {
	// 	resp := new(datastruct.ResponseInviteCount)
	// 	resp.Avatar = v.Avatar
	// 	resp.SpeedFactor =
	// 	arr = append(arr, resp)
	// }
	return arr, datastruct.NULLError
}

func (handle *DBHandler) SetPlayerData(p_data *datastruct.PlayerData) int {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	var userinfo datastruct.UserInfo
	userinfo.IdentityId = p_data.Token
	userinfo.CreatedAt = p_data.CreatedAt
	userinfo.PermissionId = p_data.PermissionId
	userinfo.UpdateTime = p_data.UpdateTime
	userinfo.Avatar = p_data.Avatar
	userinfo.NickName = p_data.NickName

	var tmp datastruct.UserInfo
	var has bool
	var err error

	var referrer datastruct.UserInfo
	if p_data.Id <= 0 {
		if p_data.Referrer > 0 {
			has, _ = session.Id(p_data.Referrer).Get(&referrer)
			if has {
				userinfo.Referrer = p_data.Referrer
			} else {
				userinfo.Referrer = 0
			}
		} else {
			userinfo.Referrer = 0
		}
	}

	if p_data.Id <= 0 {
		_, err = session.Insert(&userinfo)
	} else {
		has, _ = session.Id(p_data.Id).Get(&tmp)
		if has {
			_, err = session.Id(p_data.Id).Update(&userinfo)
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

	if p_data.PermissionId == int(datastruct.Guest) {
		p_data.NickName = tools.GetGuestName(userinfo.Id)
		p_data.Avatar = tools.GetGuestAvatar()
	}

	if p_data.Id <= 0 {
		if p_data.PermissionId == int(datastruct.Player) && userinfo.Referrer > 0 && userinfo.Referrer < userinfo.Id && referrer.PermissionId == int(datastruct.Player) {
			var inviteInfo datastruct.InviteInfo
			inviteInfo.Received = userinfo.Id
			inviteInfo.Sended = userinfo.Referrer
			inviteInfo.CreatedAt = time.Now().Unix()
			_, err := session.Insert(&inviteInfo)
			if err != nil {
				log.Debug("SetPlayerData Insert InviteInfo error:%v", err.Error())
			}
		}
		var rewardStamina datastruct.RewardStamina
		rewardStamina.Id = userinfo.Id
		rewardStamina.GetTime = time.Now().Unix()
		_, err := session.Insert(&rewardStamina)
		if err != nil {
			log.Debug("SetPlayerData Insert RewardStamina error:%v", err.Error())
		}
	}

	sql := fmt.Sprintf("REPLACE INTO player_info (id,honey_count,gold_count,soil_level,stamina,shield,invite_speed_factor)VALUES(%d,%d,%d,%d,%d,%d,%d)", userinfo.Id, p_data.HoneyCount, p_data.GoldCount, p_data.SoilLevel, p_data.Stamina, p_data.Shield, p_data.InviteSpeedFactor)
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->SetPlayerData REPLACE PlayerInfo :%s", err.Error())
		rollback(str, session)
		return userinfo.Id
	}

	if p_data.SpeedUp != nil {
		sql := fmt.Sprintf("REPLACE INTO player_speed_up (id,factor,start,end)VALUES(%d,%d,%d,%d)", userinfo.Id, p_data.SpeedUp.Factor, p_data.SpeedUp.Starting, p_data.SpeedUp.Ending)
		_, err = session.Exec(sql)
		if err != nil {
			str := fmt.Sprintf("DBHandler->SetPlayerData REPLACE player_speed_up :%s", err.Error())
			rollback(str, session)
			return userinfo.Id
		}
	} else {
		var speedup datastruct.PlayerSpeedUp
		has, _ = session.Id(userinfo.Id).Get(&speedup)
		if has {
			_, err = session.Id(userinfo.Id).Delete(&speedup)
			if err != nil {
				str := fmt.Sprintf("DBHandler->SetPlayerData Delete player_speed_up :%s", err.Error())
				rollback(str, session)
				return userinfo.Id
			}
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

func (handle *DBHandler) IsGetStamina(player_id int) bool {
	engine := handle.mysqlEngine
	var rewardStamina datastruct.RewardStamina
	has, _ := engine.Where("id=?", player_id).Get(&rewardStamina)
	if !has {
		rewardStamina.Id = player_id
		rewardStamina.GetTime = time.Now().Unix()
		_, err := engine.Insert(&rewardStamina)
		if err != nil {
			log.Debug("IsGetStamina insert error:%v", err.Error())
		}
	}
	return has
}

func (handle *DBHandler) LotterySteal(player_id int) []*datastruct.UserInfo {
	//compute
	engine := handle.mysqlEngine
	users := make([]*datastruct.UserInfo, 0)
	engine.Where("id <> ?", player_id).Find(&users)
	return users
}

func (handle *DBHandler) GetPlayerDataFromDataBase(user *datastruct.UserInfo) *datastruct.PlayerData {
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
	p_data.InviteSpeedFactor = playerInfo.InviteSpeedFactor

	soil_mp := make(map[int]*datastruct.PlayerSoil)
	petBar_mp := make(map[datastruct.AnimalType]*datastruct.PlayerPetbar)

	var soil_1 datastruct.Soil1
	engine.Where("p_id = ?", p_data.Id).Get(&soil_1)
	soil_mp[1] = tools.CreatePlayerSoil1(&soil_1)

	var soil_2 datastruct.Soil2
	engine.Where("p_id = ?", p_data.Id).Get(&soil_2)
	soil_mp[2] = tools.CreatePlayerSoil2(&soil_2)

	var soil_3 datastruct.Soil3
	engine.Where("p_id = ?", p_data.Id).Get(&soil_3)
	soil_mp[3] = tools.CreatePlayerSoil3(&soil_3)

	var soil_4 datastruct.Soil4
	engine.Where("p_id = ?", p_data.Id).Get(&soil_4)
	soil_mp[4] = tools.CreatePlayerSoil4(&soil_4)

	var soil_5 datastruct.Soil5
	engine.Where("p_id = ?", p_data.Id).Get(&soil_5)
	soil_mp[5] = tools.CreatePlayerSoil5(&soil_5)

	var petbar1 datastruct.Petbar1
	engine.Where("p_id = ?", p_data.Id).Get(&petbar1)
	petBar_mp[datastruct.Sea] = tools.CreatePetbar1(&petbar1)

	var petbar2 datastruct.Petbar2
	engine.Where("p_id = ?", p_data.Id).Get(&petbar2)
	petBar_mp[datastruct.Land] = tools.CreatePetbar2(&petbar2)

	var petbar3 datastruct.Petbar3
	engine.Where("p_id = ?", p_data.Id).Get(&petbar3)
	petBar_mp[datastruct.Fly] = tools.CreatePetbar3(&petbar3)

	var petbar4 datastruct.Petbar4
	engine.Where("p_id = ?", p_data.Id).Get(&petbar4)
	petBar_mp[datastruct.Deity] = tools.CreatePetbar4(&petbar4)

	var playerSpeedUp datastruct.PlayerSpeedUp
	has, _ := engine.Where("id = ?", p_data.Id).Get(&playerSpeedUp)
	if has {
		speedup := new(datastruct.SpeedUpData)
		speedup.Factor = playerSpeedUp.Factor
		speedup.Starting = playerSpeedUp.Start
		speedup.Ending = playerSpeedUp.End
	} else {
		p_data.SpeedUp = nil
	}
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
