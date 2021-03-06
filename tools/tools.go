package tools

import (
	"crypto/md5"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"farm/conf"
	"farm/datastruct"
	"farm/log"
	"fmt"
	"io"
	"math/rand"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func Int64ToString(tmp int64) string {
	return strconv.FormatInt(tmp, 10)
}

func StringToInt64(tmp string) int64 {
	rs, _ := strconv.ParseInt(tmp, 10, 64)
	return rs
}

func IntToString(tmp int) string {
	return strconv.Itoa(tmp)
}

func StringToInt(tmp string) int {
	rs, _ := strconv.Atoi(tmp)
	return rs
}

func StringToFloat64(tmp string) float64 {
	rs, _ := strconv.ParseFloat(tmp, 64)
	return rs
}

func BoolToString(tmp bool) string {
	if tmp == false {
		return "0"
	} else {
		return "1"
	}
}

func StringToBool(tmp string) bool {
	if tmp == "0" {
		return false
	} else {
		return true
	}
}

func GetPlantsInfo() []datastruct.Plant {
	xlsx, err := excelize.OpenFile("conf/shop_data.xlsx")
	if err != nil {
		log.Fatal("Excel error is %v", err.Error())
	}
	index := 2
	tableName := "Sheet1"
	plants := make([]datastruct.Plant, 0)
	for {
		cell_Name := fmt.Sprintf("A%d", index)
		cell_ClassId := fmt.Sprintf("B%d", index)
		cell_Price := fmt.Sprintf("C%d", index)
		cell_Income := fmt.Sprintf("D%d", index)
		cell_AddExp := fmt.Sprintf("E%d", index)
		cell_Level := fmt.Sprintf("F%d", index)
		cell_ChName := fmt.Sprintf("G%d", index)
		cell_honey := fmt.Sprintf("H%d", index)
		name := xlsx.GetCellValue(tableName, cell_Name)
		cid := xlsx.GetCellValue(tableName, cell_ClassId)
		price := xlsx.GetCellValue(tableName, cell_Price)
		income := xlsx.GetCellValue(tableName, cell_Income)
		exp := xlsx.GetCellValue(tableName, cell_AddExp)
		level := xlsx.GetCellValue(tableName, cell_Level)
		chName := xlsx.GetCellValue(tableName, cell_ChName)
		honeyCount := xlsx.GetCellValue(tableName, cell_honey)
		if name == "" {
			break
		}
		var plant datastruct.Plant
		plant.Name = name
		plant.ClassId = StringToInt(cid)
		plant.Price = StringToInt64(price)
		plant.InCome = StringToInt64(income)
		plant.ExpForAnimal = StringToInt64(exp)
		plant.Level = StringToInt(level)
		plant.HoneyCount = StringToInt64(honeyCount)
		plant.CName = chName
		plants = append(plants, plant)
		index++
	}
	return plants
}

func GetAnimalInfo() []datastruct.Animal {
	xlsx, err := excelize.OpenFile("conf/shop_data.xlsx")
	if err != nil {
		log.Fatal("Excel error is %v", err.Error())
	}
	index := 2
	tableName := "Sheet2"
	animals := make([]datastruct.Animal, 0)
	for {
		cell_Name := fmt.Sprintf("A%d", index)
		cell_ClassId := fmt.Sprintf("B%d", index)
		cell_InCome := fmt.Sprintf("C%d", index)
		cell_Exp := fmt.Sprintf("D%d", index)
		cell_Number := fmt.Sprintf("E%d", index)
		cell_Honey := fmt.Sprintf("F%d", index)
		cell_CName := fmt.Sprintf("G%d", index)
		name := xlsx.GetCellValue(tableName, cell_Name)
		cid := xlsx.GetCellValue(tableName, cell_ClassId)
		income := xlsx.GetCellValue(tableName, cell_InCome)
		exp := xlsx.GetCellValue(tableName, cell_Exp)
		number := xlsx.GetCellValue(tableName, cell_Number)
		honey := xlsx.GetCellValue(tableName, cell_Honey)
		cname := xlsx.GetCellValue(tableName, cell_CName)
		if name == "" {
			break
		}
		var animal datastruct.Animal
		animal.Name = name
		animal.ClassId = StringToInt(cid)
		animal.InCome = StringToInt64(income)
		animal.Exp = StringToInt64(exp)
		animal.Number = StringToInt(number)
		animal.HoneyCount = StringToInt64(honey)
		animal.CName = cname
		animals = append(animals, animal)
		index++
	}
	return animals
}

func GetSoildInfo() (map[int]datastruct.SoilData, map[datastruct.AnimalType]datastruct.PetbarData) {
	xlsx, err := excelize.OpenFile("conf/soil_data.xlsx")
	if err != nil {
		log.Fatal("Excel error is %v", err.Error())
	}
	index := 2
	soildtTableName := "Sheet1"
	soils := make(map[int]datastruct.SoilData)
	for {
		cell_index := fmt.Sprintf("A%d", index)
		cell_price := fmt.Sprintf("B%d", index)
		cell_factor := fmt.Sprintf("C%d", index)
		cell_require := fmt.Sprintf("D%d", index)
		cell_lastId := fmt.Sprintf("E%d", index)
		location := xlsx.GetCellValue(soildtTableName, cell_index)
		price := xlsx.GetCellValue(soildtTableName, cell_price)
		factor := xlsx.GetCellValue(soildtTableName, cell_factor)
		require := xlsx.GetCellValue(soildtTableName, cell_require)
		last_id := xlsx.GetCellValue(soildtTableName, cell_lastId)
		if location == "" {
			break
		}
		var soil datastruct.SoilData
		soil_id := StringToInt(location)
		soil.Price = StringToInt64(price)
		soil.Factor = StringToInt(factor)
		soil.Require = StringToInt(require)
		soil.LastId = StringToInt(last_id)
		soil.Level = 0
		soils[soil_id] = soil
		index++
	}

	index = 2
	petbarTableName := "Sheet2"
	petbars := make(map[datastruct.AnimalType]datastruct.PetbarData)
	for {
		cell_class := fmt.Sprintf("A%d", index)
		cell_price := fmt.Sprintf("B%d", index)
		cell_require := fmt.Sprintf("C%d", index)
		cell_id := fmt.Sprintf("D%d", index)
		cell_lastId := fmt.Sprintf("E%d", index)
		class := xlsx.GetCellValue(petbarTableName, cell_class)
		price := xlsx.GetCellValue(petbarTableName, cell_price)
		require := xlsx.GetCellValue(petbarTableName, cell_require)
		id := xlsx.GetCellValue(petbarTableName, cell_id)
		last_id := xlsx.GetCellValue(petbarTableName, cell_lastId)
		if class == "" {
			break
		}
		var petbar datastruct.PetbarData
		petbar_type := datastruct.AnimalType(StringToInt(class))
		petbar.Price = StringToInt64(price)
		petbar.Require = StringToInt(require)
		petbar.Id = StringToInt(id)
		petbar.LastId = StringToInt(last_id)
		petbars[petbar_type] = petbar
		index++
	}
	return soils, petbars
}

// func PlayerSoilToString(playerSoil []datastruct.PlayerSoil)(string,bool){
// 	 jsons, err := json.Marshal(playerSoil) //转换成JSON返回的是byte[]
// 	 if err != nil {
// 		log.Debug("PlayerSoilToString error:%s",err.Error())
// 		return "",true
// 	 }
//      return string(jsons),false
// }

func PlayerSoilToString(tmp *datastruct.PlayerSoil) (string, bool) {
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
		log.Debug("PlayerSoilToString error:%s", err.Error())
		return "", true
	}
	return string(jsons), false
}

func BytesToPlayerSoil(bytes []byte) (*datastruct.PlayerSoil, bool) {
	tmp := new(datastruct.PlayerSoil)
	err := json.Unmarshal(bytes, tmp)
	if err != nil {
		log.Debug("StringToPlayerSoil error:%s", err.Error())
		return nil, true
	}
	return tmp, false
}

func PlayerPetbarToString(tmp *datastruct.PlayerPetbar) (string, bool) {
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
		log.Debug("PlayerPetbarToString error:%s", err.Error())
		return "", true
	}
	return string(jsons), false
}

func BytesToPlayerPetbar(bytes []byte) (*datastruct.PlayerPetbar, bool) {
	tmp := new(datastruct.PlayerPetbar)
	err := json.Unmarshal(bytes, tmp)
	if err != nil {
		log.Debug("StringToPlayerSoil error:%s", err.Error())
		return nil, true
	}
	return tmp, false
}

func SpeedUpToString(tmp *datastruct.SpeedUpData) (string, bool) {
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
		log.Debug("PlayerSoilToString error:%s", err.Error())
		return "", true
	}
	return string(jsons), false
}

func BytesToSpeedUp(bytes []byte) (*datastruct.SpeedUpData, bool) {
	tmp := new(datastruct.SpeedUpData)
	err := json.Unmarshal(bytes, tmp)
	if err != nil {
		log.Debug("StringToPlayerSoil error:%s", err.Error())
		return nil, true
	}
	return tmp, false
}

func SliceIntToString(tmp []int) (string, bool) {
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
		log.Debug("SliceIntToString error:%s", err.Error())
		return "", true
	}
	return string(jsons), false
}

func BytesToSliceInt(bytes []byte) ([]int, bool) {
	var tmp []int
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		log.Debug("BytesToSliceInt error:%s", err.Error())
		return nil, true
	}
	return tmp, false
}

//level为升到当前多少级,price为土地初始价格
func ComputeSoilLevelPrice(gold int64, level int, current *datastruct.PlayerSoil) (int64, *datastruct.ResponseUpgradeSoil) {
	//compute
	resp_upsoil := new(datastruct.ResponseUpgradeSoil)
	rs_UpgradePrice := current.UpgradeLevelPrice
	rs_factor := current.Factor
	rs_level := current.Level
	for i := current.Level; i <= level; i++ {
		if gold < rs_UpgradePrice {
			break
		}
		gold -= rs_UpgradePrice
		rs_level += 1
		rs_factor += +20
		rs_UpgradePrice += 20
	}
	resp_upsoil.GoldCount = gold
	resp_upsoil.Level = rs_level
	resp_upsoil.UpgradePrice = rs_UpgradePrice
	resp_upsoil.Factor = rs_factor
	current.Level = rs_level
	current.UpgradeLevelPrice = rs_UpgradePrice
	current.Factor = rs_factor
	return gold, resp_upsoil
}

func EnableSpeedUp(ending int64, current int64) int64 {
	last := ending - current
	var CD int64
	CD = 0
	h4 := int64(4 * 3600)
	h24 := int64(24 * 3600)
	if last+h4 > h24 {
		CD = h4 - (h24 - last)
	}
	return CD
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func GetUpgradeLevelPriceForSoil(currentLevel int) int64 {
	return 100
}

func CreatePlayerSoil1(soil *datastruct.Soil1) *datastruct.PlayerSoil {
	rs := new(datastruct.PlayerSoil)
	rs.Factor = soil.Factor
	rs.Level = soil.Level
	rs.PlantId = soil.PlantId
	rs.PlantLevel = soil.PlantLevel
	rs.State = soil.State
	rs.UpgradeLevelPrice = soil.UpgradeLevelPrice
	return rs
}

func CreatePlayerSoil2(soil *datastruct.Soil2) *datastruct.PlayerSoil {
	rs := new(datastruct.PlayerSoil)
	rs.Factor = soil.Factor
	rs.Level = soil.Level
	rs.PlantId = soil.PlantId
	rs.PlantLevel = soil.PlantLevel
	rs.State = soil.State
	rs.UpgradeLevelPrice = soil.UpgradeLevelPrice
	return rs
}

func CreatePlayerSoil3(soil *datastruct.Soil3) *datastruct.PlayerSoil {
	rs := new(datastruct.PlayerSoil)
	rs.Factor = soil.Factor
	rs.Level = soil.Level
	rs.PlantId = soil.PlantId
	rs.PlantLevel = soil.PlantLevel
	rs.State = soil.State
	rs.UpgradeLevelPrice = soil.UpgradeLevelPrice
	return rs
}

func CreatePlayerSoil4(soil *datastruct.Soil4) *datastruct.PlayerSoil {
	rs := new(datastruct.PlayerSoil)
	rs.Factor = soil.Factor
	rs.Level = soil.Level
	rs.PlantId = soil.PlantId
	rs.PlantLevel = soil.PlantLevel
	rs.State = soil.State
	rs.UpgradeLevelPrice = soil.UpgradeLevelPrice
	return rs
}

func CreatePlayerSoil5(soil *datastruct.Soil5) *datastruct.PlayerSoil {
	rs := new(datastruct.PlayerSoil)
	rs.Factor = soil.Factor
	rs.Level = soil.Level
	rs.PlantId = soil.PlantId
	rs.PlantLevel = soil.PlantLevel
	rs.State = soil.State
	rs.UpgradeLevelPrice = soil.UpgradeLevelPrice
	return rs
}

func CreatePetbar1(petbar *datastruct.Petbar1) *datastruct.PlayerPetbar {
	rs := new(datastruct.PlayerPetbar)
	rs.AnimalNumber = petbar.AnimalNumber
	rs.CurrentExp = petbar.CurrentExp
	rs.State = petbar.State
	return rs
}

func CreatePetbar2(petbar *datastruct.Petbar2) *datastruct.PlayerPetbar {
	rs := new(datastruct.PlayerPetbar)
	rs.AnimalNumber = petbar.AnimalNumber
	rs.CurrentExp = petbar.CurrentExp
	rs.State = petbar.State
	return rs
}

func CreatePetbar3(petbar *datastruct.Petbar3) *datastruct.PlayerPetbar {
	rs := new(datastruct.PlayerPetbar)
	rs.AnimalNumber = petbar.AnimalNumber
	rs.CurrentExp = petbar.CurrentExp
	rs.State = petbar.State
	return rs
}

func CreatePetbar4(petbar *datastruct.Petbar4) *datastruct.PlayerPetbar {
	rs := new(datastruct.PlayerPetbar)
	rs.AnimalNumber = petbar.AnimalNumber
	rs.CurrentExp = petbar.CurrentExp
	rs.State = petbar.State
	return rs
}

//生成32位md5字串
func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func UniqueId() string {
	// 生成节点实例
	b := make([]byte, 48)
	if _, err := io.ReadFull(crypto_rand.Reader, b); err != nil {
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}

func ComputeCurrentGold(speedfactor int, soil map[int]*datastruct.PlayerSoil, petbar map[datastruct.AnimalType]*datastruct.PlayerPetbar, factor int, sec int64, plants []datastruct.Plant, animals map[datastruct.AnimalType]map[int]datastruct.Animal) int64 {
	//compute
	var addGold int64
	addGold = 0
	for _, v := range soil {
		if v.PlantId > 0 {
			rs_factor := float64(v.Factor / 100.0)
			addGold += int64(float64(plants[v.PlantId-1].InCome*sec*int64(factor*speedfactor)) * rs_factor)
		}
	}
	for k, v := range petbar {
		if v.AnimalNumber > 0 {
			animal := animals[k][v.AnimalNumber]
			addGold += animal.InCome * sec * int64(factor*speedfactor)
		}
	}
	return addGold
}

func GetGuestName(userId int) string {
	return fmt.Sprintf("游客%d", userId)
}
func GetGuestAvatar() string {
	return conf.Server.Domain + "/guest/avatar.png"
}
