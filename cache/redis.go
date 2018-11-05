package cache

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"farm/log"
)

const DB_IP = "localhost:6379"
const DB_NAME = "farm"
const DB_Pwd = "Zsly3n@s"


type CACHEHandler struct {
	redisClient *redis.Pool
}

func getRedisPool() *redis.Pool{
    RedisClient := &redis.Pool{
        // 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
        MaxIdle:     1,//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
        MaxActive:   10,//最大的激活连接数，表示同时最多有N个连接
        IdleTimeout: 180 * time.Second,
        Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",DB_IP)
			errhandle(err)
			log.Debug("----aaa---")
			// if _, err := conn.Do("AUTH", DB_Pwd); err != nil {		
			//  conn.Close()
			//  errhandle(err)
			// }	
            // c.Do("SELECT",DB_NAME)
            return c, nil
        },
    }
    return RedisClient
}

func CreateCACHEHandler()*CACHEHandler {
	cacheHandler:=new(CACHEHandler)
	cacheHandler.redisClient = getRedisPool()
	return cacheHandler
}


func errhandle(err error){
	if err != nil {
		log.Fatal("cache error is %v", err.Error())
	}
}