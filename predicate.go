package pmml2lua

import (
	"strings"

	"github.com/kelindar/pmml2lua/schema"
)

// ----------------------------------------------------------------------------

// Predicate generates the LUA code for the element.
func (s *Scope) Predicate(v *schema.Predicate, global *Scope) Compiler {
	switch {
	case v.SimplePredicate != nil:
		return NewStatement().SimplePredicate(*v.SimplePredicate)
	case v.CompoundPredicate != nil:
		return NewStatement().CompoundPredicate(*v.CompoundPredicate, global)
	case v.SimpleSetPredicate != nil:
		return NewStatement().SimpleSetPredicate(*v.SimpleSetPredicate, global)
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
func (s *Statement) CompoundPredicate(v schema.CompoundPredicate, global *Scope) *Statement {
	s.Append("eval.%s({", strings.Title(v.Operator))
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
	s.Append("; n=%d})", len(v.Predicates))
	return s
}

// ----------------------------------------------------------------------------

// BinaryOperator generates the LUA code for the element.
func (s *Statement) BinaryOperator(v string) *Statement {
	var operator string
	switch v {
	case "equal":
		operator = " == "
	case "notEqual":
		operator = " ~= "
	case "lessThan":
		operator = " < "
	case "lessOrEqual":
		operator = " <= "
	case "greaterThan":
		operator = " > "
	case "greaterOrEqual":
		operator = " >= "
	case "isMissing":
		operator = " == nil "
	case "isNotMissing":
		operator = " ~= nil "
	default:
		return s.Error("binary operator %v is not supported", v)
	}

	// Write the operator
	return s.Append(operator)
}

// SimplePredicate generates the LUA code for the element.
func (s *Statement) SimplePredicate(v schema.SimplePredicate) *Statement {
	if v.Operator == "isMissing" || v.Operator == "isNotMissing" {
		return s.Field(v.Field).BinaryOperator(v.Operator)
	}

	return s.Field(v.Field).
		Append(" and ").
		Field(v.Field).
		BinaryOperator(v.Operator).
		Value(v.Value)
}

// ----------------------------------------------------------------------------

// SimpleSetPredicate generates the LUA code for the element.
func (s *Statement) SimpleSetPredicate(v schema.SimpleSetPredicate, global *Scope) *Statement {
	if v.Array == nil {
		return s.Error("array must not be nil")
	}

	values, err := v.Array.CSV()
	if err != nil {
		return s.Error(err.Error())
	}

	/*_ = global.With(
		NewStatement().Append("local array = {%s; n=%d}", values, v.Array.Length),
	)
	return s.Append("eval.%s(", strings.Title(v.Operator)).
		Field(v.Field).
		Append(", array)")*/

	return s.Append("eval.%s(", strings.Title(v.Operator)).
		Field(v.Field).
		Append(", {%s; n=%d})", values, v.Array.Length)
}
