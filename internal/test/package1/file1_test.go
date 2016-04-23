package package1

import (
	"testing"
)

func TestFizz(t *testing.T) {
	if Fizz() != "Fizz" {
		t.Error("not a Fizz")
	}
}
