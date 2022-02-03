package user

import (
	"unicode"
)

func VerifyPwd(s string, minLenLetter int) (lenOK, number, upper, special bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			//return false, false, false, false
		}
	}
	return letters >= minLenLetter, number, upper, special
}

func VerifyActive(s string) bool {
	return s == "T" || s == "F"
}

func VerifyRegtime(s string) bool {
	return s != ""
}

func VerifyTel(s string) bool {
	return len(s) > 3
}

func VerifyAddr(s string) bool {
	return true
}

func VerifyRole(s string) bool {
	return true
}

func VerifyLevel(s string) bool {
	return true
}

func VerifyExpire(s string) bool {
	return true
}

func VerifyAvatar(s string) bool {
	return true
}
