package util

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
)

var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Number interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
}

func Atoi[T Number](str string) T {
	num, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return T(num)
}

func GetStackTrace(err any) string {
	stackTrace := strings.Builder{}
	stackTrace.WriteString(fmt.Sprintf("%v", err))
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stackTrace.WriteString(fmt.Sprintf("\n - %s:%d (0x%x)", file, line, pc))
	}
	return stackTrace.String()
}

func RandString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
