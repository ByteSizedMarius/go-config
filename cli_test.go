package go_config

import (
	"os"
	"testing"
)

func TestCli(t *testing.T) {
	t.Run("TestUseInCli", func(t *testing.T) {
		os.Args = buildArgs([]string{"-test1=x1"})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetUseInCli(true)

		err := cs.Parse()
		if err != nil {
			t.Error(err)
		}
		t.Log("success")
	})

	t.Run("TestDoNotUseInCli", func(t *testing.T) {
		os.Args = buildArgs([]string{"-test1=x1", "-test2=x2"})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetUseInCli(false)
		cs.NewString("test2", "null")

		err := cs.Parse()
		if err == nil {
			t.Error("no error thrown")
		}
	})

	t.Run("TestUseAliasInCli", func(t *testing.T) {
		os.Args = buildArgs([]string{"-testAlias1=x1", "-test2=x2"})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetAlias([]string{"testAlias1"})
		cs.NewString("test2", "null")

		err := cs.Parse()
		if err != nil {
			t.Error(err)
		}

		if x.Option1 != "x1" {
			t.Error("value was not set correctly")
		}
	})

	t.Run("TestDoNotUseInCli", func(t *testing.T) {
		os.Args = buildArgs([]string{"-testAlias1=x1", "-test2=x2"})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetUseAliasInCli(false).SetAlias([]string{"testAlias1"})
		cs.NewString("test2", "null")

		err := cs.Parse()
		if err == nil {
			t.Error("no error thrown")
		}
	})
}
