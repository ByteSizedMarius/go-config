package go_config

import (
	"fmt"
	"github.com/zieckey/goini"
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
			sf.f.pointerToStructField.SetString(iniVal)
		}
	}

	return err
}
