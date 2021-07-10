package main

import (
	go_config "github.com/ByteSizedMarius/go-config"
)

type Config struct {
	Option1 string
	Option2 string
	Option3 string
}

func main() {
	test := &Config{}

	cs, err := go_config.Initialize(test)
	if err != nil {
		panic(err)
	}

	_ = cs.NewString("optiodn1", "default")
}
