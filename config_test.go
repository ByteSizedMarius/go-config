package go_config

import (
	"io/ioutil"
	"os"
	"testing"
)

//
//
// Util
//
//

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

func createIni(val string) string {
	file := "test_conf.ini"
	err := ioutil.WriteFile(file, []byte(val), 0777)
	if err != nil {
		panic(err)
	}
	return file
}

//
//
// Tests
//
//

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
