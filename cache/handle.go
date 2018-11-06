package cache

import (
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
	   rs = handle.ReadPlayerData(conn,code)
	}
	conn.Close()
	return rs,isExist
}

func (handle *CACHEHandler)GetConn() redis.Conn{
	 conn:=handle.redisClient.Get()
	 return conn
}

func (handle *CACHEHandler) SetPlayerData(p_data *datastruct.PlayerData) {
	conn:=handle.redisClient.Get()
	defer conn.Close()
	key:=p_data.Token
	//add
	_, err := conn.Do("hmset", key,datastruct.GoldField,p_data.GoldCount, datastruct.HoneyField, p_data.HoneyCount, datastruct.IsAuthField, p_data.IsAuth,datastruct.CreatedAtField,p_data.CreatedAt,datastruct.UpdateTimeField,p_data.UpdateTime)
	if err != nil {
		log.Debug("CACHEHandler SetPlayerData err:%s",err.Error())
	}
}

func (handle *CACHEHandler)ReadPlayerData(conn redis.Conn,key string) *datastruct.PlayerData{
	rs := new(datastruct.PlayerData)
	//add
	value, err := redis.Values(conn.Do("hmget",key, datastruct.GoldField, datastruct.HoneyField, datastruct.IsAuthField, datastruct.CreatedAtField,datastruct.UpdateTimeField))
	if err == nil {
	   for index, v := range value {
		   tmp:= v.([]byte)
		   str:= string(tmp[:])
		   switch index{
			 case 0:
				rs.GoldCount = tools.StringToInt64(str)
			 case 1:
				rs.HoneyCount = tools.StringToInt64(str)
			 case 2:
				rs.IsAuth = tools.StringToBool(str)
			 case 3:
				rs.CreatedAt = tools.StringToInt64(str)
			 case 4:
				rs.UpdateTime = tools.StringToInt64(str)
		   }
	   }
	}
	rs.Token = key
	return rs
}



