package go_config

import (
	"os"
	"testing"
)

func initConfig(t *testing.T, x interface{}) Config {
	cs, err := Initialize(x)
	if err != nil {
		t.Error(err)
	}
	return cs
}

func buildArgs(targs []string) []string {
	args := []string{os.Args[0]}
	for _, v := range targs {
		args = append(args, v)
	}
	return args
}

type TestCaseSensitivityStruct struct {
	Option1 string `name:"a"`
	Option2 string `name:"A"`
}

func TestCaseSensitivity(t *testing.T) {
	os.Args = buildArgs([]string{"-a=test1", "-A=test2"})
	x := TestCaseSensitivityStruct{}

	cs := initConfig(t, &x)
	cs.NewString("a", "null")
	cs.NewString("A", "null")

	err := cs.Parse()
	if err != nil {
		t.Error(err)
	}

	if !(x.Option1 == "test1" && x.Option2 == "test2") {
		t.Errorf("Output doesn't match:\nExpected: %v\nActual:   %v", "test1, test2", x.Option1+", "+x.Option2)
	}
}

type TestDuplicateStringStruct struct {
	Option1 string `name:"myTest"`
	Option2 string `name:"A"`
}

func TestDuplicateStringFlags(t *testing.T) {
	cs := initConfig(t, &TestDuplicateStringStruct{})

	defer func() {
		if errX := recover(); errX == nil {
			t.Error("Duplicate flag was accepted")
		}
	}()

	cs.NewString("myTest", "null")
	cs.NewString("A", "null")
	cs.NewString("myTest", "null")
}

type TestUseStruct struct {
	Option1 string `name:"test1"`
	Option2 string `name:"test2"`
	Option3 string `name:"test3"`
}

func TestUseInCli(t *testing.T) {
	os.Args = buildArgs([]string{"-test1=x1"})
	x := TestUseStruct{}

	cs := initConfig(t, &x)
	cs.NewString("test1", "null")

	err := cs.Parse()
	if err != nil {
		t.Error(err)
	}
}

func TestDoNotUseInCli(t *testing.T) {
	os.Args = buildArgs([]string{"-test1=x1", "-test2=x2"})
	x := TestUseStruct{}

	cs := initConfig(t, &x)
	cs.NewString("test1", "null").SetUseInCli(false)
	cs.NewString("test2", "null")

	err := cs.Parse()
	if err == nil {
		t.Error("no error thrown")
	}
}

func TestUseAliasInCli(t *testing.T) {
	os.Args = buildArgs([]string{"-testAlias1=x1", "-test2=x2"})
	x := TestUseStruct{}

	cs := initConfig(t, &x)
	cs.NewString("test1", "null").SetAlias([]string{"testAlias1"})
	cs.NewString("test2", "null")

	err := cs.Parse()
	if err != nil {
		t.Error(err)
	}
}

func TestDoNotUseAliasInCli(t *testing.T) {
	os.Args = buildArgs([]string{"-testAlias1=x1", "-test2=x2"})
	x := TestUseStruct{}

	cs := initConfig(t, &x)
	cs.NewString("test1", "null").SetUseAliasInCli(false).SetAlias([]string{"testAlias1"})
	cs.NewString("test2", "null")

	err := cs.Parse()
	if err == nil {
		t.Error("no error thrown")
	}
}
