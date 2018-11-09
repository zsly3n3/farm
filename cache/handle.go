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


func (handle *CACHEHandler)SetPlayerSomeData(conn redis.Conn,p_data *datastruct.PlayerData) {
	key:=p_data.Token
	//add
	_, err := conn.Do("hmset", key,
	datastruct.PermissionIdField,p_data.PermissionId,
	datastruct.UpdateTimeField,p_data.UpdateTime)
	if err != nil {
	  log.Debug("CACHEHandler SetPlayerData err:%s",err.Error())
	}
}

func (handle *CACHEHandler)SetPlayerAllData(conn redis.Conn,p_data *datastruct.PlayerData) {
	key:=p_data.Token
	//add
	_, err := conn.Do("hmset", key,
	datastruct.IdField,p_data.Id,
	datastruct.GoldField,p_data.GoldCount,
	datastruct.HoneyField,p_data.HoneyCount,
	datastruct.PermissionIdField,p_data.PermissionId,
	datastruct.CreatedAtField,p_data.CreatedAt,
	datastruct.UpdateTimeField,p_data.UpdateTime,
	datastruct.NickNameField,p_data.NickName,
	datastruct.AvatarField,p_data.Avatar,
	datastruct.PlantLevelField,p_data.PlantLevel,
	datastruct.SoilLevelField,p_data.SoilLevel)
	if err != nil {
	  log.Debug("CACHEHandler SetPlayerData err:%s",err.Error())
	}
}

func (handle *CACHEHandler)ReadPlayerData(conn redis.Conn,key string) *datastruct.PlayerData{
	rs := new(datastruct.PlayerData)
	//add
	value, err := redis.Values(conn.Do("hmget",key,
	datastruct.IdField,datastruct.GoldField, datastruct.HoneyField, 
	datastruct.PermissionIdField, datastruct.CreatedAtField,datastruct.UpdateTimeField,
	datastruct.NickNameField,datastruct.AvatarField,
	datastruct.PlantLevelField,datastruct.SoilLevelField))
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
			 case 8:
				rs.PlantLevel = tools.StringToInt(str)
			 case 9:
				rs.SoilLevel = tools.StringToInt(str)
		   }
	   }
	}
	rs.Token = key
	return rs
}

func (handle *CACHEHandler)UpdatePermisson(key string,permissionId int) datastruct.CodeType{
	conn:=handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("hset", key,datastruct.PermissionIdField,permissionId)
	code:=datastruct.NULLError
	if err != nil {
	   code = datastruct.PutDataFailed
	   log.Debug("CACHEHandler UpdatePermisson err:%s",err.Error())
	}
	return code
}

func (handle *CACHEHandler)clearData(){
	conn:=handle.GetConn()
	defer conn.Close()
	conn.Do("flushdb")
}


func (handle *CACHEHandler)TestMoney(key string){
	conn:=handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("hmset", key,datastruct.GoldField,100, datastruct.HoneyField,200)
	if err != nil {
		log.Debug("CACHEHandler TestMoney err:%s",err.Error())
	}
}

