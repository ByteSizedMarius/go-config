package go_config

import (
	s "strings"
)

func sliceContains(sl []string, va string) bool {
	for _, v := range sl {
		if s.HasPrefix(v, va) {
			return true
		}
	}
	return false
}

func getStringFlagFromName(sfl []*StringFlag, name string) *StringFlag {
	for _, sf := range sfl {
		if sf.f.name == name {
			return sf
		}
	}
	return nil
}

func getStringFlagFromNameOrAlias(sfl []*StringFlag, name string) *StringFlag {
	for _, sf := range sfl {
		if sf.f.name == name || sliceContains(sf.f.alias, name) {
			return sf
		}
	}
	return nil
}
