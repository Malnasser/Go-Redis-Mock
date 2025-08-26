package command

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func ProcessCommand(value resp.Value) resp.Value {
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

	default:
		return resp.Value{Type: resp.Error, String: fmt.Sprintf("Err unkown command")}
	}
}
