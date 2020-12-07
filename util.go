package peas

import (
	"errors"
	"github.com/codnect/goo"
)

func CreateInstance(typ goo.Type, args []interface{}) (interface{}, error) {
	if typ.IsFunction() {
		fun := typ.ToFunctionType()
		if fun.GetFunctionReturnTypeCount() != 1 {
			return nil, errors.New("it only supports the construction functions with one return parameter")
		}
		results := fun.Call(args)
		return results[0], nil
	} else if typ.IsStruct() {
		if len(args) > 0 {
			return nil, errors.New("struct type does not support args")
		}
		return typ.ToStructType().NewInstance(), nil
	}
	return nil, errors.New("you can only pass Struct or Func types")
}

func getStringMapKeys(mapObj interface{}) []string {
	if mapObj == nil {
		return nil
	}
	mapType := goo.GetType(mapObj)
	if !mapType.IsMap() {
		panic("It is not an instance of map")
	}
	keyType := mapType.ToMapType().GetKeyType()
	if keyType.GetName() != "string" {
		panic("the key type of the given map is not string")
	}
	argMapKeys := mapType.GetGoValue().MapKeys()
	mapKeys := make([]string, len(argMapKeys))
	for i := 0; i < len(argMapKeys); i++ {
		mapKeys[i] = argMapKeys[i].String()
	}
	return mapKeys
}
