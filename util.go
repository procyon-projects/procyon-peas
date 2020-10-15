package peas

import (
	"errors"
	"github.com/codnect/goo"
	"reflect"
)

func CreateInstance(typ goo.Type, args []interface{}) (interface{}, error) {
	if typ.IsFunction() {
		in := make([]reflect.Value, 0)
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		/*result := typ.(goo.Function).Call(in)
		if len(result) != 1 {
			return nil, errors.New("it only supports the construction functions with one return parameter")
		}
		return result[0].Interface(), nil*/
		return nil, nil
	} else if typ.IsStruct() {
		if len(args) > 0 {
			return nil, errors.New("struct type does not support args")
		}
		return typ.(goo.Struct).NewInstance(), nil
	}
	return nil, errors.New("you can only pass Struct or Func types")
}
