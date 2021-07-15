package go_config

import (
	"fmt"
	"github.com/zieckey/goini"
	"strconv"
)

func (c *Config) parseINI() error {
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
			f := getFlagFromNameOrAlias(c.flags, iniKey)
			if f == nil {
				return fmt.Errorf("flag provided but not definied")
			}
			tmp := f

			// If usage of aliases in ini was disabled, check again
			if f.doNotUseAliasInIni {
				f = getFlagFromName(c.flags, iniKey)
			}

			// ini contains the alias of a flag where alias usage was disabled
			if f == nil {
				return fmt.Errorf("flag " + tmp.name + " was not configured to be used with aliases")
			}

			// ini contains flag where ini functionality was disabled
			if f.doNotUseInIni {
				return fmt.Errorf("flag " + f.name + " cannot be set via ini")
			}

			switch f.flagT.(type) {

			case *stringFlag:
				f.pointerToStructField.SetString(iniVal)

			case *intFlag:
				val, err := strconv.Atoi(iniVal)
				if err != nil {
					return err
				}

				f.pointerToStructField.SetInt(int64(val))

			case *boolFlag:
				val, err := strconv.ParseBool(iniVal)
				if err != nil {
					return err
				}
				f.pointerToStructField.SetBool(val)
			}
		}
	}

	return err
}

