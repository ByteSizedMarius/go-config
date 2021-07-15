package go_config

import (
	cli "flag"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type Config struct {
	structToFill interface{}
	flags        []*Flag
	Usage        func()
	iniFile      string
}

type Flag struct {
	flagT                interface{}
	name                 string
	alias                []string
	description          string
	doNotUseInCli        bool
	doNotUseInIni        bool
	doNotUseAliasInCli   bool
	doNotUseAliasInIni   bool
	pointerToStructField *reflect.Value // the pointer to the field in the struct
}

type stringFlag struct {
	values   []*string // stores the flags values during parsing
	defValue string    // the default-value
}

type intFlag struct {
	values   []*int
	defValue int
}

type boolFlag struct {
	values   []*bool
	defValue bool
}

/*
todo: check if field is correct type for current flag
todo: setuseinini setusealiasinini

Tests:
sicherstellen dass die ini-setter funktionieren
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
	err := c.parseINI()
	if err != nil {
		return errors.Wrap(err, "INI")
	}

	err = c.parseCLI()
	return errors.Wrap(err, "CLI")
}

// NewString adds a new string-parameter to the config-handler
func (c *Config) NewString(name string, defaultValue string) *Flag {
	if getFlagFromNameOrAlias(c.flags, name) != nil {
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

	fl := Flag{
		flagT:                &stringFlag{defValue: defaultValue},
		name:                 name,
		pointerToStructField: &field,
	}
	c.flags = append(c.flags, &fl)
	return &fl
}

// SetINI sets the ini-File to parse
func (c *Config) SetINI(iniFile string) *Config {
	c.iniFile = iniFile
	return c
}

// SetAlias sets the aliases for the string-parameter
func (f *Flag) SetAlias(alias []string) *Flag {
	f.alias = alias
	return f
}

// SetDescription sets the description used by the flag-pkg for help
func (f *Flag) SetDescription(desc string) *Flag {
	f.description = desc
	return f
}

func (f Flag) Name() string {
	return f.name
}

func (f Flag) Alias() []string {
	return f.alias
}

func (f Flag) Description() string {
	return f.description
}

func (f Flag) Default() interface{} {
	switch f.flagT.(type) {
	case *stringFlag:
		return f.flagT.(*stringFlag).defValue

	case *intFlag:
		return f.flagT.(*intFlag).defValue

	case *boolFlag:
		return f.flagT.(*boolFlag).defValue
	}
	return nil
}

func (f Flag) Type() string {
	switch f.flagT.(type) {
	case *stringFlag:
		return "string"

	case *intFlag:
		return "int"

	case *boolFlag:
		return "bool"

	default:
		return ""
	}
}
