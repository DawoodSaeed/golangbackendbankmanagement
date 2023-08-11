package main

import (
	"fmt"
	"log"
)

func main() {
	// Database Connection;
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal("Connection to the database was failed")
	}

	fmt.Printf("%+v\n", store)

	// Server;
	//store.Init()
	server := NewAPIServer(":3000", store)
	server.Run()
}
