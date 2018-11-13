package tools

import (
	"strconv"
	"github.com/360EntSecGroup-Skylar/excelize"
	"farm/log"
	"farm/datastruct"
	"fmt"
	"encoding/json"
)

func Int64ToString(tmp int64) string{
return strconv.FormatInt(tmp,10)
}

func StringToInt64(tmp string) int64{
rs, _ := strconv.ParseInt(tmp, 10, 64)
return rs
}

func IntToString(tmp int) string{
return strconv.Itoa(tmp)
}

func StringToInt(tmp string) int{
rs,_:=strconv.Atoi(tmp)
return rs
}

func StringToFloat64(tmp string)float64{
	rs,_ := strconv.ParseFloat(tmp,64)
	return rs
}

func BoolToString(tmp bool) string{
	if tmp == false{
		return "0"	
	}else{
		return "1"
	}
}

func StringToBool(tmp string) bool{
	if tmp == "0"{
		return false
	}else{
		return true
	}
}

func GetPlantsInfo()[]datastruct.Plant{
    xlsx, err := excelize.OpenFile("conf/shop_data.xlsx")
    if err != nil {
        log.Fatal("Excel error is %v", err.Error())
    }
	index:=2
	tableName:="Sheet1"
    plants:=make([]datastruct.Plant, 0)
    for {
		cell_Name  := fmt.Sprintf("A%d",index)
		cell_ClassId := fmt.Sprintf("B%d",index)
		cell_Price := fmt.Sprintf("C%d",index)
		cell_Income:= fmt.Sprintf("D%d",index)
		cell_AddExp:= fmt.Sprintf("E%d",index)
		cell_Level:= fmt.Sprintf("F%d",index)
		cell_ChName:= fmt.Sprintf("G%d",index)
		name := xlsx.GetCellValue(tableName, cell_Name)
		cid := xlsx.GetCellValue(tableName, cell_ClassId)
		price := xlsx.GetCellValue(tableName, cell_Price)
		income := xlsx.GetCellValue(tableName, cell_Income)
		exp := xlsx.GetCellValue(tableName, cell_AddExp)
		level := xlsx.GetCellValue(tableName, cell_Level)
		chName := xlsx.GetCellValue(tableName, cell_ChName)
        if name == "" {
            break
        }
        var plant datastruct.Plant
		plant.Name = name
		plant.ClassId = StringToInt(cid)
		plant.Price = StringToInt(price)
		plant.InCome = StringToInt(income)
		plant.ExpForAnimal = StringToInt(exp)
		plant.Level = StringToInt(level)
		plant.CName = chName
        plants = append(plants,plant)
        index++
    }
    return plants
}

func GetAnimalInfo()[]datastruct.Animal{
	xlsx, err := excelize.OpenFile("conf/shop_data.xlsx")
    if err != nil {
        log.Fatal("Excel error is %v", err.Error())
    }
	index:=2
	tableName:="Sheet2"
    animals:=make([]datastruct.Animal, 0)
    for {
		cell_Name  := fmt.Sprintf("A%d",index)
		cell_ClassId := fmt.Sprintf("B%d",index)
		cell_InCome:= fmt.Sprintf("C%d",index)
		cell_Exp:= fmt.Sprintf("D%d",index)
		cell_Number:= fmt.Sprintf("E%d",index)
		name := xlsx.GetCellValue(tableName, cell_Name)
		cid := xlsx.GetCellValue(tableName, cell_ClassId)
		income := xlsx.GetCellValue(tableName, cell_InCome)
		exp := xlsx.GetCellValue(tableName, cell_Exp)
		number := xlsx.GetCellValue(tableName, cell_Number)
        if name == "" {
            break
        }
        var animal datastruct.Animal
		animal.Name = name
		animal.ClassId = StringToInt(cid)
		animal.InCome = StringToInt(income)
		animal.Exp = StringToInt64(exp)
		animal.Number = StringToInt(number)
        animals = append(animals,animal)
        index++
    }
	return animals
}


func GetSoildInfo()(map[int]datastruct.SoilData,map[datastruct.AnimalType]datastruct.PetbarData){
	xlsx, err := excelize.OpenFile("conf/soil_data.xlsx")
    if err != nil {
        log.Fatal("Excel error is %v", err.Error())
    }
	index:=2
	soildtTableName:="Sheet1"
	soils:=make(map[int]datastruct.SoilData)
    for {
		cell_index  := fmt.Sprintf("A%d",index)
		cell_price := fmt.Sprintf("B%d",index)
		cell_factor := fmt.Sprintf("C%d",index)
		cell_require := fmt.Sprintf("D%d",index)
		location := xlsx.GetCellValue(soildtTableName, cell_index)
		price := xlsx.GetCellValue(soildtTableName, cell_price)
		factor := xlsx.GetCellValue(soildtTableName, cell_factor)
		require := xlsx.GetCellValue(soildtTableName, cell_require)
        if location == "" {
            break
		}
		var soil datastruct.SoilData
		soil_id := StringToInt(location)
		soil.Price = StringToInt(price)
		soil.Factor = StringToInt(factor)
		soil.Require = StringToInt(require)
		soil.Level = 0
		soils[soil_id]=soil
        index++
	}

	index=2
	petbarTableName:="Sheet2"
	petbars:=make(map[datastruct.AnimalType]datastruct.PetbarData)
    for {
		cell_class  := fmt.Sprintf("A%d",index)
		cell_price := fmt.Sprintf("B%d",index)
		cell_require := fmt.Sprintf("C%d",index)
		class := xlsx.GetCellValue(petbarTableName, cell_class)
		price := xlsx.GetCellValue(petbarTableName, cell_price)
		require := xlsx.GetCellValue(petbarTableName, cell_require)
        if class == "" {
            break
		}
		var petbar datastruct.PetbarData
		petbar_type:= datastruct.AnimalType(StringToInt(class))
		petbar.Price = StringToInt(price)
		petbar.Require = StringToInt(require)
		petbars[petbar_type]=petbar
        index++
    }
	return soils,petbars
}

// func PlayerSoilToString(playerSoil []datastruct.PlayerSoil)(string,bool){
// 	 jsons, err := json.Marshal(playerSoil) //转换成JSON返回的是byte[]
// 	 if err != nil {
// 		log.Debug("PlayerSoilToString error:%s",err.Error())
// 		return "",true
// 	 }
//      return string(jsons),false
// }

func PlayerSoilToString(tmp *datastruct.PlayerSoil)(string,bool){
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
	   log.Debug("PlayerSoilToString error:%s",err.Error())
	   return "",true
	}
	return string(jsons),false
}

func BytesToPlayerSoil(bytes []byte)(*datastruct.PlayerSoil,bool){
	tmp:=new(datastruct.PlayerSoil)
    err := json.Unmarshal(bytes,tmp)
    if err != nil {
        log.Debug("StringToPlayerSoil error:%s",err.Error())
		return nil,true
    }
	return tmp,false
}

func PlayerPetbarToString(tmp *datastruct.PlayerPetbar)(string,bool){
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
	   log.Debug("PlayerPetbarToString error:%s",err.Error())
	   return "",true
	}
	return string(jsons),false
}

func BytesToPlayerPetbar(bytes []byte)(*datastruct.PlayerPetbar,bool){
	tmp:=new(datastruct.PlayerPetbar)
    err := json.Unmarshal(bytes,tmp)
    if err != nil {
        log.Debug("StringToPlayerSoil error:%s",err.Error())
		return nil,true
    }
	return tmp,false
}


func SliceIntToString(tmp []int)(string,bool){
	jsons, err := json.Marshal(tmp) //转换成JSON返回的是byte[]
	if err != nil {
	   log.Debug("SliceIntToString error:%s",err.Error())
	   return "",true
	}
	return string(jsons),false
}

func BytesToSliceInt(bytes []byte)([]int,bool){
	var tmp []int
    err := json.Unmarshal(bytes,&tmp)
    if err != nil {
        log.Debug("BytesToSliceInt error:%s",err.Error())
		return nil,true
    }
	return tmp,false
}


func GetUpgradeLevelPriceForSoil(currentLevel int)int{
	 return 100
}