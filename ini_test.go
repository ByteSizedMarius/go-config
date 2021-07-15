package go_config

import (
	"os"
	"testing"
)

func TestIni(t *testing.T) {
	t.Run("TestUseInIni", func(t *testing.T) {
		os.Args = buildArgs([]string{})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetUseInINI(true)
		cs.SetINI(createIni("[MISC]\ntest1 = testxxx"))

		err := cs.Parse()
		if err != nil {
			t.Error(err)
		}

		if x.Option1 != "testxxx" {
			t.Error("value was not set correctly")
		}
	})

	t.Run("TestDoNotUseInIni", func(t *testing.T) {
		os.Args = buildArgs([]string{})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetUseInINI(false)
		cs.NewString("test2", "null")
		cs.SetINI(createIni("[MISC]\ntest1 = testxxx\ntest2 = xxx"))

		err := cs.Parse()
		if err == nil {
			t.Error("no error thrown")
		}
	})

	t.Run("TestUseAliasInIni", func(t *testing.T) {
		os.Args = buildArgs([]string{})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetAlias([]string{"testAlias1"})
		cs.NewString("test2", "null")
		cs.SetINI(createIni("[MISC]\ntestAlias1 = testxxx\ntest2 = xxx"))

		err := cs.Parse()
		if err != nil {
			t.Error(err)
		}

		if x.Option1 != "testxxx" {
			t.Error("value was not set correctly")
		}
	})

	t.Run("TestDoNotUseInIni", func(t *testing.T) {
		os.Args = buildArgs([]string{})
		x := TestUseStruct{}

		cs := initConfig(t, &x)
		cs.NewString("test1", "null").SetUseAliasInINI(false).SetAlias([]string{"testAlias1"})
		cs.NewString("test2", "null")
		cs.SetINI(createIni("[MISC]\ntestAlias1 = testxxx\ntest2 = xxx"))

		err := cs.Parse()
		if err == nil {
			t.Error("no error thrown")
		}
	})

	_ = os.Remove("test_conf.ini")
}
