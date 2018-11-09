package event

import(
    "time"
    "farm/datastruct"
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
    player.PlantLevel = 1
    soils:=handle.soils
    petbars:=handle.petbars
    player.Soil=make([]datastruct.SoilData,len(soils),len(soils))
    player.PetBar=make([]datastruct.PetbarData,len(petbars),len(petbars))
    copy(player.Soil,soils)
    copy(player.PetBar,petbars)
	return player
}