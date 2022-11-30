package main

import (
	"fmt"
	"os"
)

func sqrt(i int) int {
	return i * 2
}

func main() {

	fmt.Printf("2 ^ 2 = %d", sqrt(2))
	os.Exit(0) // want "using os.Exit in main func of main package"
}
