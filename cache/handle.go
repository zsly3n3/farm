package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"farm/datastruct"
	"farm/log"
	"farm/tools"
)

func (handle *CACHEHandler) GetPlayerData(code string) (*datastruct.PlayerData,bool){
	isExist:=false
	var rs *datastruct.PlayerData
	conn:=handle.redisClient.Get()
	ilen, err := conn.Do("hlen", code)
    if err == nil && (ilen.(int64)) > 0{
	   isExist = true
	   rs = readPlayerData(conn,code)
	}
	conn.Close()
	return rs,isExist
}

func (handle *CACHEHandler) SetPlayerData(p_data *datastruct.PlayerData) {
	conn:=handle.redisClient.Get()
	defer conn.Close()
	key:=p_data.IdentityId
	//add
	_, err := conn.Do("hmset", key, datastruct.GoldField, p_data.GoldCount, datastruct.HoneyField, p_data.HoneyCount, datastruct.IsAuthField, p_data.IsAuth,datastruct.CreatedAtField,p_data.CreatedAt,datastruct.UpdateTimeField,p_data.UpdateTime)
	if err == nil {
	}
}

func readPlayerData(conn redis.Conn,key string) *datastruct.PlayerData{
	rs := new(datastruct.PlayerData)
	//add
	value, err := redis.Values(conn.Do("hmget",key, datastruct.GoldField, datastruct.HoneyField, datastruct.IsAuthField, datastruct.CreatedAtField,datastruct.UpdateTimeField))
	if err == nil {
	   for index, v := range value {
		   switch index{
			 case 0:
				tmp:= v.([]byte)
				rs.HoneyCount = tools.ByteArrToInt64(&tmp)
				log.Debug(rs.HoneyCount)
			//  case 1:
			// 	rs.HoneyCount = v.(int64)
			//  case 2:
			// 	rs.IsAuth = v.(int64)
			//  case 3:
			// 	rs.CreatedAt = v.(int64)
			//  case 4:
			// 	rs.UpdateTime = v.(int64)
		   }
	   }
	}
	return rs
}

