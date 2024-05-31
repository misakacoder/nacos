package router

import (
	"fmt"
	"strings"
)

func BlurQuery(field, value string, blur bool) (string, string) {
	joiner := "="
	if blur && value != "" {
		value = strings.TrimFunc(value, func(r rune) bool {
			return r == '*'
		})
		joiner = "like"
		value = "%" + value + "%"
	}
	return fmt.Sprintf("%s %s ?", field, joiner), value
}
