package tools

import (
	"strconv"
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

