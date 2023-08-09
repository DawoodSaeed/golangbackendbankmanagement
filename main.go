package main

import "fmt"

func main() {
	fmt.Println("YEAH BUDDY")
	server := NewAPIServer(":3000")
	server.Run()
}
