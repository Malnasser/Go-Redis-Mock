package server

import (
	"log"
	"net"
	"os"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func RunRedisServer() (err error) {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		os.Exit(1)
	}

	log.Printf("Server start listening on port: %v ", l.Addr())

	defer func() {
		if err := l.Close(); err != nil {
			log.Fatal("Error closing socket: ", err.Error())
		}
	}()

	var connectionWG sync.WaitGroup
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		connectionWG.Add(1)
		go handleConnection(conn, &connectionWG)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	// Create a buffered reader with the sample data
	reader := resp.NewReader(conn)
	writer := resp.NewWriter(conn)

	for {
		parsedValue, err := reader.Parse()
		if err != nil {
			return
		}
		response := command.ProcessCommand(parsedValue)
		if err := writer.WriteValue(response); err != nil {
			return
		}
	}
}
