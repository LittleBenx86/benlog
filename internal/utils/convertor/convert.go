package convertor

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func Map2Struct(src map[string]interface{}, tarPtr interface{}) error {
	err := mapstructure.Decode(src, tarPtr)
	if err != nil {
		return err
	}
	return nil
}

func Struct2Map(objCopy interface{}) map[string]interface{} {
	objT := reflect.TypeOf(objCopy)
	objV := reflect.ValueOf(objCopy)

	var tar = make(map[string]interface{})
	for i := 0; i < objT.NumField(); i++ {
		tar[objT.Field(i).Name] = objV.Field(i).Interface()
	}
	return tar
}
