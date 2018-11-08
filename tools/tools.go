package tools

import (
	"strconv"
	"github.com/360EntSecGroup-Skylar/excelize"
	"farm/log"
	"farm/datastruct"
	"fmt"
	"time"
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
		name := xlsx.GetCellValue(tableName, cell_Name)
		cid := xlsx.GetCellValue(tableName, cell_ClassId)
		price := xlsx.GetCellValue(tableName, cell_Price)
		income := xlsx.GetCellValue(tableName, cell_Income)
		exp := xlsx.GetCellValue(tableName, cell_AddExp)
        if name == "" {
            break
        }
        var plant datastruct.Plant
		plant.N = name
		plant.C = StringToInt(cid)
		plant.P = StringToInt(price)
		plant.I = StringToInt(income)
		plant.E = StringToInt(exp)
        plants = append(plants,plant)
        index++
    }
    return plants
}


func GetSoildInfo()[]datastruct.SoilData{
	xlsx, err := excelize.OpenFile("conf/soil_data.xlsx")
    if err != nil {
        log.Fatal("Excel error is %v", err.Error())
    }
	index:=2
	tableName:="Sheet1"
	soils:=make([]datastruct.SoilData, 0)
    for {
		cell_index  := fmt.Sprintf("A%d",index)
		cell_price := fmt.Sprintf("B%d",index)
		cell_factor := fmt.Sprintf("C%d",index)
		location := xlsx.GetCellValue(tableName, cell_index)
		price := xlsx.GetCellValue(tableName, cell_price)
		factor := xlsx.GetCellValue(tableName, cell_factor)
        if location == "" {
            break
		}
		var soil datastruct.SoilData
		soil.Index = StringToInt(location)
		soil.Price = StringToInt(price)
		soil.Factor = StringToInt(factor)
		soil.Level = 1
		soil.Isbought = 0
		soil.PlantID =0
        soils = append(soils,soil)
        index++
    }
	return soils
}


func CreateUser(code string,permissionId int)*datastruct.PlayerData{
	player:=new(datastruct.PlayerData)
	timestamp:=time.Now().Unix()
	player.PermissionId = permissionId
	player.CreatedAt = timestamp
	player.UpdateTime = timestamp
	player.Token = code
	player.GoldCount = 0
	player.HoneyCount = 0
	player.NickName = "test1"
	player.Avatar = "avatar"
	player.Soil = GetSoildInfo()
	return player
}