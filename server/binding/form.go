package binding

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/goccy/go-reflect"
	"github.com/wrhz/solid/util"
)

type FormBinding struct{}

func (FormBinding) Name() string {
	return "form"
}

func (j FormBinding) Bind(req *http.Request, s any) error {
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("parse form: %w", err)
	}

	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("BindParams: expected struct, got %v", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		formTag := field.Tag.Get("form")

		if formTag == "" {
			formTag = util.LowerFirst(field.Name)
		}

		paramType := field.Type

		if paramType == reflect.TypeOf(multipart.FileHeader{}) {
			files := req.MultipartForm.File[formTag]

			if len(files) == 0 {
				continue
			}

			fh := files[0]
			if paramType == reflect.TypeOf(multipart.FileHeader{}) {
				v.Field(i).Set(reflect.ValueOf(*fh))
			} else if paramType == reflect.TypeOf((*multipart.FileHeader)(nil)) {
				v.Field(i).Set(reflect.ValueOf(fh))
			} else {
				return errors.New("unsupported file field type")
			}
		} else if paramType == reflect.TypeOf([]multipart.FileHeader{}) {
			files := req.MultipartForm.File[formTag]

			sliceValue := reflect.MakeSlice(paramType, 0, len(files))

			for _, fh := range files {
				if paramType.Elem() == reflect.TypeOf(multipart.FileHeader{}) {
					sliceValue = reflect.Append(sliceValue, reflect.ValueOf(*fh))
				} else if paramType.Elem() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
					sliceValue = reflect.Append(sliceValue, reflect.ValueOf(fh))
				} else {
					return errors.New("unsupported file slice field type")
				}
			}

			v.Field(i).Set(sliceValue)
		} else {
			data := req.Form[formTag]
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
					return fmt.Errorf("parse field %q as %v: %w", formTag, paramType, err)
				}

				v.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}

	return nil
}