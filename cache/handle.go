package cache

import (
	"fmt"
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


	conn.Send("MULTI")
	conn.Send("hmset", key,
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
	
    
	for i,v := range p_data.Soil{
		soiltableName:=fmt.Sprintf("soil%d",i+1)
		value,isError:=tools.PlayerSoilToString(&v)
		if isError{
		   log.Debug("CACHEHandler SetPlayerData PlayerSoilToString err:%s player:%s",soiltableName,key)	
		   return
		}
		conn.Send("hset", soiltableName,key,value)
	}

	for i,v := range p_data.PetBar{
		petbartableName:=fmt.Sprintf("petbar%d",i+1)
		value,isError:=tools.PlayerPetbarToString(&v)
		if isError{
		   log.Debug("CACHEHandler SetPlayerData PlayerPetbarToString err:%s player:%s",petbartableName,key)	
		   return
		}
		conn.Send("hset", petbartableName,key,value)
	}

	_, err := conn.Do("EXEC")
	
	if err != nil {
	  log.Debug("CACHEHandler SetPlayerData err:%s",err.Error())
	}
}



func (handle *CACHEHandler)ReadPlayerData(conn redis.Conn,key string) *datastruct.PlayerData{
	rs := new(datastruct.PlayerData)
	//add
	value, err := redis.Values(conn.Do("hmget",key,
	datastruct.IdField,datastruct.GoldField, datastruct.HoneyField, 
	datastruct.PermissionIdField,datastruct.CreatedAtField,datastruct.UpdateTimeField,
	datastruct.NickNameField,datastruct.AvatarField,
	datastruct.PlantLevelField,datastruct.SoilLevelField))
	if err!=nil{	
	   log.Debug("CACHEHandler ReadPlayerData err:%s ,player:%s",err.Error(),key)
	   return rs
	}
	for i:=0;i<len(value);i++{
		   tmp:= value[i].([]byte)
		   str:= string(tmp[:])
		   switch i{
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

	len_soil:=5
    rs.Soil=make([]datastruct.PlayerSoil,0,len_soil)
	for i:=1;i<=len_soil;i++{
		soiltableName:=fmt.Sprintf("soil%d",i)
		value, err := redis.String(conn.Do("hget",soiltableName,key))
		if err == nil{
			tmp,_:=tools.BytesToPlayerSoil([]byte(value))
			rs.Soil=append(rs.Soil,*tmp)
		}
	}


	len_petbar:=4
    rs.PetBar=make([]datastruct.PlayerPetbar,0,len_petbar)
	for i:=1;i<=len_petbar;i++{
		petbartableName:=fmt.Sprintf("petbar%d",i)
		value, err := redis.String(conn.Do("hget",petbartableName,key))
		if err == nil{
			tmp,_:=tools.BytesToPlayerPetbar([]byte(value))
			rs.PetBar=append(rs.PetBar,*tmp)
		}
	}

	rs.Token = key
	return rs
}

func (handle *CACHEHandler)UpdatePermisson(key string,permissionId int) datastruct.CodeType{
	conn:=handle.GetConn()
	defer conn.Close()
	rep, err := conn.Do("hset", key,datastruct.PermissionIdField,permissionId)
	log.Debug("rep:%v",rep)
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

func (handle *CACHEHandler)GetPlantLevel(key string)(int,datastruct.CodeType){
	conn:=handle.GetConn()
	defer conn.Close()
	value, err := redis.String(conn.Do("hget",key,datastruct.PlantLevelField))
	code:=datastruct.NULLError
	if err != nil && value == "" {
		code = datastruct.GetDataFailed
		log.Debug("CACHEHandler GetPlantLevel err:%s",err.Error())
		return -1,code
	}
	return tools.StringToInt(value),code
}









func (handle *CACHEHandler)TestMoney(key string){
	conn:=handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("hmset", key,datastruct.GoldField,100, datastruct.HoneyField,200)
	if err != nil {
		log.Debug("CACHEHandler TestMoney err:%s",err.Error())
	}
}

