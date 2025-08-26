package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type RedisServer struct {
	address string
	storage map[string]string
	exp     map[string]time.Time
	mu      sync.RWMutex
}

func NewRedisServer(address string) *RedisServer {
	return &RedisServer{
		address: address,
		storage: make(map[string]string),
		exp:     make(map[string]time.Time),
	}
}

func (s *RedisServer) Start() (err error) {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		os.Exit(1)
	}

	defer listener.Close()

	for {
		conn, errer := listener.Accept()
		if errer != nil {
			return
		}
		go s.handleConnection(conn)
	}
}

func (s *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New connection %v", conn.RemoteAddr())

	// Create a buffered reader with the sample data
	reader := resp.NewReader(conn)
	writer := resp.NewWriter(conn)

	for {
		parsedValue, err := reader.Parse()
		if err != nil {
			return
		}
		response := s.processCommand(parsedValue)
		if err := writer.WriteValue(response); err != nil {
			return
		}
	}
}

func (s *RedisServer) processCommand(value resp.Value) resp.Value {
	if value.Type != resp.Array || len(value.Array) == 0 {
		return resp.Value{Type: resp.Error, String: "Err invlid command"}
	}

	args := make([]string, len(value.Array))
	for i, arg := range value.Array {
		if arg.Type != resp.BulkString {
			return resp.Value{Type: resp.Error, String: "ERR invlid command format"}
		}
		args[i] = arg.String
	}

	command := strings.ToUpper(args[0])

	switch command {

	case "PING":
		if len(args) == 1 {
			return resp.Value{Type: resp.SimpleString, String: "PONG"}
		} else if len(args) == 2 {
			return resp.Value{Type: resp.BulkString, String: args[1]}
		}
		return resp.Value{Type: resp.Error, String: "Error wrong number of args"}

	case "ECHO":
		if len(args) != 2 {
			return resp.Value{Type: resp.Error, String: "Error wrong number of args for 'ping' commad"}
		}
		return resp.Value{Type: resp.BulkString, String: args[1]}

	case "SET":
		if len(args) < 3 {
			return resp.Value{Type: resp.Error, String: "Err wrong number of args"}
		}
		key, value := args[1], args[2]
		s.mu.Lock()
		defer s.mu.Unlock()
		s.storage[key] = value

		if len(args) == 5 && strings.ToUpper(args[3]) == "PX" {
			ms, err := strconv.Atoi(args[4])
			if err != nil || ms <= 0 {
				return resp.Value{Type: resp.Error, String: "Error invlid expire time"}
			}
			s.exp[key] = time.Now().Add(time.Duration(ms) * time.Millisecond)
		}

		return resp.Value{Type: resp.SimpleString, String: "OK"}

	case "GET":
		if len(args) != 2 {
			return resp.Value{Type: resp.Error, String: "Error wrong number of args"}
		}
		key := args[1]
		s.mu.RLock()

		val, exists := s.storage[key]
		if !exists {
			s.mu.RUnlock()
			return resp.Value{Type: resp.BulkString, String: ""}
		}

		if expireTime, hasExpiry := s.exp[key]; hasExpiry {
			if time.Now().After(expireTime) {
				s.mu.RUnlock()
				s.mu.Lock()
				delete(s.storage, key)
				delete(s.exp, key)
				s.mu.Unlock()
				return resp.Value{Type: resp.BulkString, String: ""}
			}
		}

		s.mu.RUnlock()
		return resp.Value{Type: resp.BulkString, String: val}

	default:
		return resp.Value{Type: resp.Error, String: fmt.Sprintf("Err unknown command")}
	}
}
