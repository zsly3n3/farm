package tools

import (
	"strconv"
	"github.com/360EntSecGroup-Skylar/excelize"
	"farm/log"
	"farm/datastruct"
	"fmt"
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
		name := xlsx.GetCellValue(tableName, cell_Name)
		cid := xlsx.GetCellValue(tableName, cell_ClassId)
		price := xlsx.GetCellValue(tableName, cell_Price)
		income := xlsx.GetCellValue(tableName, cell_Income)
		exp := xlsx.GetCellValue(tableName, cell_AddExp)
		level := xlsx.GetCellValue(tableName, cell_Level)
        if name == "" {
            break
        }
        var plant datastruct.Plant
		plant.N = name
		plant.C = StringToInt(cid)
		plant.P = StringToInt(price)
		plant.I = StringToInt(income)
		plant.E = StringToInt(exp)
		plant.L = StringToInt(level)
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
		cell_Income:= fmt.Sprintf("C%d",index)
		cell_Exp:= fmt.Sprintf("D%d",index)
		cell_Level:= fmt.Sprintf("E%d",index)
		name := xlsx.GetCellValue(tableName, cell_Name)
		cid := xlsx.GetCellValue(tableName, cell_ClassId)
		income := xlsx.GetCellValue(tableName, cell_Income)
		exp := xlsx.GetCellValue(tableName, cell_Exp)
		level := xlsx.GetCellValue(tableName, cell_Level)
        if name == "" {
            break
        }
        var animal datastruct.Animal
		animal.N = name
		animal.C = StringToInt(cid)
		animal.I = StringToInt(income)
		animal.E = StringToInt(exp)
		animal.L = StringToInt(level)
        animals = append(animals,animal)
        index++
    }
	return animals
}


func GetSoildInfo()([]datastruct.SoilData,[]datastruct.PetbarData){
	xlsx, err := excelize.OpenFile("conf/soil_data.xlsx")
    if err != nil {
        log.Fatal("Excel error is %v", err.Error())
    }
	index:=2
	soildtTableName:="Sheet1"
	soils:=make([]datastruct.SoilData, 0,5)
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
		soil.Index = StringToInt(location)
		soil.Price = StringToInt(price)
		soil.Factor = StringToInt(factor)
		soil.Require = StringToInt(require)
		soil.Level = 1
		soil.Isbought = 0
		soil.PlantId =0
        soils = append(soils,soil)
        index++
	}
	
	petbarTableName:="Sheet2"
	petbars:=make([]datastruct.PetbarData, 0,4)
    for {
		cell_index  := fmt.Sprintf("A%d",index)
		cell_price := fmt.Sprintf("B%d",index)
		cell_require := fmt.Sprintf("C%d",index)
		location := xlsx.GetCellValue(petbarTableName, cell_index)
		price := xlsx.GetCellValue(petbarTableName, cell_price)
		require := xlsx.GetCellValue(petbarTableName, cell_require)
        if location == "" {
            break
		}
		var petbar datastruct.PetbarData
		petbar.Index = StringToInt(location)
		petbar.Price = StringToInt(price)
		petbar.Require = StringToInt(require)
		petbar.Isbought = 0
		petbar.AnimalId = 0
        petbars = append(petbars,petbar)
        index++
    }
	return soils,petbars
}