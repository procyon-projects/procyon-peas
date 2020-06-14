package peas

import (
	"errors"
	core "github.com/procyon-projects/procyon-core"
	"reflect"
)

func CreateInstance(typ *core.Type, args []interface{}) (interface{}, error) {
	if core.IsFunc(typ) {
		in := make([]reflect.Value, 0)
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		result := typ.Val.Call(in)
		if len(result) != 1 {
			return nil, errors.New("it only supports the construction functions with one return parameter")
		}
		return result[0].Interface(), nil
	} else if core.IsStruct(typ) {
		if len(args) > 0 {
			return nil, errors.New("struct type does not support args")
		}
		return reflect.New(reflect.TypeOf(typ.Typ)), nil
	}
	return nil, errors.New("you can only pass Struct or Func types")
}
