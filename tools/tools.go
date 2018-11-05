package tools

import (
	"bytes"
	"encoding/binary"
)


func ByteArrToInt64(b *[]byte) int32{
	b_buf := bytes.NewBuffer(*b)
	var x int32
	binary.Read(b_buf, binary.BigEndian, &x)
	return x  
}

func ByteArrToBool(b *[]byte) bool{
	b_buf := bytes.NewBuffer(*b)
	var x int32
	binary.Read(b_buf, binary.BigEndian, &x)
	if x == 0{
	   return false
	}
	return true
}