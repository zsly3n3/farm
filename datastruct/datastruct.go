package datastruct

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
	GetDataFailed//获取数据失败
	PutDataFailed//修改数据失败
	VersionError//客户端与服务器版本不一致
	TokenNull//没有Token或者值为空
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

//plantlevel , soil 保存到redis和mysql

type UserInfo struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	IdentityId string   `xorm:"VARCHAR(128) not null"` //标识id
	PermissionId int `xorm:"not null INT(11)"` //权限id
	CreatedAt int64 `xorm:"bigint not null"` //创建用户的时间
	UpdateTime int64 `xorm:"bigint not null"` //最近一次离开或者登陆的时间
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
	GoldCount  int64 `xorm:"bigint not null"`//金币数量
	PlantLevel int `xorm:"not null INT(11) "`//可购买商店植物的等级
}

//植物类型表
type PlantClass struct {
	Id   int       `xorm:"not null pk autoincr INT(11)"` 
	Desc string `xorm:"VARCHAR(32) not null"`//描述
}

type Soil struct{
	Id   int       `xorm:"not null pk autoincr INT(11)"` 
	Index int `xorm:"VARCHAR(32) not null"`//土地位置索引
	Price int `xorm:"not null INT(11) "`//初始价格
	InCome int `xorm:"not null INT(11) "`//初始收益
	DefaultLevel int `xorm:"not null INT(11) "` //默认等级
 }


//植物表
type Plant struct {
	Id int    `xorm:"not null pk autoincr INT(11)"`
    N  string `xorm:"VARCHAR(64) not null 'name'"` //植物名称
	P int `xorm:"not null INT(11) 'price'"`//价格
	I int `xorm:"not null INT(11) 'in_come'"`//初始收益
	E int `xorm:"not null INT(11) 'exp_for_animal'"`//增加动物经验
	C int `xorm:"not null INT(11) 'class_id'"` //关联PlantClass中id
}

type ShopData struct{
	Plants []Plant
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
	PlantLevel int //可购买商店植物的等级
	Soil []SoilData //玩家土地信息
}

type SoilData struct{
	Index int //土地索引
	Level int //土地等级
	Isbought int//是否购买
	PlantID int //0表示没有种植
	Price int  //当前价格
	Factor int //生产系数
}



type PermissionType int //错误码
const (
	Guest PermissionType = 1 +iota //游客
	Player //普通玩家
)


func ReponseLoginData(p_data *PlayerData)map[string]interface{}{
	if p_data == nil{
	   return nil
	}
	mp:=make(map[string]interface{})
	mp[PermissionIdField] = &(p_data.PermissionId)
	mp["Token"] = &(p_data.Token)
	mp[GoldField] = &(p_data.GoldCount)
	mp[HoneyField] = &(p_data.HoneyCount)
	mp["PlantLevel"] = &(p_data.PlantLevel)
	mp["Soil"] = p_data.Soil
	return mp
}

