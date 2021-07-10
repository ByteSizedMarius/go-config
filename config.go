package go_config

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type Config struct {
	structToFill interface{}
	stringFlags  []StringFlag
}

type flag struct {
	name  string
	alias []string
}

type StringFlag struct {
	f                    flag
	value                string
	defValue             string
	pointerToStructField *reflect.Value
}

func Initialize(stf interface{}) (Config, error) {
	val := reflect.ValueOf(stf)

	if val.Kind() != reflect.Ptr {
		return Config{}, fmt.Errorf("structToFill must be a pointer")
	}

	if val.Elem().Kind() != reflect.Struct {
		return Config{}, fmt.Errorf("structToFill must be a struct")
	}

	return Config{structToFill: stf}, nil
}

// NewString adds a new string-parameter to the config-handler
func (c *Config) NewString(name string, defaultValue string) *StringFlag {

	// looks up the given string in the given struct
	f := getFieldFromName(*c, name)

	// check if the field exists in the struct and whether it's valid
	field, err := checkFieldValidity(f)
	if err != nil {
		panic(errors.Wrap(err, "error when setting string-parameter "+name))
	}

	sf := StringFlag{f: flag{name: name}, defValue: defaultValue, pointerToStructField: &field}
	c.stringFlags = append(c.stringFlags, sf)
	return &sf
}

// SetAlias sets the aliases for the string-parameter
func (sf *StringFlag) SetAlias(alias []string) *StringFlag {
	sf.f.alias = alias
	return sf
}
