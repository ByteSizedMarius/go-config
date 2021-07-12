package go_config

import (
	"os"
	"testing"
)

type TestConfig struct {
	Option1 string `name:"optionEins"`
	Option2 string `name:"optionZwei"`
	Option3 string `name:"optionDrei"`
}

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
