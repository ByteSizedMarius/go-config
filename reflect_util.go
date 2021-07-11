package go_config

import (
	"fmt"
	"reflect"
)

func checkFieldValidity(f interface{}) (reflect.Value, error) {
	if f == nil {
		return reflect.Value{}, fmt.Errorf("field could not be found")
	}

	field := f.(reflect.Value)

	if !field.IsValid() {
		return field, fmt.Errorf("field is not valid")
	}

	if !field.CanSet() {
		return field, fmt.Errorf("field cannot be set")
	}

	return field, nil
}

func getFieldFromName(c Config, name string) interface{} {
	val := reflect.ValueOf(c.structToFill).Elem()

	for i := 0; i < val.NumField(); i++ {
		tagName, ok := val.Type().Field(i).Tag.Lookup("name")

		if ok && tagName == name {
			return val.Field(i)
		}
	}

	return nil
}
