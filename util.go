package peas

import (
	"errors"
	"github.com/codnect/goo"
)

func CreateInstance(typ goo.Type, args []interface{}) (interface{}, error) {
	if typ.IsFunction() {
		fun := typ.(goo.Function)
		if fun.GetFunctionReturnTypeCount() != 1 {
			return nil, errors.New("it only supports the construction functions with one return parameter")
		}
		results := fun.Call(args)
		return results[0], nil
	} else if typ.IsStruct() {
		if len(args) > 0 {
			return nil, errors.New("struct type does not support args")
		}
		return typ.(goo.Struct).NewInstance(), nil
	}
	return nil, errors.New("you can only pass Struct or Func types")
}
