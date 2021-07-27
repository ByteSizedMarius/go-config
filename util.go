package goconfig

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

func getFlagFromName(fl []*Flag, name string) *Flag {
	for _, f := range fl {
		if f.name == name {
			return f
		}
	}
	return nil
}

func getFlagFromNameOrAlias(fl []*Flag, name string) *Flag {
	for _, f := range fl {
		if f.name == name || sliceContains(f.alias, name) {
			return f
		}
	}
	return nil
}
