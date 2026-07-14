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
	instanceMap   = make(map[uint64]interface{})
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
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    names := []string{}
    for i := 0; i < t.NumMethod(); i++ {
        names = append(names, t.Method(i).Name)
    }
    return names
}

func getObjectClass() js.Value {
	return js.Global().Get("Object")
}

func getGoExports() js.Value {
	goExports := js.Global().Get("goExports")

	if goExports.IsUndefined() || goExports.IsNull() {
		js.Global().Set("goExports", js.Global().Get("Map").New())
		goExports = js.Global().Get("goExports")
	}

	return goExports
}

func createValueMethod(method reflect.Value) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		results := method.Call([]reflect.Value{
			reflect.ValueOf(this),
			reflect.ValueOf(args),
		})

		if len(results) > 0 {
			return results[0].Interface()
		}

		return nil
	})
}

func createPrototypeMethod(method reflect.Method) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
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
	})
}

func createClass(classValue reflect.Value) (js.Func, error) {
	classType := classValue.Type()

	if classType.Kind() == ptrKind {
		for classType.Kind() == ptrKind {
			classType =classType.Elem()
		}
	}

	if classType.Kind() != reflect.Struct { return js.FuncOf(func(this js.Value, args []js.Value) any { return nil }), fmt.Errorf("The class type should be a struct") }

	for classValue.Kind() == ptrKind {
		classValue = classValue.Elem()
	}

	if classValue.Kind() != reflect.Struct { return js.FuncOf(func(this js.Value, args []js.Value) any { return nil }), fmt.Errorf("The class value should be a struct") }

	structType := classType

	type fieldInfo struct {
        index int
        name  string
    }
    var fields []fieldInfo
    for i := 0; i < structType.NumField(); i++ {
        tag := structType.Field(i).Tag.Get("wasm")
        if tag == "" {
            tag = lowerFirst(structType.Field(i).Name)
        }
        fields = append(fields, fieldInfo{index: i, name: tag})
    }

	prototypeMethods := map[string]any{}
	valueMethods := map[string]any{}
	classTypePtr := reflect.PointerTo(classType)
	methodNames := getValueMethodNames(classType)

	var constructorFunc reflect.Value

	for i := 0; i < classType.NumMethod(); i++ {
		method := classType.Method(i)
		bound := classValue.Method(i)

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

	constructor := js.FuncOf(func(this js.Value, args []js.Value) any {
		newPtr := reflect.New(structType)
        newStruct := newPtr.Elem()

		for _, f := range fields {
            if newStruct.Field(f.index).CanSet() {
                newStruct.Field(f.index).Set(classValue.Field(f.index))
            }
        }

        for _, f := range fields {
            this.Set(f.name, js.ValueOf(newStruct.Field(f.index).Interface()))
        }

		id := atomic.AddUint64(&instanceIDSeq, 1)
		instanceMap[id] = newPtr.Interface()
		this.Set("_goId", js.ValueOf(id))

		if constructorFunc.IsValid() {
			constructorFunc.Call([]reflect.Value{
				newPtr,
				reflect.ValueOf(this),
				reflect.ValueOf(args),
			})
		}

		return nil
	})

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
	objectClass := getObjectClass()

	goExports := getGoExports()

	if !objectClass.Call("hasOwn", goExports, js.ValueOf(moduleName)).Bool() {
		object := objectClass.New()

		goExports.Set(moduleName, object)
	}

	goExports.Get(moduleName).Set(name, js.ValueOf(value))
}

func ExportFunc(name string, callFunc func(this js.Value, args []js.Value) any) {
	objectClass := getObjectClass()

	goExports := getGoExports()

	if !objectClass.Call("hasOwn", goExports, js.ValueOf(moduleName)).Bool() {
		object := objectClass.New()

		goExports.Set(moduleName, object)
	}

	goExports.Get(moduleName).Set(name, js.FuncOf(callFunc))
}

func ExportClass(name string, class any) error {
	objectClass := getObjectClass()

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
