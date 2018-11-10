package event

import(
    "time"
    "farm/datastruct"
    // "farm/tools"
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

func (handle *EventHandler)createUser(code string,permissionId int)*datastruct.PlayerData{
	player:=new(datastruct.PlayerData)
	timestamp:=time.Now().Unix()
	player.PermissionId = permissionId
	player.CreatedAt = timestamp
	player.UpdateTime = timestamp
	player.Token = code
	player.GoldCount = 0
	player.HoneyCount = 0
	player.NickName = "test1"
	player.Avatar = "avatar"
    player.PlantLevel = 0
    player.SoilLevel = 0
    player.Soil = createSoil(handle.soils)
    player.PetBar = createPetbar(handle.petbars)
    return player
}

func createSoil(soils []datastruct.SoilData)[]datastruct.PlayerSoil{
    rs:=make([]datastruct.PlayerSoil,0,len(soils))
    for i,v := range soils{
        var tmp datastruct.PlayerSoil
        state:=datastruct.Locked
        if i == 0{
          state = datastruct.Unlocked
        }
        tmp.Id = v.Id
        tmp.Factor = v.Factor
        tmp.Level = v.Level
        tmp.PlantId = 0
        tmp.Price = v.Price
        tmp.State = state
        rs = append(rs,tmp)
    }
    return rs
}

func createPetbar(petbars []datastruct.PetbarData)[]datastruct.PlayerPetbar{
    rs:=make([]datastruct.PlayerPetbar,0,len(petbars))
    for _,v := range petbars{
        var tmp datastruct.PlayerPetbar
        tmp.Id = v.Id
        tmp.AnimalId = 0
        tmp.Price = v.Price
        tmp.State = datastruct.Locked
        rs = append(rs,tmp)
    }
    return rs
}
