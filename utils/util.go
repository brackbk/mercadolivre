package utils

import (
	"regexp"
	"unicode/utf8"
)

func Ean(str string) string {
	var newEan string
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	submatchall := re.FindAllString(str, -1)
	for _, element := range submatchall {
		newEan = newEan + element
	}

	if utf8.RuneCountInString(newEan) < 13 {
		newEan = " "
	}
	return newEan
}

func BoolNotNull(v interface{}) bool {
	if v == nil {
		return false
	} else {
		return v.(bool)
	}
}

func StringNotNull(v interface{}) string {
	if v == nil {
		return " "
	}
	return v.(string)
}

func FloatNotNull(v interface{}) float64 {

	if v == nil {
		return 0.0
	}
	return v.(float64)

}

func IntNotNull(v interface{}) int {

	if v == nil {
		return 0
	}
	return int(v.(float64))

}
