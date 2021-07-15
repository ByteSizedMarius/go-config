package go_config

import "reflect"

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

// Name returns the name of the referenced flag
func (f Flag) Name() string {
	return f.name
}

// Alias returns all aliases of the referenced flag
func (f Flag) Alias() []string {
	return f.alias
}

// Description returns the description of the referenced flag
func (f Flag) Description() string {
	return f.description
}

// Default returns the default-Value of the referenced flag
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

// Type returns the type of the referenced flag as a reflect.Kind
func (f Flag) Type() reflect.Kind {
	switch f.flagT.(type) {
	case *stringFlag:
		return reflect.String

	case *intFlag:
		return reflect.Int

	case *boolFlag:
		return reflect.Bool

	default:
		return reflect.Invalid
	}
}

// SetUseInCli sets whether the parameter will be accessible via commandline-parameters
// Default: True
func (f *Flag) SetUseInCli(use bool) *Flag {
	f.doNotUseInCli = !use
	return f
}

// SetUseAliasInCli allows setting whether the provided aliases will be accessible via the cli.
// Default: True
func (f *Flag) SetUseAliasInCli(use bool) *Flag {
	f.doNotUseAliasInCli = !use
	return f
}

// SetUseInINI sets whether the flag will be accessible in the ini-file
// Default: True
func (f *Flag) SetUseInINI(use bool) *Flag {
	f.doNotUseInIni = !use
	return f
}

// SetUseAliasInINI sets whether the flags aliases will be accessible in the ini-file
// Default: True
func (f *Flag) SetUseAliasInINI(use bool) *Flag {
	f.doNotUseAliasInIni = !use
	return f
}
