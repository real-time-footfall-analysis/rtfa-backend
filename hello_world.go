package main

import "fmt"

func main() {
    fmt.Println(verboseAdder(3,4))
}

func verboseAdder(x, y int) (int) {
	sum := x + y
	return sum
}