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
	 has, _ := handle.mysqlEngine.Where("identity_id=?", code).Get(user)
	 if has{
		isExist = true
		//add
		log.Debug("DBHandler GetPlayerData true")
	 }
	 return rs,isExist
}

type UserInfo struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	IdentityId string   `xorm:"VARCHAR(128) not null"` //标识id
	IsAuth int8 `xorm:"TINYINT(1) not null"` //是否授权
	CreatedAt int64 `xorm:"bigint not null"` //创建用户的时间
	UpdateTime int64 `xorm:"bigint not null"` //最近一次登录的时间
}

type PlayerInfo struct {
	Id    int       `xorm:"not null pk INT(11)"` //关联UserInfo中id
	HoneyCount int64 `xorm:"bigint not null"`//蜂蜜数量
	GoldCount int64 `xorm:"bigint not null"`//金币数量
}

func (handle *DBHandler) SetPlayerData(p_data *datastruct.PlayerData) {
	engine:=handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	//add
	var userinfo datastruct.UserInfo
	userinfo.IdentityId = p_data.Token
	userinfo.CreatedAt = p_data.CreatedAt
	var isauth int8
	isauth = 0
	if p_data.IsAuth {
	   isauth = 1
	}
	userinfo.IsAuth = isauth
	userinfo.UpdateTime = p_data.UpdateTime
	// p_data.Id = 1
	// userinfo.Id = p_data.Id
	var err error
	if p_data.Id <= 0{
		_, err = session.Insert(&userinfo)  	
	} else {
		_, err = session.Id(p_data.Id).Update(&userinfo)
	}
	if err != nil{
		str:=fmt.Sprintf("DBHandler->SetPlayerData Update UserInfo :%s",err.Error())
		rollback(str,session)
	    return
	}
	sql:=fmt.Sprintf("REPLACE INTO player_info (id,honey_count,gold_count)VALUES(%d,%d,%d)",userinfo.Id,p_data.HoneyCount,p_data.GoldCount)
	_, err=session.Exec(sql)
	if err != nil{
	  str:=fmt.Sprintf("DBHandler->SetPlayerData Update PlayerInfo :%s",err.Error())
	  rollback(str,session)
	  return
	}
	err=session.Commit()
	if err != nil{
	  str:=fmt.Sprintf("DBHandler->SetPlayerData Commit :%s",err.Error())
	  rollback(str,session)	
	}
}

func rollback(err_str string,session *xorm.Session){
	log.Debug("will rollback,err_str:%v",err_str)
	session.Rollback()
}