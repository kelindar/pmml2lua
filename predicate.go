package pmml2lua

import (
	"fmt"
	"strings"

	"github.com/kelindar/pmml2lua/schema"
)

// ----------------------------------------------------------------------------

// Predicate generates the LUA code for the element.
func (s *Scope) Predicate(v *schema.Predicate, global *Scope) Compiler {
	switch {
	case v.SimplePredicate != nil:
		return NewStatement().SimplePredicate(v.SimplePredicate)
	case v.CompoundPredicate != nil:
		return NewStatement().CompoundPredicate(v.CompoundPredicate, global)

	//case v.SimpleSetPredicate != nil:
	case v.True != nil:
		return NewStatement().Boolean(true)
	case v.False != nil:
		return NewStatement().Boolean(false)
	}

	return NewStatement()

	//id := "func_" + xid.New().String()
	//fn := s.Function(id, "v")
	//return fn
}

// CompoundPredicate generates the LUA code for the element.
func (s *Statement) CompoundPredicate(v *schema.CompoundPredicate, global *Scope) *Statement {
	s.Append("eval.%s(", strings.Title(v.Operator))
	for i, p := range v.Predicates {
		switch fn := global.Predicate(&p, global).(type) {
		case *Scope:
			s.Call(fn.Name(), "v")
		case *Statement:
			s.Statement(fn)
		}

		if i+1 < len(v.Predicates) {
			s.Append(", ")
		}
	}
	s.Append(")")
	return s
}

// BinaryOperator generates the LUA code for the element.
func (s *Statement) BinaryOperator(v string) *Statement {
	var operator string
	switch v {
	case "equal":
		operator = "=="
	case "notEqual":
		operator = "~="
	case "lessThan":
		operator = "<"
	case "lessOrEqual":
		operator = "<="
	case "greaterThan":
		operator = ">"
	case "greaterOrEqual":
		operator = ">="

	// TODO: support the two other operators
	case "isMissing":
		fallthrough
	case "isNotMissing":
		fallthrough
	default:
		s.err = fmt.Errorf("binary operator %v is not supported", v)
	}

	// Write the operator
	return s.Append(operator)
}

// SimplePredicate generates the LUA code for the element.
func (s *Statement) SimplePredicate(v *schema.SimplePredicate) *Statement {
	if v == nil {
		s.err = errNilElement
		return s
	}

	return s.
		Field(v.Field).
		Append(" and ").
		Field(v.Field).
		Whitespace().
		BinaryOperator(v.Operator).
		Whitespace().
		Value(v.Value)
}
