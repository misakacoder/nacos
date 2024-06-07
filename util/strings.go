package util

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

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
