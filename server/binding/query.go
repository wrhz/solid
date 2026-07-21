package binding

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-reflect"
	"github.com/wrhz/solid/util"
)

type QueryBinding struct{}

func (QueryBinding) Name() string {
	return "query"
}

func (j QueryBinding) Bind(req *http.Request, s any) error {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("BindQuery: expected struct, got %v", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		pathTag := field.Tag.Get("path")

		if pathTag == "" {
			pathTag = util.LowerFirst(field.Name)
		}

		paramValue := req.URL.Query()[pathTag]

		paramType := field.Type.Kind()

		if paramType != reflect.Slice {
			if len(paramValue) == 0 {
				continue
			}
			
			data, err := util.ParseType(paramValue[0], paramType)

			if err != nil {
				return fmt.Errorf("parse field %q as %v: %w", pathTag, paramType, err)
			}

			v.Field(i).Set(reflect.ValueOf(data))
		} else {
			elementType := field.Type.Elem().Kind()

			sliceValue := reflect.MakeSlice(field.Type, 0, len(paramValue))

			for _, valStr := range paramValue {
				elemValue, err := util.ParseType(valStr, elementType)
				if err != nil {
					return fmt.Errorf("parse slice element %q as %v: %w", valStr, elementType, err)
				}
				sliceValue = reflect.Append(sliceValue, reflect.ValueOf(elemValue))
			}

			v.Field(i).Set(sliceValue)
		}
	}
	return nil
}