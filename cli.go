package go_config

import (
	cli "flag"
	"fmt"
	"os"
	s "strings"
)

func (c *Config) parseCLI() error {
	c.initializeCommandline()

	// make sure that no parameter has been set twice in the cli
	err := c.checkForDuplicates()
	if err != nil {
		return err
	}

	// evaluate commandline-parameters
	cli.Parse()

	// iterate all string flags
	for _, sf := range c.stringFlags {

		// Take the first non-empty value from the cli and set it into the struct
		for _, sfv := range sf.values {
			if *sfv != "" {
				sf.pointerToStructField.SetString(*sfv)
			}
		}
	}

	return nil
}

func (c *Config) checkForDuplicates() error {
	// Remove cli prefixes
	args := append(os.Args[:0:0], os.Args...) // copy slice
	for i := range args {
		args[i] = s.Trim(args[i], "-")
	}

	// iterate string flags
	for _, sf := range c.stringFlags {
		inc := false

		// iterate all possible names
		for _, sfn := range append(sf.f.alias, sf.f.name) {
			con := sliceContains(args, sfn)

			// if two names of one flag are included, throw error
			if inc && con {
				return fmt.Errorf("There were two values proviced for the commandline-parameter " + sfn + " via its aliases. Every commandline-parameter should only be set once.")
			} else if con {
				// if this is the first one included, set inc to true
				inc = true
			}
		}
	}

	// todo do this for other flags
	return nil
}

func (c *Config) initializeCommandline() {
	// iterate string flags
	for _, sf := range c.stringFlags {

		// results to each parameter will be in the map, using their name as a key
		sf.values = append(sf.values, cli.String(sf.f.name, "", sf.f.description))

		// since one parameter theoretically can have multiple values, they are stored in an array
		for _, sfa := range sf.f.alias {
			sf.values = append(sf.values, cli.String(sfa, "", sf.f.description))
		}
	}

	// todo iterate all other possible flag types
}
