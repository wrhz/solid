package binding

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-reflect"
	"github.com/wrhz/solid/util"
)

type HeaderBinding struct{}

func (HeaderBinding) Name() string {
	return "header"
}

func (h HeaderBinding) Bind(req *http.Request, s any) error {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("BindHeaders: expected struct, got %v", v.Kind())
	}

	headers := req.Header

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		headerTag := field.Tag.Get("header")

		if headerTag == "" {
			headerTag = util.LowerFirst(field.Name)
		}

		paramType := field.Type

		data := headers[headerTag]
		paramTypeKind := paramType.Kind()

		if paramTypeKind == reflect.Slice {
			elementType := paramType.Elem().Kind()

			sliceValue := reflect.MakeSlice(field.Type, 0, len(data))

			for _, valStr := range data {
				
				elemValue, err := util.ParseType(valStr, elementType)

				if err != nil {
					return fmt.Errorf("parse slice element %q as %v: %w", valStr, elementType, err)
				}

				sliceValue = reflect.Append(sliceValue, reflect.ValueOf(elemValue))
			}

			v.Field(i).Set(sliceValue)
		} else {
			if len(data) == 0 {
				continue
			}

			value, err := util.ParseType(data[0], paramTypeKind)

			if err != nil {
				return fmt.Errorf("parse field %q as %v: %w", headerTag, paramType, err)
			}

			v.Field(i).Set(reflect.ValueOf(value))
		}
	}

	return nil
}