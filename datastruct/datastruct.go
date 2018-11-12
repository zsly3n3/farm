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

const PlayerPetbarField = "PlayerPetbar"


type UserInfo struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	IdentityId string   `xorm:"VARCHAR(128) not null"` //标识id
	PermissionId int `xorm:"not null INT(11)"` //权限id
	CreatedAt int64 `xorm:"bigint not null"` //创建用户的时间
	UpdateTime int64 `xorm:"bigint not null"` //最近一次离开或者登陆的时间
	NickName string `xorm:"VARCHAR(255) not null"` //昵称
	Avatar string `xorm:"VARCHAR(255) not null"`//头像
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
}



//植物类型表
type PlantClass struct {
	Id   int       `xorm:"not null pk autoincr INT(11)"` 
	Desc string `xorm:"VARCHAR(32) not null"`//描述
}


//植物表
type Plant struct {
	Id int    `xorm:"not null pk autoincr INT(11)" json:"id"` 
    Name  string `xorm:"VARCHAR(64) not null" json:"name"`//植物名称
	Price int `xorm:"not null INT(11)" json:"price"`//价格
	InCome int `xorm:"not null INT(11)" json:"income"`//初始收益
	ExpForAnimal int `xorm:"not null INT(11)" json:"exp"`//增加动物经验
	Classid int `xorm:"not null INT(11)" json:"type"`//关联PlantClass中id
	Level int `xorm:"not null INT(11)" json:"level"`//要求玩家种植等级
}


//动物表
type Animal struct {
    Id int    `xorm:"not null pk autoincr INT(11)"`
    Name  string `xorm:"VARCHAR(64) not null "` //名称
	Factor int `xorm:"not null INT(11) "`//增益系数
	Exp int64 `xorm:"not null bigint"`//升级所需经验
	ClassId int `xorm:"not null INT(11) "` //关联AnimalClass中id
	Number int `xorm:"not null INT(11) "` //动物编号
}

//动物类型表
type AnimalClass struct {
	Id   AnimalType  `xorm:"not null pk INT(11)"` 
	Desc string `xorm:"VARCHAR(32) not null"`//描述
}

type ShopData struct{
	Plants []*ResponsePlant
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
	Soil map[int]PlayerSoil //玩家土地信息
	PetBar map[AnimalType]PlayerPetbar //宠物栏信息
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
	AnimalNumber int `xorm:"not null INT(11)"`//0表示没有种植
	CurrentExp int64 `xorm:"not null bigint"` //当前宠物栏经验
	State GoodsState `xorm:"not null INT(11)"` //状态
}

type Petbar2 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	AnimalNumber int `xorm:"not null INT(11)"`//0表示没有种植
	CurrentExp int64 `xorm:"not null bigint"` //当前宠物栏经验
	State GoodsState `xorm:"not null INT(11)"` //状态
}

type Petbar3 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	AnimalNumber int `xorm:"not null INT(11)"`//0表示没有种植
	CurrentExp int64 `xorm:"not null bigint"` //当前宠物栏经验
	State GoodsState `xorm:"not null INT(11)"` //状态
}

type Petbar4 struct{
	PId int `xorm:"not null pk INT(11)"` //玩家id
	AnimalNumber int `xorm:"not null INT(11)"`//0表示没有种植
	CurrentExp int64 `xorm:"not null bigint"` //当前宠物栏经验
	State GoodsState `xorm:"not null INT(11)"` //状态
}


type SoilData struct{
	Level int //土地当前等级
	Price int //当前价格
	Factor int //生产系数
	Require int //开启条件
}

type PlayerSoilBase struct{
	Level int //土地等级
	Price int  //当前价格
	Factor int //生产系数
	State GoodsState //土地状态
}

type PlayerSoil struct{
	PlayerSoilBase
	PlantId int //0表示没有种植
}

type ResponsePlayerSoilBase struct{
	PlayerSoilBase
	Id int //土地id
}

type ResponsePlayerSoil struct{
	*ResponsePlayerSoilBase
	PlantId int //0表示没有种植
}





type PetbarData struct{
	Price int //单价
	Require int //开启条件
}

type PlayerPetbar struct{
	AnimalNumber int//为0,表示没有养动物
	CurrentExp int64 //当前宠物栏经验
	State GoodsState
}

type ResponsePetbar struct{
	Animal *ResponseAnimal//动物
	*ResponsePetbarBase
}

type ResponsePetbarBase struct{
	Type AnimalType //宠物栏类型
	Price int //单价
	State GoodsState
}

type ResponseAnimal struct{
	Name  string //名称
	Factor int //增益系数
	CurrentExp int64 //当前经验
	Exp int64 //升级所需经验
}

type ResponsePlant struct{
	Plant
	State GoodsState `json:"state"`
}

type PermissionType int //错误码
const (
	Guest PermissionType = 1 +iota //游客
	Player //普通玩家
)

type AnimalType int 
const (
	Sea AnimalType = iota + 1 //海
	Land//陆
	Fly//空
	Deity//神
)

func ResponseLoginData(p_data *PlayerData,petbars map[AnimalType]PetbarData,ani_mp map[AnimalType]map[int]Animal)map[string]interface{}{
	if p_data == nil{
	   return nil
	}  	
	mp:=make(map[string]interface{})
	mp[PermissionIdField] = &(p_data.PermissionId)
	mp["Token"] = &(p_data.Token)
	mp[GoldField] = &(p_data.GoldCount)
	mp[HoneyField] = &(p_data.HoneyCount)
	mp["Soil"] = responsePlayerSoil(p_data)
	mp["Petbar"] = responsePetbarData(p_data,petbars,ani_mp)
	return mp
}

func responsePlayerSoil(p_data *PlayerData)[]interface{}{
	rs:=make([]interface{},0,len(p_data.Soil))
	for k,v:=range p_data.Soil{
		var interface_var interface{}
		resp_base:=new(ResponsePlayerSoilBase)
		resp_base.Id = k
		resp_base.Level = v.Level
		resp_base.Price = v.Price
		resp_base.Factor = v.Factor
		resp_base.State = v.State
        if v.PlantId == 0{
		  interface_var = resp_base
		} else {
			resp:=new(ResponsePlayerSoil)
			resp.ResponsePlayerSoilBase = resp_base
			resp.PlantId = v.PlantId
			interface_var = resp
		}
		rs = append(rs,interface_var)
	}
	return rs
}

func responsePetbarData(p_data *PlayerData,petbars map[AnimalType]PetbarData,ani_mp map[AnimalType]map[int]Animal)[]interface{}{
	 rs:=make([]interface{}, 0,len(p_data.PetBar))
	 for k,v:= range p_data.PetBar{
		var interface_var interface{}
		base:=new(ResponsePetbarBase)
		base.Type = k
		base.Price = petbars[k].Price
		base.State = v.State

		if v.AnimalNumber == 0{
		  interface_var = base
		} else {
		  resp:=new(ResponsePetbar)
		  resp.ResponsePetbarBase = base
		  var tf bool
		  var anis map[int]Animal
		  anis,tf= ani_mp[k]
		  if tf {
			var ani Animal
			ani,tf=anis[v.AnimalNumber]
			if tf{
			   resp.Animal = new(ResponseAnimal)
			   resp.Animal.CurrentExp = v.CurrentExp
			   resp.Animal.Factor = ani.Factor
			   resp.Animal.Exp = ani.Exp
			   resp.Animal.Name = ani.Name
			} else {
			   resp.Animal = nil 
			}
		  } else {	  
			resp.Animal = nil 
		  }
		  interface_var = resp
		}
		rs = append(rs,interface_var)
	 }
	 return rs
}



