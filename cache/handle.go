package cache

import (
	"github.com/gomodule/redigo/redis"
	"farm/datastruct"
	"farm/log"
)

func (handle *CACHEHandler) GetPlayerData(code string) (*datastruct.PlayerData,bool){
	isExist:=false
	var rs *datastruct.PlayerData
	ilen, err := handle.redisClient.Get().Do("hlen", code)
    if err == nil && (ilen.(int64)) > 0{
	   isExist = true
	   //add
    } 
	return rs,isExist
}

func (handle *CACHEHandler) SetPlayerData(p_data *datastruct.PlayerData) {
	conn:=handle.redisClient.Get()
	key:=p_data.IdentityId
	_, err := conn.Do("hmset", key, datastruct.GoldField, p_data.GoldCount, datastruct.HoneyField, p_data.HoneyCount, datastruct.IsAuthField, p_data.IsAuth,datastruct.CreatedAtField,p_data.CreatedAt,datastruct.UpdateTimeField,p_data.UpdateTime)
	if err == nil {
	   value, err := redis.Values(conn.Do("hmget",key, datastruct.GoldField, datastruct.HoneyField, datastruct.IsAuthField, datastruct.CreatedAtField,datastruct.UpdateTimeField))
	   if err == nil {
		for _, v := range value {
			log.Debug("%s ", v.([]byte))
		}
	   }
	}
}