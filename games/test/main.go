package main

import "fmt"

func main() {}

//export add
func add(a, b int32) int32 {
	return a + b
}

//export greet
func greet(name string) string {
	return fmt.Sprintf("Hello %s from WASM!", name)
}
