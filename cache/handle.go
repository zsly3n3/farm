package cache

import (
	"github.com/gomodule/redigo/redis"
	"farm/datastruct"
	"farm/log"
)

func (handle *CACHEHandler) GetPlayerData(code string) (*datastruct.PlayerData,bool){
	isExist:=true
	result, err := redis.Values(handle.redisClient.Get().Do("hgetall", code))
    if err != nil{
	  isExist = false
	  log.Debug("CACHEHandler GetPlayerData err:%s",err.Error())
	} else {
	  log.Debug("all keys and values are:")
      for _, v := range result {
        log.Debug("%s ", v.([]byte))
      }	
	}
	return nil,isExist
}