package pmml2lua

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kelindar/pmml2lua/schema"
)

var (
	errNilElement = errors.New("unable to convert a nil element")
)

// Compiler represents somthing that can be compiled (statement or scope)
type Compiler interface {
	Compile() ([]byte, error)
}

// Assert various compiler implementations
var (
	_ Compiler = new(Statement)
	_ Compiler = new(Scope)
)

// ----------------------------------------------------------------------------

// Scope represents a scope that can be rendered.
type Scope struct {
	ref string     // The reference of the scope (e.g. name of the function)
	dst []Compiler // The list of statements
	tab int        // The number of tabs for indentation
}

// NewScope prepares a new scope.
func NewScope() *Scope {
	return &Scope{
		dst: make([]Compiler, 0, 8),
	}
}

// Compile returns the compiled statement.
func (s *Scope) Compile() ([]byte, error) {
	var buf bytes.Buffer
	for _, v := range s.dst {
		switch v := v.(type) {
		case *Statement:
			v.tab = s.tab
		case *Scope:
			v.tab = s.tab + 1
		}

		compiled, err := v.Compile()
		if err != nil {
			return nil, err
		}
		buf.Write(compiled)
	}

	return buf.Bytes(), nil
}

// Name returns the name of the scope
func (s *Scope) Name() string {
	return s.ref
}

// With adds the children to the scope.
func (s *Scope) With(body ...Compiler) *Scope {
	return s.WithIf(true, body...)
}

// WithIf conditionally adds the children to the scope.
func (s *Scope) WithIf(condition bool, body ...Compiler) *Scope {
	if condition {
		s.dst = append(s.dst, body...)
	}
	return s
}

// Scope creates a child scope without sub-identation
func (s *Scope) Scope() *Scope {
	child := NewScope()
	s.With(child)
	child.tab = s.tab // Make sure we're at the same level
	return child
}

// Function creates a LUA function scope
func (s *Scope) Function(name string, args ...string) *Scope {
	body := NewScope()
	body.ref = name
	s.With(
		NewStatement().Append("function %s(%s)", name, strings.Join(args, ", ")),
		body,
		NewStatement().Append("end"),
	)
	return body
}

// ----------------------------------------------------------------------------

// Statement represents a single line statement
type Statement struct {
	buf  *bytes.Buffer // The destination writer
	tab  int           // The number of tabs for indentation
	err  error         // The last error that has occured
	cond bool          // The condition to evaluate for the statement
}

// NewStatement prepares a new statement.
func NewStatement() *Statement {
	return &Statement{
		buf:  bytes.NewBuffer(nil),
		cond: true,
	}
}

// Append creates a new formatting statement.
func Append(format string, args ...interface{}) *Statement {
	return AppendIf(true, format, args...)
}

// AppendIf creates a new conditional statement.
func AppendIf(condition bool, format string, args ...interface{}) *Statement {
	return (&Statement{
		buf:  bytes.NewBuffer(nil),
		cond: condition,
	}).Append(format, args...)
}

// Compile returns the compiled statement.
func (s *Statement) Compile() ([]byte, error) {
	if s.err != nil {
		return nil, fmt.Errorf("statement: error at %s... due to %s",
			s.err.Error(),
			string(s.buf.Bytes()),
		)
	}

	// If the statement is disabled, do not generate anything
	if !s.cond {
		return nil, nil
	}

	// Writte identation to the buffer
	var buffer bytes.Buffer
	for i := 0; i < s.tab; i++ {
		buffer.WriteRune('\t')
	}
	buffer.Write(s.buf.Bytes())
	buffer.WriteRune('\n')
	return buffer.Bytes(), nil
}

// Statement appends one statement.
func (s *Statement) Statement(statement *Statement) *Statement {
	if s.err != nil {
		return s
	}

	_, s.err = s.buf.Write(statement.buf.Bytes())
	return s
}

// Append writes a formatted string.
func (s *Statement) Append(format string, args ...interface{}) *Statement {
	if s.err == nil {
		_, s.err = s.buf.WriteString(fmt.Sprintf(format, args...))
	}
	return s
}

// String writes an escaped LUA string.
func (s *Statement) String(v string) *Statement {
	return s.Append(`'%v'`, v)
}

// Whitespace writes an white space character.
func (s *Statement) Whitespace() *Statement {
	return s.Append(` `)
}

// Value generates the LUA code for the element.
func (s *Statement) Value(v schema.Value) *Statement {
	if s.err != nil {
		return s
	}

	// Try to write the number first
	s.Number(string(v))
	if s.err != nil {
		s.err = nil
		return s.String(string(v))
	}

	return s
}

// Boolean writes a LUA boolean value
func (s *Statement) Boolean(v bool) *Statement {
	if v {
		return s.Append("true")
	}
	return s.Append("false")
}

// Number writes a LUA number.
func (s *Statement) Number(v interface{}) *Statement {
	switch f := v.(type) {
	case float64:
		return s.Append("%v", strconv.FormatFloat(f, 'f', 24, 64))
	case string:
		if _, s.err = strconv.ParseFloat(f, 64); s.err == nil {
			return s.Append(f)
		}
	default:
		s.err = fmt.Errorf("WriteNumber: unsupported type %T", f)
	}
	return s
}

// Field generates the LUA code for the element.
func (s *Statement) Field(fieldName string) *Statement {
	return s.Append(`v.%v`, fieldName)
}

// Return writes a return keyword.
func (s *Statement) Return() *Statement {
	return s.Append("return ")
}

// Call writes a function call statement.
func (s *Statement) Call(name string, args ...string) *Statement {
	return s.Append("%s(%s)", name, strings.Join(args, ", "))
}

// Error sets the erroor internally and returns the statement
func (s *Statement) Error(format string, args ...interface{}) *Statement {
	if s.err == nil {
		s.err = fmt.Errorf(format, args...)
	}
	return s
}
