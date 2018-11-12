package event

import(
    "time"
    "farm/datastruct"
    "farm/tools"
    // "farm/log"
)

func (handle *EventHandler)createTicker(times time.Duration){
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

func (handle *EventHandler)selectTicker(){
    for{
        select {
         case <-handle.ticker.C:
            handle.checkOnlinePlayer()
        }
    }
}

func (handle *EventHandler)checkOnlinePlayer(){

}

func (handle *EventHandler)createUser(code string,permissionId int,nickName string,avatar string)*datastruct.PlayerData{
	player:=new(datastruct.PlayerData)
	timestamp:=time.Now().Unix()
	player.PermissionId = permissionId
	player.CreatedAt = timestamp
	player.UpdateTime = timestamp
	player.Token = code
	player.GoldCount = 0
	player.HoneyCount = 0
	player.NickName = nickName
	player.Avatar = avatar
    player.PlantLevel = 0
    player.Soil = createSoil(handle.soils)
    player.PetBar = createPetbar(handle.petbars)
    return player
}

func createSoil(soils map[int]datastruct.SoilData)map[int]datastruct.PlayerSoil{
    rs:=make(map[int]datastruct.PlayerSoil)
    for k,v := range soils{
        var tmp datastruct.PlayerSoil
        state:=datastruct.Locked
        tmp.PlantId = 0
        if k == 1{
          state = datastruct.Owned
          tmp.PlantId = 1
        } else {
          tmp.PlantId = 0  
        }
        tmp.Factor = v.Factor
        tmp.Level = v.Level
        tmp.Price = v.Price
        tmp.UpgradeLevelPrice = tools.GetUpgradeLevelPriceForSoil(tmp.Level)
        tmp.State = state
        rs[k]=tmp
    }
    return rs
}

func createPetbar(petbars map[datastruct.AnimalType]datastruct.PetbarData)map[datastruct.AnimalType]datastruct.PlayerPetbar{
    rs:=make(map[datastruct.AnimalType]datastruct.PlayerPetbar)
    for k,_ := range petbars{
        var tmp datastruct.PlayerPetbar
        tmp.AnimalNumber = 0
        tmp.State = datastruct.Locked
        tmp.CurrentExp = 0
        rs[k]=tmp
    }
    return rs
}
