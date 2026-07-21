package wasm

import "syscall/js"

func JSValueToGo(v js.Value) any {
	switch v.Type() {
	case js.TypeUndefined, js.TypeNull:
		return nil
	case js.TypeBoolean:
		return v.Bool()
	case js.TypeNumber:
		return v.Float()
	case js.TypeString:
		return v.String()
	case js.TypeSymbol:
		return v.String()
	case js.TypeObject, js.TypeFunction:
		if v.InstanceOf(js.Global().Get("Array")) {
			length := v.Length()
			arr := make([]any, length)
			for i := range arr {
				arr[i] = JSValueToGo(v.Index(i))
			}
			return arr
		}

		if v.Type() == js.TypeObject {
			keys := js.Global().Get("Object").Call("keys", v)
			length := keys.Length()
			obj := make(map[string]any, length)
			for i := 0; i < length; i++ {
				key := keys.Index(i).String()
				val := v.Get(key)
				obj[key] = JSValueToGo(val)
			}
			return obj
		}
		
		return nil
	default:
		return nil
	}
}