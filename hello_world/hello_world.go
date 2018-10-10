package hello_world

import "fmt"

func Run() {
	fmt.Println(VerboseAdder(3, 4))
}

func VerboseAdder(x, y int) int {
	sum := x + y
	return sum
}
