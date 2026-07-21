package wasm

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"sync/atomic"
	"syscall/js"
	"unicode"
)

var moduleName = os.Args[0]
const ptrKind = reflect.Kind(22)

var (
	instanceMap   = make(map[uint64]any)
	instanceIDSeq uint64
)

func lowerFirst(s string) string {
    if s == "" {
        return s
    }
    r := []rune(s)
    r[0] = unicode.ToLower(r[0])
    return string(r)
}

func getValueMethodNames(t reflect.Type) []string {
    if t.Kind() == ptrKind {
        t = t.Elem()
    }
    names := []string{}
    for i := 0; i < t.NumMethod(); i++ {
        names = append(names, t.Method(i).Name)
    }
    return names
}

func getType(dataType reflect.Type) reflect.Type {
	if dataType.Kind() == ptrKind {
		for dataType.Kind() == ptrKind {
			dataType = dataType.Elem()
		}
	}

	return dataType
}

func getValue(dataValue reflect.Value) reflect.Value {
	if dataValue.Kind() == ptrKind {
		for dataValue.Kind() == ptrKind {
			dataValue = dataValue.Elem()
		}
	}

	return dataValue
}

func createValueMethod(method reflect.Value) js.Func {
	return addFunc(js.FuncOf(func(this js.Value, args []js.Value) any {
		results := method.Call([]reflect.Value{
			reflect.ValueOf(this),
			reflect.ValueOf(args),
		})

		if len(results) > 0 {
			return results[0].Interface()
		}

		return nil
	}))
}

func createPrototypeMethod(method reflect.Method) js.Func {
	return addFunc(js.FuncOf(func(this js.Value, args []js.Value) any {
		id := this.Get("_goId").Int()
		structPtr := instanceMap[uint64(id)]

		results := method.Func.Call([]reflect.Value{
			reflect.ValueOf(structPtr),
			reflect.ValueOf(this),
			reflect.ValueOf(args),
		})

		if len(results) > 0 {
			return results[0].Interface()
		}

		return nil
	}))
}

func valueGetter(instance reflect.Value, fieldName string) js.Func {
	return addFunc(js.FuncOf(func(this js.Value, args []js.Value) any {
		return instance.Elem().FieldByName(fieldName).Interface()
	}))
}

func valueSetter(tag string, instance reflect.Value, fieldName string) js.Func {
    return addFunc(js.FuncOf(func(this js.Value, args []js.Value) any {
        target := args[0]
        goVal := JSValueToGo(target)

        field := instance.Elem().FieldByName(fieldName)
        if !field.IsValid() || !field.CanSet() {
            return fmt.Errorf("cannot set field %s", fieldName)
        }

        converted := reflect.ValueOf(goVal)
        if converted.Type() != field.Type() {
            if converted.Type().ConvertibleTo(field.Type()) {
                converted = converted.Convert(field.Type())
            } else {
                switch field.Kind() {
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                    if f, ok := goVal.(float64); ok {
                        converted = reflect.ValueOf(int(f)).Convert(field.Type())
                    }
                case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                    if f, ok := goVal.(float64); ok {
                        converted = reflect.ValueOf(uint(f)).Convert(field.Type())
                    }
                case reflect.Float32, reflect.Float64:
                default:
                    return fmt.Errorf("unsupported field type for JS number")
                }
            }
        }

        field.Set(converted)
        this.Set("_" + tag, target)
        return nil
    }))
}

func classGetter(class js.Func) js.Func {
	return addFunc(js.FuncOf(func(this js.Value, args []js.Value) any {
		return class
	}))
}

func createClass(classValue reflect.Value) (js.Func, error) {
	classType := getType(classValue.Type())

	if classType.Kind() != reflect.Struct { return addFunc(js.FuncOf(func(this js.Value, args []js.Value) any { return nil })), fmt.Errorf("The class type should be a struct") }

	classValue = getValue(classValue)

	if classValue.Kind() != reflect.Struct { return addFunc(js.FuncOf(func(this js.Value, args []js.Value) any { return nil })), fmt.Errorf("The class value should be a struct") }

	prototypeMethods := map[string]any{}
	valueMethods := map[string]any{}
	classTypePtr := reflect.PointerTo(classType)
	methodNames := getValueMethodNames(classType)

	var constructorFunc reflect.Value

	for _, methodName := range methodNames {
		method, _ := classType.MethodByName(methodName)
		bound := classValue.MethodByName(methodName)

		valueMethods[lowerFirst(method.Name)] = createValueMethod(bound)
	}

	for i := 0; i < classTypePtr.NumMethod(); i++ {
		method := classTypePtr.Method(i)

		if slices.Contains(methodNames, method.Name) { continue }

		name := method.Name
		
		if name == "Constructor" {
			constructorFunc = method.Func
		} else {
			prototypeMethods[lowerFirst(name)] = createPrototypeMethod(method)
		}
	}

	constructor := addFunc(js.FuncOf(func(this js.Value, args []js.Value) any {
		newPtr := reflect.New(classType)
    	newStruct := newPtr.Elem()

		id := atomic.AddUint64(&instanceIDSeq, 1)
		instanceMap[id] = newPtr.Interface()
		this.Set("_goId", js.ValueOf(id))

		for i := 0; i < classType.NumField(); i++ {
			field := classType.Field(i)
			fieldValue := classValue.Field(i)

			newStruct.Field(i).Set(fieldValue)

			tag := field.Tag.Get("wasm")

			if tag == "" {
				tag = lowerFirst(classType.Field(i).Name)
			}

			fieldKind := getType(field.Type)

			var descriptor = Object().New()

			if fieldKind.Kind() == reflect.Struct {
				data, err := createClass(fieldValue)

				if err != nil {
					return err
				}

				descriptor.Set("get", classGetter(data))
				descriptor.Set("set", js.FuncOf(func(this js.Value, args []js.Value) any {
					Console().Call("error", "The class can't get anything")

					return nil
				}))

				this.Set(tag, data)
			} else {
				fieldValue := newStruct.Field(i)
				fieldName := field.Name

				descriptor.Set("get", valueGetter(newPtr, fieldName))
				descriptor.Set("set", valueSetter(tag, newPtr, fieldName))

				this.Set("_" + tag, js.ValueOf(fieldValue.Interface()))
			}

			descriptor.Set("enumerable", true)
    		descriptor.Set("configurable", true)

			Object().Call("defineProperty", this, tag, descriptor)
		}

		if constructorFunc.IsValid() {
			constructorFunc.Call([]reflect.Value{
				newPtr,
				reflect.ValueOf(this),
				reflect.ValueOf(args),
			})
		}

		return nil
	}))

	prototype := constructor.Get("prototype")

	for name, valueMethod := range valueMethods {
		constructor.Set(name, valueMethod)
	}

	for name, prototypeMethod := range prototypeMethods {
		prototype.Set(name, prototypeMethod)
	}

	return constructor, nil
}

func ExportVar(name string, value any) {
	objectClass := Object()

	goExports := getGoExports()

	if !objectClass.Call("hasOwn", goExports, js.ValueOf(moduleName)).Bool() {
		object := objectClass.New()

		goExports.Set(moduleName, object)
	}

	goExports.Get(moduleName).Set(name, js.ValueOf(value))
}

func ExportFunc(name string, callFunc func(this js.Value, args []js.Value) any) {
	objectClass := Object()

	goExports := getGoExports()

	if !objectClass.Call("hasOwn", goExports, js.ValueOf(moduleName)).Bool() {
		object := objectClass.New()

		goExports.Set(moduleName, object)
	}

	goExports.Get(moduleName).Set(name, addFunc(js.FuncOf(callFunc)))
}

func ExportClass(name string, class any) error {
	objectClass := Object()

	goExports := getGoExports()

	if !objectClass.Call("hasOwn", goExports, js.ValueOf(moduleName)).Bool() {
		object := objectClass.New()

		goExports.Set(moduleName, object)
	}

	classData, err := createClass(reflect.ValueOf(class))

	if err != nil {
		return err
	}

	goExports.Get(moduleName).Set(name, classData)

	return nil
}
