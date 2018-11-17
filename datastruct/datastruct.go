package datastruct

const NULLSTRING = ""
const NULLID = -1

type TestData struct {
	UserName string
	Avatar   string
}

const MaxStamina = 30
const MaxShield = 3

// DBSessionGetError//xorm事务中Get方法执行出错
// DBSessionExecError//xorm事务中Exec方法执行出错
// DBSessionInsertError//xorm事务中Insert方法执行出错
// DBSessionCommitError//xorm事务中Commit方法执行出错
// DBSessionUpdateError//xorm事务中Update方法执行出错

type CodeType int //错误码
const (
	NULLError                             CodeType = iota //无错误
	ParamError                                            //参数错误,数据为空或者类型不对等
	LoginFailed                                           //登录失败,如无此账号或者密码错误等
	JsonParseFailedFromPostBody                           //来自post请求中的Body解析json失败
	GetDataFailed                                         //获取数据失败
	UpdateDataFailed                                      //修改数据失败
	VersionError                                          //客户端与服务器版本不一致
	TokenError                                            //没有Token或者值为空,或者不存在此Token
	JsonParseFailedFromPutBody                            //来自put请求中的Body解析json失败
	GoldIsNotEnoughForPlant                               //购买植物金币不足
	PlantRequireUnlock                                    //植物未到达解锁条件
	GoldIsNotEnoughForSoil                                //购买土地金币不足
	SoilRequireUnlock                                     //土地未到达解锁条件
	ExpIsNotFullForUpgradeAnimal                          //升级动物失败,经验值不满足 value=13
	HoneyCountIsNotEnoughForUpgradeAnimal                 //升级动物失败,蜂蜜不足 value=14
	AddHoneyCD                                            //采集时间未到
)

type RewardType int //奖励类型
const (
	Gold_10k   RewardType = iota //10k金币  0
	Gold_103k                    //103k金币 1
	Energy_UI1                   //能量     2
	Gold_48k                     //48k金币  3
	Gold_16k                     //16k金币  4
	Gog                          //看护狗   5
	Steal                        //偷取     6
	Energy_UI2                   //能量     7
)

type Platform int //平台
const (
	WX_Platform Platform = iota //微信平台
	PC_Platform                 //pc平台
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

const SoilLevelField = "SoilLevel"
const SpeedUpField = "SpeedUp"
const StaminaField = "Stamina"
const ShieldField = "Shield"

const PlayerPetbarField = "PlayerPetbar"

type UserInfo struct {
	Id           int    `xorm:"not null pk autoincr INT(11)"`
	IdentityId   string `xorm:"VARCHAR(128) pk not null"` //标识id
	PermissionId int    `xorm:"not null INT(11)"`         //权限id
	CreatedAt    int64  `xorm:"bigint not null"`          //创建用户的时间
	UpdateTime   int64  `xorm:"bigint not null"`          //最近一次离开或者登陆的时间
	NickName     string `xorm:"VARCHAR(255) not null"`    //昵称
	Avatar       string `xorm:"VARCHAR(255) not null"`    //头像
}

type Permission struct {
	Id   int    `xorm:"not null pk autoincr INT(11)"`
	Name string `xorm:"VARCHAR(32) not null"` //权限名称
}

type PlayerInfo struct {
	Id         int   `xorm:"not null pk INT(11)"` //关联UserInfo中id
	HoneyCount int64 `xorm:"bigint not null"`     //蜂蜜数量
	GoldCount  int64 `xorm:"bigint not null"`     //金币数量
	SoilLevel  int   `xorm:"not null INT(11) "`   //玩家的土地购买等级
	Stamina    int   `xorm:"not null INT(11) "`   //玩家当前体力值
	Shield     int   `xorm:"not null INT(11) "`   //玩家当前的盾牌，可防止被偷取数据，UI上显示为狗
}

type PlayerSpeedUp struct {
	Id       int   `xorm:"not null pk INT(11)"` //关联UserInfo中id
	Factor   int   `xorm:"not null INT(11)"`    //加速系数
	Starting int64 `xorm:"not null bigint"`     //开始时间，时间戳
	Ending   int64 `xorm:"not null bigint"`     //结束时间，时间戳
}

/*用户领取的体力记录*/
type RewardStamina struct {
	Id      int   `xorm:"not null pk INT(11)"` //关联UserInfo中id
	GetTime int64 `xorm:"not null bigint"`     //领取时间，时间戳
}

//植物类型表
type PlantClass struct {
	Id   int    `xorm:"not null pk autoincr INT(11)"`
	Desc string `xorm:"VARCHAR(32) not null"` //描述
}

//植物表
type Plant struct {
	Id           int    `xorm:"not null pk autoincr INT(11)" json:"id"`
	Name         string `xorm:"VARCHAR(64) not null" json:"name"`    //植物名称
	Price        int64  `xorm:"not null INT(11)" json:"price"`       //价格
	InCome       int64  `xorm:"not null INT(11)" json:"income"`      //初始收益
	ExpForAnimal int64  `xorm:"not null bigint" json:"exp"`          //增加动物经验
	ClassId      int    `xorm:"not null INT(11)" json:"type"`        //关联PlantClass中id
	Level        int    `xorm:"not null INT(11)" json:"level"`       //要求玩家种植等级
	CName        string `xorm:"VARCHAR(64) not null" json:"c_name"`  //植物中文名称
	HoneyCount   int64  `xorm:"not null INT(11)" json:"honey_count"` //蜂蜜产量
}

//动物表
type Animal struct {
	Id         int    `xorm:"not null pk autoincr INT(11)"`
	Name       string `xorm:"VARCHAR(64) not null "` //名称
	InCome     int64  `xorm:"not null INT(11) "`     //初始收益
	Exp        int64  `xorm:"not null bigint"`       //升级所需经验
	HoneyCount int64  `xorm:"not null bigint"`       //升级所需的蜂蜜
	ClassId    int    `xorm:"not null INT(11) "`     //关联AnimalClass中id
	Number     int    `xorm:"not null INT(11) "`     //动物编号
}

//动物类型表
type AnimalClass struct {
	Id   AnimalType `xorm:"not null pk INT(11)"`
	Desc string     `xorm:"VARCHAR(32) not null"` //描述
}

type ShopData struct {
	Plants []*ResponsePlant `json:"plants"`
}

//save reids,save mysql
type PlayerData struct {
	Id           int    //对应数据库中userinfo表中的id
	PermissionId int    //权限id
	Token        string //标识id IdentityId
	CreatedAt    int64  //创建用户的时间
	UpdateTime   int64  //最近一次操作的时间
	GoldCount    int64  //金币数量
	HoneyCount   int64  //蜂蜜数量
	Stamina      int    //玩家当前体力值
	Shield       int    //当前盾牌数量
	NickName     string
	Avatar       string
	SoilLevel    int                          //可购买土地的等级
	Soil         map[int]*PlayerSoil          //玩家土地信息
	PetBar       map[AnimalType]*PlayerPetbar //宠物栏信息
	SpeedUp      *SpeedUpData                 //全局加速数据
}

type SpeedUpData struct {
	Factor   int   //加速系数
	Starting int64 //加速开始时间,时间戳
	Ending   int64 //加速结束时间,时间戳
}

type GoodsState int

const (
	Locked   GoodsState = iota //有锁
	Unlocked                   //解锁未购买
	Owned                      //已拥有
)

//土地表1,2,3,4,5
type Soil1 struct {
	PId               int        `xorm:"not null pk INT(11)"` //玩家id
	Level             int        `xorm:"not null INT(11)"`    //土地等级
	PlantId           int        `xorm:"not null INT(11)"`    //0表示没有种植
	UpgradeLevelPrice int64      `xorm:"not null INT(11)"`    //升下一级的价格
	Factor            int        `xorm:"not null INT(11)"`    //生产系数
	State             GoodsState `xorm:"not null INT(11)"`    //土地状态
	PlantLevel        int        `xorm:"not null INT(11)"`    //可购买商店植物的等级
}

type Soil2 struct {
	PId               int        `xorm:"not null pk INT(11)"` //玩家id
	Level             int        `xorm:"not null INT(11)"`    //土地等级
	PlantId           int        `xorm:"not null INT(11)"`    //0表示没有种植
	UpgradeLevelPrice int64      `xorm:"not null INT(11)"`    //升下一级的价格
	Factor            int        `xorm:"not null INT(11)"`    //生产系数
	State             GoodsState `xorm:"not null INT(11)"`    //土地状态
	PlantLevel        int        `xorm:"not null INT(11)"`    //可购买商店植物的等级
}

type Soil3 struct {
	PId               int        `xorm:"not null pk INT(11)"` //玩家id
	Level             int        `xorm:"not null INT(11)"`    //土地等级
	PlantId           int        `xorm:"not null INT(11)"`    //0表示没有种植
	UpgradeLevelPrice int64      `xorm:"not null INT(11)"`    //升下一级的价格
	Factor            int        `xorm:"not null INT(11)"`    //生产系数
	State             GoodsState `xorm:"not null INT(11)"`    //土地状态
	PlantLevel        int        `xorm:"not null INT(11)"`    //可购买商店植物的等级
}

type Soil4 struct {
	PId               int        `xorm:"not null pk INT(11)"` //玩家id
	Level             int        `xorm:"not null INT(11)"`    //土地等级
	PlantId           int        `xorm:"not null INT(11)"`    //0表示没有种植
	UpgradeLevelPrice int64      `xorm:"not null INT(11)"`    //升下一级的价格
	Factor            int        `xorm:"not null INT(11)"`    //生产系数
	State             GoodsState `xorm:"not null INT(11)"`    //土地状态
	PlantLevel        int        `xorm:"not null INT(11)"`    //可购买商店植物的等级
}

type Soil5 struct {
	PId               int        `xorm:"not null pk INT(11)"` //玩家id
	Level             int        `xorm:"not null INT(11)"`    //土地等级
	PlantId           int        `xorm:"not null INT(11)"`    //0表示没有种植
	UpgradeLevelPrice int64      `xorm:"not null INT(11)"`    //升下一级的价格
	Factor            int        `xorm:"not null INT(11)"`    //生产系数
	State             GoodsState `xorm:"not null INT(11)"`    //土地状态
	PlantLevel        int        `xorm:"not null INT(11)"`    //可购买商店植物的等级
}

//宠物栏1,2,3,4 海，陆，空，神
type Petbar1 struct {
	PId          int        `xorm:"not null pk INT(11)"` //玩家id
	AnimalNumber int        `xorm:"not null INT(11)"`    //0表示没有种植
	CurrentExp   int64      `xorm:"not null bigint"`     //当前宠物栏经验
	State        GoodsState `xorm:"not null INT(11)"`    //状态
}

type Petbar2 struct {
	PId          int        `xorm:"not null pk INT(11)"` //玩家id
	AnimalNumber int        `xorm:"not null INT(11)"`    //0表示没有种植
	CurrentExp   int64      `xorm:"not null bigint"`     //当前宠物栏经验
	State        GoodsState `xorm:"not null INT(11)"`    //状态
}

type Petbar3 struct {
	PId          int        `xorm:"not null pk INT(11)"` //玩家id
	AnimalNumber int        `xorm:"not null INT(11)"`    //0表示没有种植
	CurrentExp   int64      `xorm:"not null bigint"`     //当前宠物栏经验
	State        GoodsState `xorm:"not null INT(11)"`    //状态
}

type Petbar4 struct {
	PId          int        `xorm:"not null pk INT(11)"` //玩家id
	AnimalNumber int        `xorm:"not null INT(11)"`    //0表示没有种植
	CurrentExp   int64      `xorm:"not null bigint"`     //当前宠物栏经验
	State        GoodsState `xorm:"not null INT(11)"`    //状态
}

type SoilData struct {
	Level   int   //土地默认等级
	Price   int64 //购买价格
	Factor  int   //生产系数
	Require int   //开启条件
	LastId  int   //上一个土地id
}

type PlayerSoilBase struct {
	Level             int        `json:"level"`             //土地等级
	Price             int64      `json:"price"`             //购买价格
	UpgradeLevelPrice int64      `json:"upgradelevelprice"` //升下一级的价格
	Factor            int        `json:"factor"`            //生产系数
	State             GoodsState `json:"state"`             //土地状态
}

type PlayerSoil struct {
	PlayerSoilBase
	PlantLevel int
	PlantId    int //0表示没有种植
}

type ResponsePlayerSoilBase struct {
	PlayerSoilBase
	Id int `json:"id"` //土地id
}

type ResponsePlayerSoil struct {
	*ResponsePlayerSoilBase
	Plant *ResponseSoilPlant `json:"plant"`
}

type ResponseSoilPlant struct {
	Name         string `json:"name"`
	InCome       int64  `json:"income"`
	ExpForAnimal int64  `json:"expforanimal"`
	Type         int    `json:"type"`
}

type ResponseSoil struct {
	Level  int //土地默认等级
	Factor int //生产系数
}

type PetbarData struct {
	Price   int64 //单价
	Require int   //开启条件
	Id      int   //土地id
	LastId  int   //上一个土地id
}

type PlayerPetbar struct {
	AnimalNumber int   //为0,表示没有养动物
	CurrentExp   int64 //当前宠物栏经验
	State        GoodsState
}

type ResponsePetbar struct {
	Animal *ResponseAnimal `json:"animal"` //动物
	*ResponsePetbarBase
}

type ResponsePetbarBase struct {
	Type  AnimalType `json:"type"`  //宠物栏类型
	Price int64      `json:"price"` //单价
	State GoodsState `json:"state"`
	Id    int        `json:"id"` //宠物栏id
}

type ResponseAnimal struct {
	Name       string `json:"name"`       //名称
	InCome     int64  `json:"income"`     //基本收益
	CurrentExp int64  `json:"currentexp"` //当前经验
	Exp        int64  `json:"exp"`        //升级所需经验
	HoneyCount int64  `json:"honeycount"` //升级所需蜂蜜
	IsLast     int    `json:"islast"`     //是否为最后一个动物 ,1是，0不是
}

type ResponsePlant struct {
	Plant
	State GoodsState `json:"state"`
}

type PermissionType int //错误码
const (
	Guest  PermissionType = 1 + iota //游客
	Player                           //普通玩家
)

type AnimalType int

const (
	Sea   AnimalType = iota + 1 //海
	Land                        //陆
	Fly                         //空
	Deity                       //神
)

func ResponseLoginData(tmp *TmpLoginData, p_data *PlayerData, plants []Plant, petbars map[AnimalType]PetbarData, ani_mp map[AnimalType]map[int]Animal) map[string]interface{} {
	if p_data == nil {
		return nil
	}
	mp := make(map[string]interface{})
	farm_mp := make(map[string]interface{})
	mp["farm"] = farm_mp
	mp["nickname"] = &p_data.NickName
	mp["avatar"] = &p_data.Avatar
	mp["permissionid"] = &(p_data.PermissionId)
	mp["token"] = &(p_data.Token)
	mp["lottery"] = responesLottery(p_data)
	farm_mp["goldcount"] = &(p_data.GoldCount)
	farm_mp["honeycount"] = &(p_data.HoneyCount)
	farm_mp["dogs"] = &(p_data.Shield)
	farm_mp["soil"] = GetResponsePlayerSoil(p_data, plants)
	farm_mp["petbar"] = GetResponsePetbarData(p_data, petbars, ani_mp)
	if tmp == nil {
		farm_mp["speedcd"] = 0
	} else {
		farm_mp["speedcd"] = tmp.CD
		if p_data.SpeedUp != nil {
			resp_speed := new(ResponesSpeedUpData)
			resp_speed.Factor = p_data.SpeedUp.Factor
			resp_speed.Ending = tmp.Sec_EndingSpeedUp
			farm_mp["speedup"] = resp_speed
		}
	}

	return mp
}

func GetResponsePlayerSoil(p_data *PlayerData, plants []Plant) []interface{} {
	length := len(p_data.Soil)
	start_index := 1
	rs := make([]interface{}, length, length)
	for k, v := range p_data.Soil {
		var interface_var interface{}
		resp_base := new(ResponsePlayerSoilBase)
		resp_base.Id = k
		resp_base.Level = v.Level
		resp_base.Price = v.Price
		resp_base.UpgradeLevelPrice = v.UpgradeLevelPrice
		resp_base.Factor = v.Factor
		resp_base.State = v.State
		if v.PlantId <= 0 {
			interface_var = resp_base
		} else {
			resp := new(ResponsePlayerSoil)
			resp.ResponsePlayerSoilBase = resp_base
			resp.Plant = createResponseSoilPlant(v.PlantId, plants)
			interface_var = resp
		}
		rs[resp_base.Id-start_index] = interface_var
	}

	return rs
}

func createResponseSoilPlant(plant_id int, plants []Plant) *ResponseSoilPlant {
	rs := new(ResponseSoilPlant)
	for _, v := range plants {
		if v.Id == plant_id {
			rs.Name = v.Name
			rs.InCome = v.InCome
			rs.ExpForAnimal = v.ExpForAnimal
			rs.Type = v.ClassId
			break
		}
	}
	return rs
}

func GetResponsePetbarData(p_data *PlayerData, petbars map[AnimalType]PetbarData, ani_mp map[AnimalType]map[int]Animal) []interface{} {
	length := len(p_data.PetBar)
	rs := make([]interface{}, length, length)
	start_index := 6
	for k, v := range p_data.PetBar {
		var interface_var interface{}
		base := new(ResponsePetbarBase)
		base.Type = k
		base.Price = petbars[k].Price
		base.State = v.State
		base.Id = petbars[k].Id
		if v.AnimalNumber <= 0 {
			interface_var = base
		} else {
			resp := new(ResponsePetbar)
			resp.ResponsePetbarBase = base
			var tf bool
			var anis map[int]Animal
			anis, tf = ani_mp[k]
			if tf {
				var ani Animal
				ani, tf = anis[v.AnimalNumber]
				num := len(anis)
				if tf {
					resp.Animal = new(ResponseAnimal)
					resp.Animal.CurrentExp = v.CurrentExp
					resp.Animal.InCome = ani.InCome
					resp.Animal.Exp = ani.Exp
					resp.Animal.Name = ani.Name
					resp.Animal.HoneyCount = ani.HoneyCount
					isLast := 0
					if num == ani.Number {
						isLast = 1
					}
					resp.Animal.IsLast = isLast
				} else {
					resp.Animal = nil
				}
			} else {
				resp.Animal = nil
			}
			interface_var = resp
		}
		rs[base.Id-start_index] = interface_var
	}
	return rs
}

func responesLottery(p_data *PlayerData) *ResponesLottery {
	lottery := new(ResponesLottery)
	lottery.MaxStamina = MaxStamina
	lottery.CurrentStamina = p_data.Stamina
	return lottery
}

//Tmp数据，用于传递，不保存redis和mysql
type TmpLoginData struct {
	Sec_EndingSpeedUp int64 //还剩多少秒结束加速
	CD                int64 `json:"time"` //还有多少时间，才能使用加速
}

//response
type ResponseUpgradeSoil struct {
	GoldCount    int64 `json:"goldcount"`
	Level        int   `json:"level"`
	Factor       int   `json:"factor"`
	UpgradePrice int64 `json:"upgradeprice"`
}

type ResponseAnimalUpgrade struct {
	HoneyCount int64
	Animal     *ResponseAnimal
	RightExp   int64
}

type ResponseAddHoney struct {
	HoneyCount int64                `json:"addhoney"`
	CD         int64                `json:"nextspeedcd"`
	SpeedUp    *ResponesSpeedUpData `json:"speedup"`
}

type ResponesSpeedUpData struct {
	Factor int   `json:"factor"` //加速系数
	Ending int64 `json:"ending"` //加速结束时间,多少秒之后结束
}

type ResponesLottery struct {
	MaxStamina     int `json:"maxstamina"`     //体力值
	CurrentStamina int `json:"currentstamina"` //体力值
}

type ResponesStaminaData struct {
	Stamina     int   `json:"currentstamina"` //当前体力值
	NextRequest int64 `json:"nextrequest"`    //下一次请求时间，单位是秒
}

type ResponesLotteryData struct {
	Stamina   int             `json:"currentstamina"` //当前体力值
	GoldCount int64           `json:"goldcount"`
	Shield    int             `json:"dogs"`
	Stolen    *ResponseStolen `json:"stolen"`
}

type ResponseStolen struct {
	StolenGold  int64                  `json:"stolengold"`
	StolenHoney int64                  `json:"stolenhoney"`
	Succeed     int                    `json:"succeed"`
	PlayerData  map[string]interface{} `json:"playerdata"`
}

//body
type UserLogin struct {
	PlatformId Platform `json:"platformid"` //平台
	Code       string   `json:"code"`       //身份标识
	IsAuth     int      `json:"isauth"`     //是否授权
	NickName   string   `json:"nickname"`
	Avatar     string   `json:"avatar"`
}

type PlantInSoil struct {
	PlantId int `json:"plantid"`
	SoilId  int `json:"soilid"`
}

type UpgradeSoil struct {
	Level  int `json:"level"`
	SoilId int `json:"soilid"`
}

type BuyPetbar struct {
	PetbarId int `json:"petbarid"`
}

type AddExpForAnimal struct {
	PetbarId int `json:"petbarid"`
	SoilId   int `json:"soilid"`
}

type LotteryBody struct {
	RewardType int `json:"rewardtype"`
	Expend     int `json:"expend"`
}

type UserAuthBody struct {
	NickName string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
