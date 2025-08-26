package resp

// RESP data type constant
const (
	SimpleString = '+'
	Error        = '-'
	integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

// Value represents a RESP protocol Value
type Value struct {
	Type   byte
	String string
	Int    int64
	Array  []Value
}

func (v Value) Stringify() string {
	switch v.Type {
	case SimpleString:
		return "Simeple string: " + v.String
	case Error:
		return "Error: " + v.String
	case integer:
		return "integer: " + string(rune(v.Int))
	case BulkString:
		if len(v.String) == -1 {
			return "BulkString: null"
		}
		return "BulkString: " + v.String
	case Array:
		if len(v.Array) == -1 {
			return "Array: null"
		}
		result := "Array["
		for i, val := range v.Array {
			if i > 0 {
				result += ", "
			}
			result += val.Stringify()
		}
		result += "]"
		return result
	default:
		return "Unknown type"
	}
}
