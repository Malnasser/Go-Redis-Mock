package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

func main() {
	run()
}

func run() (err error) {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal("Error establising listener to port %v", l.Addr())
		os.Exit(1)
	}

	defer closeIt(l)

	conn, err := l.Accept()
	if err != nil {
		log.Fatal("Error creating connection: ", err.Error())
		os.Exit(1)
	}

	data := make([]byte, 2046)
	for {
		n, err := conn.Read(data)
		if errors.Is(err, io.EOF) {
			break
		}

		log.Println("commend %v%s", data[:n])

		_, err = conn.Write([]byte("+PING\r\n"))
	}

	return nil
}

func closeIt(l net.Listener) {
	err := l.Close()
	if err != nil {
		log.Fatal("Fail to close application")
	}
}
