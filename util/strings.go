package util

import "strconv"

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
