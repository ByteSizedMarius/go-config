package goconfig

import (
	cli "flag"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type Config struct {
	structToFill interface{}
	flags        []*Flag
	iniFile      string
}

/*
Tests:
int/bool/string jeweils einmal testen
int/bool/string jeweils einmal mit falschem typ im struct testen
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

	return Config{structToFill: stf}, nil
}

func (c *Config) Parse() error {
	err := c.parseINI()
	if err != nil {
		return errors.Wrap(err, "INI")
	}

	err = c.parseCLI()
	return errors.Wrap(err, "CLI")
}

// New adds a new parameter to the config-handler based on the given type
func (c *Config) New(name string, defaultValue interface{}) *Flag {
	if c.structToFill == nil{
		panic("Config was not initialized correctly. Use goconfig.Initialize().")
	}

	// Look for existing flags
	if getFlagFromNameOrAlias(c.flags, name) != nil {
		panic("duplicate flag: " + name)
	}

	// get the corresponding field from the struct
	f := getFieldFromName(*c, name)

	// check if the field exists in the struct and whether it's valid
	field, err := checkFieldValidity(f)
	if err != nil {
		panic(errors.Wrap(err, "error with field "+name))
	}

	fl := Flag{
		name:                 name,
		pointerToStructField: &field,
	}
	switch defaultValue.(type) {

	case string:
		if field.Type().Kind() != reflect.String {
			panic("cannot assign string-value of flag \"" + name + "\" to field of type " + fmt.Sprint(field.Type().Kind()))
		}
		field.SetString(defaultValue.(string))
		fl.flagT = &stringFlag{defValue: defaultValue.(string)}

	case int:
		if !(field.Type().Kind() == reflect.Int || field.Type().Kind() == reflect.Int64) {
			panic("cannot assign int-value of flag \"" + name + "\" to field of type " + fmt.Sprint(field.Type().Kind()))
		}
		field.SetInt(int64(defaultValue.(int)))
		fl.flagT = &intFlag{defValue: defaultValue.(int)}

	case bool:
		if field.Type().Kind() != reflect.Bool {
			panic("cannot assign bool-value of flag \"" + name + "\" to field of type " + fmt.Sprint(field.Type().Kind()))
		}
		field.SetBool(defaultValue.(bool))
		fl.flagT = &boolFlag{defValue: defaultValue.(bool)}

	default:
		panic("type not yet implemented. visit https://github.com/ByteSizedMarius/go-config/")
	}

	c.flags = append(c.flags, &fl)
	return &fl
}

// NewString adds a new string-parameter to the config-handler
func (c *Config) NewString(name string, defaultValue string) *Flag {
	return c.New(name, defaultValue)
}

// NewInt adds a new int-parameter to the config-handler
func (c *Config) NewInt(name string, defaultValue int) *Flag {
	return c.New(name, defaultValue)
}

// NewBool adds a new bool-parameter to the config-handler
func (c *Config) NewBool(name string, defaultValue bool) *Flag {
	return c.New(name, defaultValue)
}

// SetINI sets the ini-File to parse
func (c *Config) SetINI(iniFile string) *Config {
	c.iniFile = iniFile
	return c
}

func (c *Config) SetUsage(usage func()) {
	cli.Usage = usage
}
