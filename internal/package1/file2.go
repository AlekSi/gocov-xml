package package1

import (
	"fmt"
)

func Buzz() string {
	return "Buzz"
}

func FizzBuzz() string {
	return fmt.Sprintf("%s%s", Fizz(), Buzz())
}
