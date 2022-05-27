package main

import "fmt"

func main() {

	server := NewServer("10.88.48.131", 8988)
	server.Start()

	fmt.Println("123")
}
