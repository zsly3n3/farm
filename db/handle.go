package db

import(
	"farm/datastruct"
)

func (handle *DBHandler) GetPlayerData(code string) (*datastruct.PlayerData,bool){
	 isExist:=false
	 var rs *datastruct.PlayerData
	 user := new(datastruct.UserInfo)
	 has, _ := handle.mysqlEngine.Where("identity_id=?", code).Get(user)
	 if has{
		isExist = true
		//add
	 }
	 return rs,isExist
}