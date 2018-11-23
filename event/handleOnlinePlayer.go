package event

import (
	"farm/datastruct"
	"farm/tools"
	"time"
)

func (handle *EventHandler) createTicker(times time.Duration) {
	if !handle.isExistTicker {
		handle.isExistTicker = true
		handle.ticker = time.NewTicker(times)
		go handle.selectTicker()
	}
}

// func (handle *EventHandler)stopTicker(){
//     if handle.ticker != nil{
// 	   handle.ticker.Stop()
// 	   handle.isExistTicker = false
//     }
// }

func (handle *EventHandler) selectTicker() {
	for {
		select {
		case <-handle.ticker.C:
			handle.checkOnlinePlayer()
		}
	}
}

func (handle *EventHandler) checkOnlinePlayer() {
	currentTime := time.Now().Unix()
	slice := make([]interface{}, 0)
	handle.onlinePlayers.Lock.Lock()
	defer handle.onlinePlayers.Lock.Unlock()
	for k, v := range handle.onlinePlayers.Bm {
		if v.WillDelete {
			slice = append(slice, k)
		} else {
			if v.OnlineTime < currentTime {
				onlinePlayerData := new(datastruct.OnlinePlayerData)
				onlinePlayerData.WillDelete = true
				handle.onlinePlayers.Bm[k] = *onlinePlayerData
				go handle.fromRedisToMysql(k)
			}
		}
	}
	if len(slice) > 0 {
		conn := handle.cacheHandler.GetConn()
		defer conn.Close()
		for _, v := range slice {
			key := v.(string)
			delete(handle.onlinePlayers.Bm, key)
			isRemoved, userId := handle.cacheHandler.IsRemoveGuest(key)
			if isRemoved {
				go handle.dbHandler.DeleteUser(userId, handle.soils, handle.petbars) //游客过期,删除数据
			}
		}
		handle.deletefromRedis(slice)
	}
}
func (handle *EventHandler) deletefromRedis(keys []interface{}) {
	handle.cacheHandler.DeletedKeys(keys, handle.petbars, handle.soils)
}

func (handle *EventHandler) createUser(code string, permissionId int, nickName string, avatar string) *datastruct.PlayerData {
	player := new(datastruct.PlayerData)
	timestamp := time.Now().Unix()
	player.PermissionId = permissionId
	player.CreatedAt = timestamp
	player.UpdateTime = timestamp
	player.Token = code
	player.GoldCount = 0
	player.HoneyCount = 0
	player.NickName = nickName
	player.Avatar = avatar
	player.SoilLevel = 0
	player.Soil = createSoil(handle.soils)
	player.PetBar = createPetbar(handle.petbars)
	player.SpeedUp = nil
	player.Stamina = datastruct.MaxStamina
	player.Shield = 0
	return player
}

func createSoil(soils map[int]datastruct.SoilData) map[int]*datastruct.PlayerSoil {
	rs := make(map[int]*datastruct.PlayerSoil)
	for k, v := range soils {
		tmp := new(datastruct.PlayerSoil)
		state := datastruct.Locked
		tmp.PlantId = 0
		if k == 1 {
			state = datastruct.Unlocked
		}
		tmp.Factor = v.Factor
		tmp.Level = v.Level
		tmp.UpgradeLevelPrice = tools.GetUpgradeLevelPriceForSoil(tmp.Level)
		tmp.State = state
		tmp.PlantLevel = 0
		rs[k] = tmp
	}
	return rs
}

func createPetbar(petbars map[datastruct.AnimalType]datastruct.PetbarData) map[datastruct.AnimalType]*datastruct.PlayerPetbar {
	rs := make(map[datastruct.AnimalType]*datastruct.PlayerPetbar)
	for k, _ := range petbars {
		tmp := new(datastruct.PlayerPetbar)
		tmp.AnimalNumber = 0
		tmp.State = datastruct.Locked
		tmp.CurrentExp = 0
		rs[k] = tmp
	}
	return rs
}
