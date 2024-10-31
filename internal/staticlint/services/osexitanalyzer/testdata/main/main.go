package main

import (
	"fmt"
	"os"
)

func main() {
	func2()
	fmt.Println("")
	os.Environ()
	os.Exit(100) // want "direct call of function os.Exit()"
}

func func1() {
	os.Exit(100)
}

func func2() {
	var x int
	println(x)

}
