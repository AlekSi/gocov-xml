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

func TestConvertGocov(t *testing.T) {
	for _, f := range []string{"package1", "package12"} {
		actual, err := ioutil.ReadFile(filepath.Join(rootPath, f+".gocov.xml"))
		require.NoError(t, err)
		actual = bytes.Replace(actual, []byte(rootPath), nil, -1)
		actual = regexp.MustCompile(`timestamp="\d+"`).ReplaceAll(actual, []byte(`timestamp="123456789"`))

		expected, err := ioutil.ReadFile("expected_" + f + ".xml")
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}
