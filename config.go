package go_config

import (
	cli "flag"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type Config struct {
	structToFill interface{}
	stringFlags  []*StringFlag
	Usage        func()
}

type flag struct {
	name        string
	alias       []string
	description string
}

type StringFlag struct {
	f                    flag
	values               []*string      // stores the flags values during parsing
	defValue             string         // the default-value
	pointerToStructField *reflect.Value // the pointer to the field in the struct
}

/*
Tests:
-a und -A sollten als unterschiedliche Parameter funktionieren können
sicherstellen dass die einzigartigkeit von flags richtig überprüft wird

*/

// Initialize checks parameters and creates a Config-struct
func Initialize(stf interface{}) (Config, error) {
	val := reflect.ValueOf(stf)

	if val.Kind() != reflect.Ptr {
		return Config{}, fmt.Errorf("structToFill must be a pointer")
	}

	if val.Elem().Kind() != reflect.Struct {
		return Config{}, fmt.Errorf("structToFill must be a pointer to a struct")
	}

	return Config{structToFill: stf, Usage: cli.Usage}, nil
}

func (c *Config) Parse() error {
	err := c.parseCLI()
	return err
}

// NewString adds a new string-parameter to the config-handler
func (c *Config) NewString(name string, defaultValue string) *StringFlag {
	if getStringFlagFromName(c.stringFlags, name) != nil {
		panic("duplicate flag: " + name)
	}

	// looks up the given string in the given struct
	f := getFieldFromName(*c, name)

	// check if the field exists in the struct and whether it's valid
	field, err := checkFieldValidity(f)
	if err != nil {
		panic(errors.Wrap(err, "error when setting string-parameter "+name))
	}
	field.SetString(defaultValue)

	sf := StringFlag{f: flag{name: name, alias: []string{}}, defValue: defaultValue, pointerToStructField: &field}
	c.stringFlags = append(c.stringFlags, &sf)
	return &sf
}

// SetAlias sets the aliases for the string-parameter
func (sf *StringFlag) SetAlias(alias []string) *StringFlag {
	sf.f.alias = alias
	return sf
}

// SetDescription sets the description used by the flag-pkg for help
func (sf *StringFlag) SetDescription(desc string) *StringFlag {
	sf.f.description = desc
	return sf
}
