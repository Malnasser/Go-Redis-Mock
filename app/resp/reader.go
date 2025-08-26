package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func (r *Reader) Parse() (value Value, err error) {
	typeByte, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeByte {
	case SimpleString:
		return r.parseSimpleString()
	case Array:
		return r.parseArray()
	case BulkString:
		return r.parseBulkString()
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %c", typeByte)
	}
}

func (r *Reader) parseBulkString() (value Value, err error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	length, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, err
	}
	if length == -1 {
		return Value{Type: BulkString, String: ""}, nil
	}
	bulk := make([]byte, length)
	if _, err := io.ReadFull(r.reader, bulk); err != nil {
		return Value{}, err
	}
	if _, err := r.readLine(); err != nil {
		return Value{}, err
	}
	return Value{Type: BulkString, String: string(bulk)}, nil
}

func (r *Reader) parseSimpleString() (value Value, err error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: SimpleString, String: string(line)}, nil
}

func (r *Reader) parseArray() (value Value, err error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	length, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, err
	}
	if length == -1 {
		return Value{Type: Array, Array: nil}, nil
	}

	array := make([]Value, length)
	for i := range length {
		val, err := r.Parse()
		if err != nil {
			return Value{}, err
		}
		array[i] = val
	}

	return Value{Type: Array, Array: array}, nil
}

func (r *Reader) readLine() (bytes []byte, err error) {
	line, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(line) >= 2 && line[len(line)-2] == '\r' {
		return line[:len(line)-2], nil
	}
	return line, nil
}
