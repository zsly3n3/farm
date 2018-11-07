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
