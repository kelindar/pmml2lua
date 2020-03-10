package pmml2lua

import (
	"context"
	"testing"

	"github.com/kelindar/lua"
	"github.com/kelindar/pmml2lua/schema"
	"github.com/stretchr/testify/assert"
)

func TestCompoundPredicate(t *testing.T) {
	input := `<CompoundPredicate booleanOperator="or">
		<CompoundPredicate booleanOperator="and">
			<SimplePredicate field="temperature" operator="lessThan" value="90"/>
			<SimplePredicate field="temperature" operator="greaterThan" value="50"/>
		</CompoundPredicate>
		<SimplePredicate field="humidity" operator="greaterOrEqual" value="80"/>
	</CompoundPredicate>`

	var out schema.Predicate
	body, global, code := scopeFor(input, &out)
	body.With(
		NewStatement().Return().CompoundPredicate(out.CompoundPredicate, global),
	)

	assert.Contains(t, code(),
		`eval.Or(eval.And(v.temperature and v.temperature < 90, v.temperature and v.temperature > 50), v.humidity and v.humidity >= 80)`,
	)

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
}

func TestSimplePredicate(t *testing.T) {
	td := []struct {
		xml    string      // The XML document to parse
		lua    string      // The output LUA code
		input  interface{} // The input data
		expect interface{} // The expected result
	}{
		{
			xml:    `<SimplePredicate field="wallet" operator="lessThan" value="0.08086312118570185"/>`,
			lua:    `v.wallet and v.wallet < 0.08086312118570185`,
			input:  map[string]float64{},
			expect: nil,
		},
		{
			xml:    `<SimplePredicate field="wallet" operator="lessThan" value="0.08086312118570185"/>`,
			lua:    `v.wallet and v.wallet < 0.08086312118570185`,
			input:  map[string]float64{"wallet": 0.05},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="wallet" operator="lessThan" value="0.08086312118570185"/>`,
			lua:    `v.wallet and v.wallet < 0.08086312118570185`,
			input:  map[string]float64{"wallet": 0.09},
			expect: false,
		},
		{
			xml:    `<SimplePredicate field="wallet" operator="lessOrEqual" value="0.08"/>`,
			lua:    `v.wallet and v.wallet <= 0.08`,
			input:  map[string]float64{"wallet": 0.08},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="name" operator="equal" value="Roman"/>`,
			lua:    `v.name and v.name == 'Roman'`,
			input:  map[string]string{"name": "Roman"},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="name" operator="notEqual" value="Wenbo"/>`,
			lua:    `v.name and v.name ~= 'Wenbo'`,
			input:  map[string]string{"name": "Roman"},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="name" operator="notEqual" value="Roman"/>`,
			lua:    `v.name and v.name ~= 'Roman'`,
			input:  map[string]string{"name": "Roman"},
			expect: false,
		},
		{
			xml:    `<SimplePredicate field="age" operator="greaterThan" value="30"/>`,
			lua:    `v.age and v.age > 30`,
			input:  map[string]float64{"age": 30},
			expect: false,
		},
		{
			xml:    `<SimplePredicate field="age" operator="greaterThan" value="30"/>`,
			lua:    `v.age and v.age > 30`,
			input:  map[string]float64{"age": 31},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="age" operator="greaterOrEqual" value="30"/>`,
			lua:    `v.age and v.age >= 30`,
			input:  map[string]float64{"age": 30},
			expect: true,
		},
	}

	for _, tt := range td {
		t.Run(tt.xml, func(t *testing.T) {

			// Generate the script
			var out schema.Predicate
			body, _, code := scopeFor(tt.xml, &out)
			body.With(
				NewStatement().Return().SimplePredicate(out.SimplePredicate),
			)

			// Code must contain the statement
			assert.Contains(t, code(), tt.lua)

			// Run the generated script
			s := makeScript(code())
			v, err := s.Run(context.Background(), tt.input)
			assert.NoError(t, err)

			var result interface{}
			switch v := v.(type) {
			case lua.Nil:
				result = nil
			case lua.Bool:
				result = bool(v)
			}
			assert.Equal(t, tt.expect, result)
		})
	}
}
