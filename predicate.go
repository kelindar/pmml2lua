package pmml2lua

import (
	"fmt"

	"github.com/kelindar/pmml2lua/schema"
)

// Predicate generates the LUA code for the element.
func (w *Writer) Predicate(v *schema.Predicate) error {
	return w.OneOf(
		w.SimplePredicate(v.SimplePredicate),
		w.CompoundPredicate(v.CompoundPredicate),
		w.True(v.True),
		w.False(v.False),
	)
}

// CompoundPredicate generates the LUA code for the element.
func (w *Writer) CompoundPredicate(v *schema.CompoundPredicate) error {
	if v == nil {
		return errNilElement
	}

	for i, p := range v.Predicates {
		if err := w.Each(
			w.Append("("),
			w.Predicate(&p),
			w.Append(")"),
		); err != nil {
			return err
		}

		if i+1 < len(v.Predicates) {
			w.Append(" %v ", v.Operator)
		}
	}

	// local foo = function() print(%a) end

	// Generate child predicates
	/*w.Append("%v( ", strings.Title(v.Operator))
	for i, p := range v.Predicates {
		if err := w.Predicate(&p); err != nil {
			return err
		}

		if i+1 < len(v.Predicates) {
			w.Append(", ")
		}
	}
	w.Append(" )")*/

	/*
		http://dmg.org/pmml/v4-1/TreeModel.html

		P       Q       AND     OR      XOR
		True	True	True	True	False
		True	False	False	True	True
		True	Unknown	Unknown	True	Unknown
		False	True	False	True	True
		False	False	False	False	False
		False	Unknown	False	Unknown	Unknown
		Unknown	True	Unknown	True	Unknown
		Unknown	False	False	Unknown	Unknown
		Unknown	Unknown	Unknown	Unknown	Unknown

	*/

	/*var operator string
	switch v.Operator {
	case "or":
		operator = "or"
	case "and":
		operator = "and"
	case "xor":
		operator = "xor"
	case "surrogate":
		operator = "surrogate"
	default:
		return fmt.Errorf("compound operator %v is not supported", v)
	}*/

	return nil
}

/*
   <xs:enumeration value="or"/>
   <xs:enumeration value="and"/>
   <xs:enumeration value="xor"/>
   <xs:enumeration value="surrogate"/>
*/

// SimplePredicate generates the LUA code for the element.
func (w *Writer) SimplePredicate(v *schema.SimplePredicate) error {
	if v == nil {
		return errNilElement
	}

	return w.Each(
		w.Field(v.Field),
		w.Append(" and "),
		w.Field(v.Field),
		w.Whitespace(),
		w.BinaryOperator(v.Operator),
		w.Whitespace(),
		w.Value(v.Value),
		w.Append(""),
	)
}

// BinaryOperator generates the LUA code for the element.
func (w *Writer) BinaryOperator(v string) error {
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
		return fmt.Errorf("binary operator %v is not supported", v)
	}

	// Write the operator
	return w.Append(operator)
}

// True writes a LUA boolean value
func (w *Writer) True(v *schema.True) error {
	if v == nil {
		return errNilElement
	}

	return w.Append("true")
}

// False writes a LUA boolean value
func (w *Writer) False(v *schema.False) error {
	if v == nil {
		return errNilElement
	}

	return w.Append("false")
}
