package util

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"reflect"
)

func ToJSONString(object any) (string, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func ParseObject[T any](text string) (T, error) {
	var object T
	err := json.Unmarshal([]byte(text), &object)
	return object, err
}

func Copy(source, target any) {
	err := copier.CopyWithOption(target, source, copier.Option{DeepCopy: true})
	if err != nil {
		panic(err)
	}
}

func ZeroValue(v any) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}
