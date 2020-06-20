// Package utils common utils, string util
package utils

import (
	"github.com/xwi88/kit4go/json"
)

// ToJsonString convert interface{} to json string with special json pkg
func ToJsonString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJsonStringIndent convert interface{} to json indent string with special json pkg
func ToJsonStringIndent(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SubString sub string include chinese character
func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// SubStringFromEnd sub string include chinese character
func SubStringFromEnd(str string, begin, end int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}
