package utils

import (
	"strconv"
	"fmt"
)


func ConvertToStringPointer(value interface{}) *string {
	var vsp *string
	switch value.(type) {
	case int:
		vs := strconv.FormatInt(int64(value.(int)), 10)
		vsp = &vs
	case int32:
		vs := strconv.FormatInt(int64(value.(int32)), 10)
		vsp = &vs
	case int64:
		vs := strconv.FormatInt(int64(value.(int64)), 10)
		vsp = &vs
	case float32: 
		vs := strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
		vsp = &vs
	case float64: 
		vs := strconv.FormatFloat(value.(float64), 'f', -1, 64)
		vsp = &vs
	case string:
		vs := value.(string)
		vsp = &vs
	default:
		fmt.Println("Don't recognize this type")
	}

	return vsp

}