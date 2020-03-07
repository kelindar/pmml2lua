package pmml2lua

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/kelindar/pmml2lua/schema"
)

var (
	errNilElement = errors.New("unable to convert a nil element")
)

// Writer represents a writer which converts PMML to LUA code.
type Writer struct {
	dst *bufio.Writer // The destination writer
	err error         // The last error that has occured
}

// New creates a new writer for conversion.
func New(dst io.Writer) *Writer {
	return &Writer{
		dst: bufio.NewWriter(dst),
	}
}

// Flush flushes the writer
func (w *Writer) Flush() error {
	return w.dst.Flush()
}

// Field generates the LUA code for the element.
func (w *Writer) Field(fieldName string) error {
	return w.Append(`v.%v`, fieldName)
}

// Value generates the LUA code for the element.
func (w *Writer) Value(v schema.Value) error {
	if err := w.Number(string(v)); err != nil {
		return w.String(string(v))
	}

	return nil
}

// Boolean writes a LUA boolean value
func (w *Writer) Boolean(v bool) (err error) {
	if v {
		return w.Append("true")
	}
	return w.Append("false")
}

// Number writes a LUA number.
func (w *Writer) Number(v interface{}) (err error) {
	switch f := v.(type) {
	case float64:
		w.Append("%v", strconv.FormatFloat(f, 'f', 24, 64))
	case string:
		if _, err = strconv.ParseFloat(f, 64); err == nil {
			err = w.Append(f)
		}
	default:
		err = fmt.Errorf("WriteNumber: unsupported type %T", f)
	}
	return
}

// String writes an escaped LUA string.
func (w *Writer) String(v string) error {
	return w.Append(`'%v'`, v)
}

// Whitespace writes an white space character
func (w *Writer) Whitespace() error {
	return w.Append(` `)
}

// Append writes a formatted string
func (w *Writer) Append(format string, args ...interface{}) error {
	_, err := w.dst.WriteString(fmt.Sprintf(format, args...))
	return err
}

// AppendLine writes a formatted string with a new line
func (w *Writer) AppendLine(format string, args ...interface{}) error {
	_, err := w.dst.WriteString(fmt.Sprintf(format+"\n", args...))
	return err
}

// Each stops on the first error
func (w *Writer) Each(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// OneOf stops on the first non-error
func (w *Writer) OneOf(errs ...error) error {
	for _, err := range errs {
		if err == nil {
			return nil
		}
	}
	return errNilElement
}

// Builtins writes builtin functions
func (w *Writer) Builtins() error {
	return w.Append(`
-- Checks if the value is missing
local function Unknown(v)
	return v == nil or v == ''
end

-- Performs a logical AND operation on the arguments
local And(...)
	local 
	for i, v in ipairs(arg) do
		if Unknown(v) then 
			return false
		end
		if not v then
			return false
		end
	end
	return true
end

`)
}
