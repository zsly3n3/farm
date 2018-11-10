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

const PlantLevelField = "PlantLevel"
const SoilLevelField = "SoilLevel"

const PlayerPetbarField = "PlayerPetbar"
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
	PlantLevel int `xorm:"not null INT(11) "`//玩家的种植等级
	SoilLevel int `xorm:"not null INT(11) "`//玩家的土地等级
}

//玩家购买了哪些植物
type PlantForPlayer struct {
	Id int `xorm:"not null pk autoincr INT(11)"` //关联UserInfo中id
	PlayerId int `xorm:"INT(11) not null"`//
	PlantId  int `xorm:"INT(11) not null"`//
}

//植物类型表
type PlantClass struct {
	Id   int       `xorm:"not null pk autoincr INT(11)"` 
	Desc string `xorm:"VARCHAR(32) not null"`//描述
}


//植物表
type Plant struct {
	Id int    `xorm:"not null pk autoincr INT(11)"`
    N  string `xorm:"VARCHAR(64) not null 'name'"` //植物名称
	P int `xorm:"not null INT(11) 'price'"`//价格
	I int `xorm:"not null INT(11) 'in_come'"`//初始收益
	E int `xorm:"not null INT(11) 'exp_for_animal'"`//增加动物经验
	C int `xorm:"not null INT(11) 'class_id'"` //关联PlantClass中id
	L int `xorm:"not null INT(11) 'level'"` //要求玩家种植等级
}

//动物表
type Animal struct {
    Id int    `xorm:"not null pk autoincr INT(11)"`
    N  string `xorm:"VARCHAR(64) not null 'name'"` //名称
	F int `xorm:"not null INT(11) 'factor'"`//增益系数
	E int `xorm:"not null INT(11) 'exp'"`//升级所需经验
	C int `xorm:"not null INT(11) 'class_id'"` //关联AnimalClass中id
	L int `xorm:"not null INT(11) 'level'"` //要求玩家饲养等级
}

//动物类型表
type AnimalClass struct {
	Id   int       `xorm:"not null pk autoincr INT(11)"` 
	Desc string `xorm:"VARCHAR(32) not null"`//描述
}

type ShopData struct{
	Plants []Plant
	Animals []Animal
}

type UserLogin struct{
	 PlatformId Platform //平台
	 Code string //身份标识
	 IsAuth int //是否授权
	 NickName string
	 Avatar string
}

//save reids,save mysql 
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
	SoilLevel int //可购买土地或宠物栏的等级
	Soil []PlayerSoil //玩家土地信息
	PetBar []PlayerPetbar //宠物栏信息
	OwnPlants []int //已购买的植物ID
}

type SoilData struct{
	Id int //土地id 
	Level int //土地当前等级
	Price int //当前价格
	Factor int //生产系数
	Require int //开启条件
}

type GoodsState int 
const (
	Locked GoodsState = iota //有锁
	Unlocked//解锁未购买
	Owned//已拥有
)

//土地表1,2,3,4,5
type Soil1 struct{
	 PId int `xorm:"not null pk INT(11)"` //玩家id
	 Level int `xorm:"not null INT(11)"`//土地等级
	 PlantId int `xorm:"not null INT(11)"`//0表示没有种植
	 Price int  `xorm:"not null INT(11)"`//当前价格
	 Factor int `xorm:"not null INT(11)"`//生产系数
	 State GoodsState `xorm:"not null INT(11)"` //土地状态
}

type Soil2 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	Level int `xorm:"not null INT(11)"`//土地等级
	PlantId int `xorm:"not null INT(11)"`//0表示没有种植
	Price int  `xorm:"not null INT(11)"`//当前价格
	Factor int `xorm:"not null INT(11)"`//生产系数
	State GoodsState `xorm:"not null INT(11)"` //土地状态
}

type Soil3 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	Level int `xorm:"not null INT(11)"`//土地等级
	PlantId int `xorm:"not null INT(11)"`//0表示没有种植
	Price int  `xorm:"not null INT(11)"`//当前价格
	Factor int `xorm:"not null INT(11)"`//生产系数
	State GoodsState `xorm:"not null INT(11)"` //土地状态
}

type Soil4 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	Level int `xorm:"not null INT(11)"`//土地等级
	PlantId int `xorm:"not null INT(11)"`//0表示没有种植
	Price int  `xorm:"not null INT(11)"`//当前价格
	Factor int `xorm:"not null INT(11)"`//生产系数
	State GoodsState `xorm:"not null INT(11)"` //土地状态
}

type Soil5 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	Level int `xorm:"not null INT(11)"`//土地等级
	PlantId int `xorm:"not null INT(11)"`//0表示没有种植
	Price int  `xorm:"not null INT(11)"`//当前价格
	Factor int `xorm:"not null INT(11)"`//生产系数
	State GoodsState `xorm:"not null INT(11)"` //土地状态
}

//宠物栏1,2,3,4 海，陆，空，神
type Petbar1 struct{
	 PId int `xorm:"not null pk INT(11)"` //玩家id
	 AnimalId int `xorm:"not null INT(11)"`//0表示没有种植
	 Price int  `xorm:"not null INT(11)"`//当前价格
	 State GoodsState `xorm:"not null INT(11)"` //状态
}

type Petbar2 struct{
	 PId int `xorm:"not null pk INT(11)"` //玩家id
	 AnimalId int `xorm:"not null INT(11)"`//0表示没有种植
	 Price int  `xorm:"not null INT(11)"`//当前价格
	 State GoodsState `xorm:"not null INT(11)"` //状态
}

type Petbar3 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	AnimalId int `xorm:"not null INT(11)"`//0表示没有种植
	Price int  `xorm:"not null INT(11)"`//当前价格
	State GoodsState `xorm:"not null INT(11)"` //状态
}

type Petbar4 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	AnimalId int `xorm:"not null INT(11)"`//0表示没有种植
	Price int  `xorm:"not null INT(11)"`//当前价格
	State GoodsState `xorm:"not null INT(11)"` //状态
}

type PlayerSoil struct{
	SId int //土地id
	Level int //土地等级
	PlantId int //0表示没有种植
	Price int  //当前价格
	Factor int //生产系数
	State GoodsState //土地状态
}


type PetbarData struct{
	Id int //宠物栏id
	Price int  //当前价格
	Require int //开启条件
}

type PlayerPetbar struct{
	Id int //宠物栏id
	AnimalId int //0表示没有养宠物
	Price int  //当前价格
	State GoodsState
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
	mp["Soil"] = p_data.Soil
	mp["Petbar"] = p_data.PetBar
	mp["PlantLevel"] = p_data.PlantLevel
	mp["OwnPlants"] = p_data.OwnPlants
	return mp
}

