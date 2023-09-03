package utils

import (
	"math/rand"
	"strconv"
	"strings"
)

func StrToInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func BuildKey(ts int64, id string) string {
	return strconv.FormatInt(ts, 10) + "$" + id
}

func GetId(key string) string {
	arr := strings.Split(key, "$")
	return arr[1]
}

// mock execution -> success : 200 , failure : 500 status code
func MockSuccess() bool {
	// failure to success ratio is 1:2
	return rand.Float64() >= (1.0 / 3.0)
}
