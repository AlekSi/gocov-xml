package package1_test

import (
	. "."
	"testing"
)

func TestBuzz(t *testing.T) {
	if Buzz() != "Buzz" {
		t.Error("not a Buzz")
	}
}

func TestFizzBuzz(t *testing.T) {
	if FizzBuzz() != "FizzBuzz" {
		t.Error("not a FizzBuzz")
	}
}
