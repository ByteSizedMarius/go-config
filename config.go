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
	intFlags     []*IntFlag
	boolFlags    []*BoolFlag
	Usage        func()
	iniFile      string
}

type flag struct {
	name                 string
	alias                []string
	description          string
	doNotUseInCli        bool
	doNotUseInIni        bool
	doNotUseAliasInCli   bool
	doNotUseAliasInIni   bool
	pointerToStructField *reflect.Value // the pointer to the field in the struct
}

type StringFlag struct {
	f        *flag
	values   []*string // stores the flags values during parsing
	defValue string    // the default-value
}

type IntFlag struct {
	f        *flag
	values   []*int
	defValue int
}

type BoolFlag struct {
	f        *flag
	values   []*bool
	defValue bool
}

/*
todo: check if field is correct type for current flag



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
func (c *Config) NewString(name string, defaultValue string) *StringFlag {
	if getStringFlagFromNameOrAlias(c.stringFlags, name) != nil {
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

	sf := StringFlag{f: &flag{name: name, alias: []string{}, pointerToStructField: &field}, defValue: defaultValue}
	c.stringFlags = append(c.stringFlags, &sf)
	return &sf
}

func (c *Config) SetINI(iniFile string) *Config {
	c.iniFile = iniFile
	return c
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

func (sf *StringFlag) SetUseAliasInCli(use bool) *StringFlag {
	sf.f.doNotUseAliasInCli = !use
	return sf
}

func (sf *StringFlag) SetUseInCli(use bool) *StringFlag {
	sf.f.doNotUseInCli = !use
	return sf
}