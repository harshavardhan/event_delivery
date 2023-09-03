package utils

import (
	"math/rand"
	"strconv"
)

func StrToInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// mock execution -> success : 200 , failure : 500 status code
func MockSuccess() bool {
	// failure to success ratio is 1:2
	return rand.Float64() >= (1.0 / 3.0)
}
