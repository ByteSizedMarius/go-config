package go_config

import (
	cli "flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zieckey/goini"
	"reflect"
)

type Config struct {
	structToFill interface{}
	stringFlags  []*StringFlag
	Usage        func()
	iniFile      string
}

type flag struct {
	name               string
	alias              []string
	description        string
	doNotUseInCli      bool
	doNotUseInIni      bool
	doNotUseAliasInCli bool
	doNotUseAliasInIni bool
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
sicherstellen dass die use-optionen im flag struct funken
todo: wenn alias in cli disabled aber angegeben, panic abfangen und stattdessen error zurückgeben?
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

func (c *Config) ParseINI() error {
	if c.iniFile == "" {
		return nil
	}

	ini := goini.New()
	err := ini.ParseFile(c.iniFile)

	configMap := ini.GetAll()

	// The configMap is a nested map
	// The outer map are the sections in the ini, we ignore these here.
	for _, section := range configMap {

		// The inner map is [string]string, consisting of the ini-keys and values
		for iniKey, iniVal := range section {

			// check if we have a flag with the same name
			sf := getStringFlagFromNameOrAlias(c.stringFlags, iniKey)
			if sf == nil {
				return fmt.Errorf("flag provided but not definied")
			}
			tmp := sf

			// If usage of aliases in ini was disabled, check again
			if sf.f.doNotUseAliasInIni {
				sf = getStringFlagFromName(c.stringFlags, iniKey)
			}

			// ini contains the alias of a flag where alias usage was disabled
			if sf == nil {
				return fmt.Errorf("flag " + tmp.f.name + " was not configured to be used with aliases")
			}

			// ini contains flag where ini functionality was disabled
			if sf.f.doNotUseInIni {
				return fmt.Errorf("flag " + sf.f.name + " cannot be set via ini")
			}
			sf.pointerToStructField.SetString(iniVal)
		}
	}

	return err
}

func (c *Config) Parse() error {
	err := c.ParseINI()
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

	sf := StringFlag{f: flag{name: name, alias: []string{}}, defValue: defaultValue, pointerToStructField: &field}
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

func (sf *StringFlag) SetUseAliasInCli(use bool) {
	sf.f.doNotUseAliasInCli = !use
}

func (sf *StringFlag) SetUseInCli(use bool) {
	sf.f.doNotUseInCli = !use
}
