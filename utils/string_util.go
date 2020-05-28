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
