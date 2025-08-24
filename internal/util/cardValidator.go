package util

import "strings"

func ValidCardNumber(number string) bool {
	number = strings.Replace(number, " ", "", -1)
	if len(number) != 12 {
		return false
	}
	return true
}
