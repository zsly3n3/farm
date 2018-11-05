package cache

import (
	//"github.com/gomodule/redigo/redis"
	"farm/datastruct"
	"farm/log"
)

func (handle *CACHEHandler) GetPlayerData(code string) (*datastruct.PlayerData,bool){
	isExist, err := handle.redisClient.Get().Do("hexists", code, datastruct.GoldField)
    if err != nil {
		log.Debug("hexist failed", err.Error())
    } else {
		log.Debug("exist or not:", isExist)
    }
     

	// result, err := redis.Values(handle.redisClient.Get().Do("hgetall", code))
    // if err != nil{
	//   isExist = false
	//   log.Debug("CACHEHandler GetPlayerData err:%s",err.Error())
	// } else {
	//   log.Debug("all keys and values are:")
    //   for _, v := range result {
    //     log.Debug("%s ", v.([]byte))
    //   }	
	// }
	return nil,false
}