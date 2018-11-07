package db

import(
	"github.com/go-xorm/xorm"
	"farm/datastruct"
	"farm/log"
	"fmt"
)

func (handle *DBHandler) GetPlayerData(code string) (*datastruct.PlayerData,bool){
	 isExist:=false
	 var rs *datastruct.PlayerData
	 user := new(datastruct.UserInfo)
	 engine:=handle.mysqlEngine
	 has, _ := engine.Where("identity_id=?", code).Get(user)
	 if has{
		isExist = true
		//add
		rs = new(datastruct.PlayerData)
		rs.Avatar = user.Avatar
		rs.CreatedAt = user.CreatedAt
		rs.Id = user.Id
		rs.NickName = user.NickName
		rs.PermissionId = user.PermissionId
		rs.Token = user.IdentityId
		rs.UpdateTime = user.UpdateTime
		var playerInfo datastruct.PlayerInfo
		engine.Id(rs.Id).Get(&playerInfo)
		rs.GoldCount = playerInfo.GoldCount
		rs.HoneyCount = playerInfo.HoneyCount
	 }
	 return rs,isExist
}

func (handle *DBHandler) SetPlayerData(p_data *datastruct.PlayerData) int {
	engine:=handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	//add
	var userinfo datastruct.UserInfo
	userinfo.IdentityId = p_data.Token
	userinfo.CreatedAt = p_data.CreatedAt
	userinfo.PermissionId = p_data.PermissionId
	userinfo.UpdateTime = p_data.UpdateTime
	userinfo.Avatar = p_data.Avatar
	userinfo.NickName = p_data.NickName
	var err error
	if p_data.Id <= 0{
		_, err = session.Insert(&userinfo)  	
	} else {
		var tmp datastruct.UserInfo
		var has bool
		has, err = session.Where("id=?",p_data.Id).Get(&tmp)
	    if has {
		  userinfo.Id = p_data.Id
		  _, err = session.Where("id=?",p_data.Id).Update(&userinfo)
		} else {
		  _, err = session.Insert(&userinfo)
		}
	}
	if err != nil{
		str:=fmt.Sprintf("DBHandler->SetPlayerData Update UserInfo :%s",err.Error())
		rollback(str,session)
	    return userinfo.Id
	}
	sql:=fmt.Sprintf("REPLACE INTO player_info (id,honey_count,gold_count)VALUES(%d,%d,%d)",userinfo.Id,p_data.HoneyCount,p_data.GoldCount)
	_, err=session.Exec(sql)
	if err != nil{
	  str:=fmt.Sprintf("DBHandler->SetPlayerData Update PlayerInfo :%s",err.Error())
	  rollback(str,session)
	  return userinfo.Id
	}
	err=session.Commit()
	if err != nil{
	  str:=fmt.Sprintf("DBHandler->SetPlayerData Commit :%s",err.Error())
	  rollback(str,session)	
	}
	return userinfo.Id
}

func rollback(err_str string,session *xorm.Session){
	log.Debug("will rollback,err_str:%v",err_str)
	session.Rollback()
}


func(handle *DBHandler)GetPlantsData()(datastruct.CodeType,[]datastruct.Plant){
	engine:=handle.mysqlEngine
	plants := make([]datastruct.Plant, 0)
	err := engine.Find(&plants)
	if err != nil{
	   log.Debug("GetPlantsData error:%v",err.Error())
	   return datastruct.GetDataFailed,nil
	}
    return datastruct.NULLError,plants
}