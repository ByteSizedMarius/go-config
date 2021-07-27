package main

import (
	"fmt"
	"github.com/ByteSizedMarius/go-config"
)

type Configuration struct {
	PathToFile                string `name:"file"`
	LogVerbosity              int    `name:"logverb"`
	RequireElevatedPrivileges bool   `name:"requireElevatedPrivileges"`
}

func main() {
	myconfig := Configuration{}

	c, err := goconfig.Initialize(&myconfig)
	if err != nil {
		panic(err)
	}
	c.SetINI("./config.ini")

	c.New("file", "./README.md")
	c.New("logverb", 2).SetAlias([]string{"lv"})
	c.New("requireElevatedPrivileges", false)

	err = c.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(myconfig.PathToFile)
	fmt.Println(myconfig.LogVerbosity)
	fmt.Println(myconfig.RequireElevatedPrivileges)
}
