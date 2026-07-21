package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goccy/go-reflect"
)

func ParseType(value string, targetType reflect.Kind) (any, error) {
	var err error
	var r any

	switch targetType {
	case reflect.String:
		r = value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = strconv.Atoi(value)
	case reflect.Float32, reflect.Float64:
		r, err = strconv.ParseFloat(value, 64)
	case reflect.Bool:
		r, err = strconv.ParseBool(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = strconv.ParseUint(value, 10, 64)
	case reflect.TypeOf(time.Time{}).Kind():
		r, err = time.Parse(time.RFC3339, value)
	default:
		err = fmt.Errorf("unsupported type: %v", targetType)
	}

	return r, err
}