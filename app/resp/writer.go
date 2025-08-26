package resp

import (
	"fmt"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) WriteValue(value Value) error {
	switch value.Type {
	case SimpleString:
		return w.writeSimpleString(value.String)
	case Error:
		return w.writeError(value.String)
	case integer:
		return w.writeInteger(value.Int)
	case BulkString:
		if len(value.String) == -1 {
			return w.writeNull()
		}
		return w.writeBulkString(value.String)
	case Array:
		return w.writeArray(value.Array)
	default:
		return fmt.Errorf("unknown RESP type: %c", value.Type)
	}
}

func (w *Writer) writeSimpleString(s string) error {
	_, err := fmt.Fprintf(w.writer, "+%s\r\n", s)
	return err
}

func (w *Writer) writeError(msg string) error {
	_, err := fmt.Fprintf(w.writer, "-%s\r\n", msg)
	return err
}

func (w *Writer) writeInteger(i int64) error {
	_, err := fmt.Fprintf(w.writer, ":%d\r\n", i)
	return err
}

func (w *Writer) writeBulkString(s string) error {
	_, err := fmt.Fprintf(w.writer, "$%d\r\n%s\r\n", len(s), s)
	return err
}

func (w *Writer) writeNull() error {
	_, err := fmt.Fprintf(w.writer, "$-1\r\n")
	return err
}

func (w *Writer) writeArray(values []Value) error {
	if _, err := fmt.Fprintf(w.writer, "*%d\r\n", len(values)); err != nil {
		return err
	}

	for _, value := range values {
		if err := w.WriteValue(value); err != nil {
			return err
		}
	}
	return nil
}
