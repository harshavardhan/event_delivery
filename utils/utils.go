package utils

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func StrToInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func BuildKey(ts int64, id string) string {
	// need to convert ts to a fixed size string : ignoring for now
	return strconv.FormatInt(ts, 10) + "$" + id
}

func GetId(key string) string {
	arr := strings.Split(key, "$")
	return arr[1]
}

// mock execution -> success : 200 , failure : 500 status code
// We assume there are no bad payloads. Responses from destination are either 200 or 500 with 500 being retryable
func MockSuccess(destination string) bool {
	// ratio :: failure : success

	switch destination {
	case "A":
		// 0:1
		return true
	case "B":
		// 1:4
		return rand.Float64() >= (1.0 / 5.0)
	case "C":
		// 1:3
		return rand.Float64() >= (1.0 / 4.0)
	case "D":
		// 1:2
		return rand.Float64() >= (1.0 / 3.0)
	default:
		return false
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
