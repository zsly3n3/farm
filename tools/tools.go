package tools

import (
	"bytes"
	"encoding/binary"
)


func ByteArrToInt64(b *[]byte) int64{
	b_buf := bytes.NewBuffer(*b)
	var x int64
	binary.Read(b_buf, binary.BigEndian, &x)
	return x  
}