package binding

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-reflect"
	"github.com/gorilla/mux"
	"github.com/wrhz/solid/util"
)

type ParamBinding struct{}

func (ParamBinding) Name() string {
	return "param"
}

func (j ParamBinding) Bind(req *http.Request, s any) error {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("BindParams: expected struct, got %v", v.Kind())
	}

	params := mux.Vars(req)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		paramTag := field.Tag.Get("param")

		if paramTag == "" {
			paramTag = util.LowerFirst(field.Name)
		}

		data := params[paramTag]

		paramType := field.Type.Kind()

		value, err := util.ParseType(data, paramType)

		if err != nil {
			return fmt.Errorf("parse field %q as %v: %w", paramTag, paramType, err)
		}

		v.Field(i).Set(reflect.ValueOf(value))
	}

	return nil
}