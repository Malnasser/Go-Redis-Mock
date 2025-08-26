package main

import (
	"fmt"
	"os"

	"redis-server/app/server"
)

func main() {
	server := server.NewRedisServer(":6379")
	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
