package peas

import (
	core "github.com/procyon-projects/procyon-core"
	"reflect"
)

func CreateInstance(typ *core.Type, args []interface{}) interface{} {
	if core.IsFunc(typ) {
		in := make([]reflect.Value, 0)
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		result := typ.Val.Call(in)
		if len(result) != 1 {
			core.Logger.Error("It only supports the construction functions with one return parameter")
			return nil
		}
		return result[0].Interface()
	} else if core.IsStruct(typ) {
		if len(args) > 0 {
			core.Logger.Error("Struct type does not support args")
			return nil
		}
		return reflect.New(reflect.TypeOf(typ.Typ))
	}
	core.Logger.Error("You can only pass Struct or Func types")
	return nil
}
