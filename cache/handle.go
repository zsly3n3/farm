package cache

import (
	"github.com/gomodule/redigo/redis"
	"farm/datastruct"
	"farm/log"
	"farm/tools"
)

func (handle *CACHEHandler) GetPlayerData(conn redis.Conn,code string) (*datastruct.PlayerData,bool){
	isExist:=false
	var rs *datastruct.PlayerData
	ilen, err := conn.Do("hlen", code)
    if err == nil && (ilen.(int64)) > 0{
	   isExist = true
	   rs = handle.ReadPlayerData(conn,code)
	}
	return rs,isExist
}

func (handle *CACHEHandler)GetConn() redis.Conn{
	 conn:=handle.redisClient.Get()
	 return conn
}

func (handle *CACHEHandler) SetPlayerID(conn redis.Conn,key string,p_id int){
	_, err := conn.Do("hset", key,datastruct.IdField,p_id)
	if err != nil {
		log.Debug("CACHEHandler SetPlayerID err:%s",err.Error())
	}
}

func (handle *CACHEHandler)SetPlayerData(conn redis.Conn,args ...interface{}) {
	if len(args) <= 1{
		log.Debug("CACHEHandler SetPlayerData args error")
		return
	}
	for i,v := range args{
	   index:=i+1
	   switch v.(type){
		 case string:
			log.Debug("第%d : %v",index,v.(string))
		 case int64:
			log.Debug("第%d : %v",index,v.(int64))
		 case int:
			log.Debug("第%d : %v",index,v.(int))		  
	   }
	}
	_, err := conn.Do("hmset", args)
	if err != nil {
	  log.Debug("CACHEHandler SetPlayerData err:%s",err.Error())
	}
}

// func (handle *CACHEHandler)SetPlayerData(conn redis.Conn,p_data *datastruct.PlayerData) {
// 	key:=p_data.Token
// 	//add
// 	_, err := conn.Do("hmset", key,
// 	datastruct.IdField,p_data.Id,
// 	datastruct.GoldField,p_data.GoldCount,
// 	datastruct.HoneyField,p_data.HoneyCount,
// 	datastruct.PermissionIdField,p_data.PermissionId,
// 	datastruct.CreatedAtField,p_data.CreatedAt,
// 	datastruct.UpdateTimeField,p_data.UpdateTime,
// 	datastruct.NickNameField,p_data.NickName,
// 	datastruct.AvatarField,p_data.Avatar)
// 	if err != nil {
// 	  log.Debug("CACHEHandler SetPlayerData err:%s",err.Error())
// 	}
// }

func (handle *CACHEHandler)ReadPlayerData(conn redis.Conn,key string) *datastruct.PlayerData{
	rs := new(datastruct.PlayerData)
	//add
	value, err := redis.Values(conn.Do("hmget",key,
	datastruct.IdField,datastruct.GoldField, datastruct.HoneyField, 
	datastruct.PermissionIdField, datastruct.CreatedAtField,datastruct.UpdateTimeField,
	datastruct.NickNameField,datastruct.AvatarField))
	if err == nil {
	   for index, v := range value {
		   tmp:= v.([]byte)
		   str:= string(tmp[:])
		   switch index{
			 case 0:
				rs.Id = tools.StringToInt(str)
			 case 1:
				rs.GoldCount = tools.StringToInt64(str)
			 case 2:
				rs.HoneyCount = tools.StringToInt64(str)
			 case 3:
				rs.PermissionId = tools.StringToInt(str)
			 case 4:
				rs.CreatedAt = tools.StringToInt64(str)
			 case 5:
				rs.UpdateTime = tools.StringToInt64(str)
			 case 6:
				rs.NickName = str
			 case 7:
				rs.Avatar = str
		   }
	   }
	}
	rs.Token = key
	return rs
}

func (handle *CACHEHandler)TestMoney(key string){
	conn:=handle.GetConn()
	_, err := conn.Do("hmset", key,datastruct.GoldField,100, datastruct.HoneyField,200)
	if err != nil {
		log.Debug("CACHEHandler TestMoney err:%s",err.Error())
	}
}

