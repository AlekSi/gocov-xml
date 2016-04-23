package package2

import (
	"fmt"

	"github.com/AlekSi/gocov-xml/internal/test/package1"
)

func FizzBuzz() {
	for i := 1; i <= 100; i++ {
		m3 := (i % 3) == 0
		m5 := (i % 5) == 0
		switch {
		case m3 && m5:
			fmt.Println(package1.FizzBuzz())
		case m3:
			fmt.Println(package1.Fizz())
		case m5:
			fmt.Println(package1.Buzz())
		default:
			fmt.Println(i)
		}
	}
}
