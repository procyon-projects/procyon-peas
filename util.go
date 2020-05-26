package peas

import (
	core "github.com/Rollcomp/procyon-core"
	"reflect"
)

func CreateInstance(typ *core.Type, args []interface{}) interface{} {
	if core.IsFunc(typ) {
		in := make([]reflect.Value, 0)
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		result := typ.Val.Call(in)
		if len(result) > 1 {
			panic("It only supports the construction functions with one parameter")
		}
		return result[0]
	} else if core.IsStruct(typ) {
		if len(args) > 0 {
			panic("Struct type does not support args")
		}
		return reflect.New(reflect.TypeOf(typ.Typ))
	}
	panic("You can only pass Struct or Func types")
}
