package public

import (
	"regexp"
)

const (
	regular = "^1[3|4|5|6|7|8|9][0-9]\\d{8}$"
)

func IsValidMobile(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

