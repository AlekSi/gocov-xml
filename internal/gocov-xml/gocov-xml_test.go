package gocov_xml

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	actual = bytes.Replace(actual, []byte(rootPath), nil, -1)
	actual = regexp.MustCompile(`timestamp="\d+"`).ReplaceAll(actual, []byte(`timestamp="123456789"`))

	expected, err := ioutil.ReadFile("expected_package1.xml")
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
