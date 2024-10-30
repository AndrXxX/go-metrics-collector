package main

import (
	"os"
)

func main() {
	os.Exit(100) // want "direct call of function os.Exit()"
}

func func1() {
	os.Exit(100)
}
