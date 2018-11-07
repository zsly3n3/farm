package datastruct

import (
	"time"
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

const IdField = "Id"
const GoldField = "GoldCount"
const HoneyField = "HoneyCount"
const PermissionIdField = "PermissionId"
const CreatedAtField = "CreatedAt"
const UpdateTimeField = "UpdateTime"
const IdentityIdField = "IdentityId"
const NickNameField = "NickName"
const AvatarField = "Avatar"

type UserInfo struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	IdentityId string   `xorm:"VARCHAR(128) not null"` //标识id
	PermissionId int `xorm:"not null INT(11)"` //权限id
	CreatedAt int64 `xorm:"bigint not null"` //创建用户的时间
	UpdateTime int64 `xorm:"bigint not null"` //最近一次登录的时间
	NickName string `xorm:"VARCHAR(256) not null"` //昵称
	Avatar string `xorm:"VARCHAR(256) not null"`//头像
}

type Permission struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	Name  string   `xorm:"VARCHAR(32) not null"` //权限名称
}

type PlayerInfo struct {
	Id    int       `xorm:"not null pk INT(11)"` //关联UserInfo中id
	HoneyCount int64 `xorm:"bigint not null"`//蜂蜜数量
	GoldCount int64 `xorm:"bigint not null"`//金币数量
}

//植物类型表
type PlantClass struct {
	Id   int       `xorm:"not null pk INT(11)"` 
	Desc string `xorm:"VARCHAR(32) not null"`//描述
}

//植物表
type Plants struct {
	Id    int       `xorm:"not null pk INT(11)"`
    Name  string   `xorm:"VARCHAR(64) not null"` //植物名称
	Price int `xorm:"not null INT(11)"`//价格
	Income int `xorm:"not null INT(11)"`//初始收益
	ExpForAnimal int `xorm:"not null INT(11)"`//增加动物经验
	ClassId int `xorm:"not null INT(11)"` //关联PlantClass中id
}

type PlantData struct{
	Id int //植物id
	C int //类型id
	N string //植物名称
	P int //价格
	I int //初始收益
	E int //增加动物经验
}



type UserLogin struct{
	 PlatformId Platform //平台
	 Code string //身份标识
	 IsAuth int //是否授权
	 NickName string
	 Avatar string
}

type PlayerData struct{
	Id int //对应数据库中userinfo表中的id
	PermissionId int //权限id
	Token string //标识id IdentityId
	CreatedAt int64 //创建用户的时间
	UpdateTime int64 //最近一次登录的时间
	GoldCount int64 //金币数量
	HoneyCount int64 //蜂蜜数量
	NickName string
	Avatar string
}


type PermissionType int //错误码
const (
	Guest PermissionType = 1 +iota //游客
	Player //普通玩家
)

func CreateUser(code string,permissionId int)*PlayerData{
	player:=new(PlayerData)
	timestamp:=time.Now().Unix()
	player.PermissionId = permissionId
	player.CreatedAt = timestamp
	player.UpdateTime = timestamp
	player.Token = code
	player.GoldCount = 0
	player.HoneyCount = 0
	player.NickName = "test1"
	player.Avatar = "avatar"
	return player
}