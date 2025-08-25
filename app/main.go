package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

func main() {
	if err := run(); err != nil {
		log.Fatal("Error starting the server: ", err.Error())
	}
}

func run() (err error) {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		os.Exit(1)
	}

	log.Printf("Server start listening on port: ", l.Addr())

	defer closeIt(l)

	var wg sync.WaitGroup

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error creating connection: ", err.Error())
			os.Exit(1)
		}
		wg.Add(1)
		go handleConnection(conn, &wg)
	}

	return nil
}

func closeIt(l net.Listener) {
	err := l.Close()
	if err != nil {
		log.Fatal("Fail to close application")
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()
	data := make([]byte, 2046)
	for {
		n, err := conn.Read(data)
		if errors.Is(err, io.EOF) {
			break
		}

		log.Printf("New read with length %v", n)

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			break
		}
	}
}
