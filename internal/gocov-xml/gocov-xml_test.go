package gocov_xml

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

var rootPath string

func init() {
	var err error
	rootPath, err = filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		panic(err)
	}
}

func TestPackage1(t *testing.T) {
	actual, err := ioutil.ReadFile(filepath.Join(rootPath, "package1.xml"))
	if err != nil {
		t.Fatal(err)
	}
	actual = bytes.Replace(actual, []byte(rootPath), nil, -1)
	actual = regexp.MustCompile(`timestamp="\d+"`).ReplaceAll(actual, []byte(`timestamp="123456789"`))
	t.Logf("actual:\n%s", actual)

	expected, err := ioutil.ReadFile("expected_package1.xml")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("expected:\n%s", expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Error("expected != actual")
	}
}
