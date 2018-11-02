package cache

import (
	"github.com/go-redis/redis"
	"farm/log"
)

const DB_IP = "localhost:6379"
const DB_Index = 0 // use default DB
const DB_Pwd = "Zsly3n@s"


type CACHEHandler struct {
	redisClient *redis.Client
}


func CreateCACHEHandler()*CACHEHandler {
	cacheHandler:=new(CACHEHandler)
	client := redis.NewClient(&redis.Options{
		Addr:     DB_IP,
		Password: DB_Pwd,
		DB:       DB_Index,
	})
	_, err := client.Ping().Result()
	errhandle(err)
    cacheHandler.redisClient = client
	return cacheHandler
}


func errhandle(err error){
	if err != nil {
		log.Fatal("cache error is %v", err.Error())
	}
}