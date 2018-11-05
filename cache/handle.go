package cache

import (
	//"github.com/gomodule/redigo/redis"
	"farm/datastruct"
	//"farm/log"
)

func (handle *CACHEHandler) GetPlayerData(code string) (*datastruct.PlayerData,bool){
	isExist:=false
	var rs *datastruct.PlayerData
	ilen, err := handle.redisClient.Get().Do("hlen", code)
    if err == nil || ((ilen.(int64)) > 0){
	   isExist = true
	   //add
    } 
	return rs,isExist
}

func (handle *CACHEHandler) SetPlayerData(p_data *datastruct.PlayerData) {
     
}