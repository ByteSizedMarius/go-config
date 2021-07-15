package go_config

import (
	"errors"
	cli "flag"
	"fmt"
	"os"
	s "strings"
)

func (c *Config) parseCLI() (err error) {
	c.initializeCommandline()

	// make sure that no parameter has been set twice in the cli
	err = c.checkForDuplicates()
	if err != nil {
		return err
	}

	// recover panic from flag-package
	defer func() {
		if errX := recover(); errX != nil {
			errY := errX.(error)
			if !errors.Is(cli.ErrHelp, errY) {
				err = fmt.Errorf(fmt.Sprint(errX))
			}
		}
	}()

	// evaluate commandline-parameters
	cli.Parse()

	// iterate all string flags
	for _, f := range c.flags {

		switch f.flagT.(type) {

		case *stringFlag:
			// Take the first non-empty value from the cli and set it into the struct
			for _, sfv := range f.flagT.(*stringFlag).values {
				if *sfv != "" {
					f.pointerToStructField.SetString(*sfv)
				}
			}

		case *intFlag:
			// Take the first non-empty value from the cli and set it into the struct
			for _, ifv := range f.flagT.(*intFlag).values {
				if *ifv != -1 {
					f.pointerToStructField.SetInt(int64(*ifv))
				}
			}

		case *boolFlag:
			// Take the first non-empty value from the cli and set it into the struct
			value := false
			for _, bfv := range f.flagT.(*boolFlag).values {
				value = value || *bfv
			}
			f.pointerToStructField.SetBool(value)
		}
	}

	return nil
}

func (c *Config) initializeCommandline() {
	cli.CommandLine = cli.NewFlagSet(os.Args[0], cli.PanicOnError)

	for _, f := range c.flags {
		if f.doNotUseInCli {
			continue
		}

		switch f.flagT.(type) {

		case *stringFlag:
			sf := f.flagT.(*stringFlag)

			// results to each parameter will be in the map, using their name as a key
			sf.values = append(sf.values, cli.String(f.name, "", f.description))

			// since one parameter theoretically can have multiple values, they are stored in an array
			if !f.doNotUseAliasInCli {
				for _, sfa := range f.alias {
					sf.values = append(sf.values, cli.String(sfa, "", f.description))
				}
			}

		case *intFlag:
			inf := f.flagT.(*intFlag)

			inf.values = append(inf.values, cli.Int(f.name, -1, f.description))
			if !f.doNotUseAliasInCli {
				for _, infa := range f.alias {
					inf.values = append(inf.values, cli.Int(infa, -1, f.description))
				}
			}

		case *boolFlag:
			bf := f.flagT.(*boolFlag)

			bf.values = append(bf.values, cli.Bool(f.name, false, f.description))
			if !f.doNotUseAliasInCli {
				for _, bfa := range f.alias {
					bf.values = append(bf.values, cli.Bool(bfa, false, f.description))
				}
			}
		}
	}
}

func (c *Config) checkForDuplicates() error {
	used := make(map[string]bool)

	// Remove cli prefixes
	args := append(os.Args[:0:0], os.Args...) // copy slice
	for i := range args {
		args[i] = s.Trim(args[i], "-")
	}

	for _, f := range c.flags {
		for _, fn := range append(f.alias, f.name) {
			if used[fn] {
				return fmt.Errorf("There were two values provided for the commandline-parameter " + fn + " via its aliases. Every commandline-parameter should only be set once.")
			}

			used[fn] = true
		}
	}

	return nil
}
