package datastruct

import (
)


const NULLSTRING = ""
const NULLID = -1

type TestData struct {
	 UserName string
	 Avatar string
}



// DBSessionGetError//xorm事务中Get方法执行出错
// DBSessionExecError//xorm事务中Exec方法执行出错
// DBSessionInsertError//xorm事务中Insert方法执行出错
// DBSessionCommitError//xorm事务中Commit方法执行出错
// DBSessionUpdateError//xorm事务中Update方法执行出错


type CodeType int //错误码
const (
	NULLError CodeType = iota //无错误
	ParamError//参数错误,数据为空或者类型不对等
	LoginFailed//登录失败,如无此账号或者密码错误等
	JsonParseFailedFromPostBody//来自post请求中的Body解析json失败
)

type Platform int //平台
const (
	WX_Platform Platform = iota //微信平台
    PC_Platform //pc平台
)

const GoldField = "GoldCount"
const HoneyField = "HoneyCount"

type UserInfo struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	IdentityId string   `xorm:"VARCHAR(128) not null"` //标识id
	IsAuth int8 `xorm:"TINYINT(1) not null"` //是否授权
	CreatedAt int64 `xorm:"bigint not null"` //创建用户的时间
	UpdateTime int64 `xorm:"bigint not null"` //最近一次登录的时间
}

type PlayerInfo struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	IdentityId string   `xorm:"VARCHAR(128) not null"` //标识id
	HoneyCount int64 `xorm:"bigint not null"`//蜂蜜数量
	GoldCount int64 `xorm:"bigint not null"`//金币数量
}


type UserLogin struct{
	 PlatformId Platform //平台
	 Code string //身份标识
}


type PlayerData struct{
	GoldCount int64 //金币数量
	HoneyCount int64 //蜂蜜数量
}






