package main

import (
	"fmt"
	go_config "github.com/ByteSizedMarius/go-config"
)

type Config struct {
	Option1 string `name:"optionEins"`
	Option2 string `name:"optionZwei"`
	Option3 string `name:"optionDrei"`
}

func main() {
	test := &Config{}

	cs, err := go_config.Initialize(test)
	if err != nil {
		panic(err)
	}

	cs.NewString("optionEins", "default").SetAlias([]string{"o1"})
	cs.NewString("optionZwei", "default")
	cs.NewString("optionDrei", "default")

	err = cs.Parse()
	fmt.Println(err)
	fmt.Println(test.Option1)
}
