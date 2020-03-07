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
	w, r := writerFor(input, &out)
	assert.NoError(t, w.CompoundPredicate(out.CompoundPredicate))

	println(r())
	assert.Fail(t, "x")
}

func TestSimplePredicate(t *testing.T) {
	td := []struct {
		xml    string      // The XML document to parse
		lua    string      // The output LUA code
		input  interface{} // The input data
		expect bool        // The expected result
	}{
		{
			xml:    `<SimplePredicate field="wallet" operator="lessThan" value="0.08086312118570185"/>`,
			lua:    `v.wallet < 0.08086312118570185`,
			input:  map[string]float64{"wallet": 0.05},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="wallet" operator="lessThan" value="0.08086312118570185"/>`,
			lua:    `v.wallet < 0.08086312118570185`,
			input:  map[string]float64{"wallet": 0.09},
			expect: false,
		},
		{
			xml:    `<SimplePredicate field="wallet" operator="lessOrEqual" value="0.08"/>`,
			lua:    `v.wallet <= 0.08`,
			input:  map[string]float64{"wallet": 0.08},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="name" operator="equal" value="Roman"/>`,
			lua:    `v.name == 'Roman'`,
			input:  map[string]string{"name": "Roman"},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="name" operator="notEqual" value="Wenbo"/>`,
			lua:    `v.name ~= 'Wenbo'`,
			input:  map[string]string{"name": "Roman"},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="name" operator="notEqual" value="Roman"/>`,
			lua:    `v.name ~= 'Roman'`,
			input:  map[string]string{"name": "Roman"},
			expect: false,
		},
		{
			xml:    `<SimplePredicate field="age" operator="greaterThan" value="30"/>`,
			lua:    `v.age > 30`,
			input:  map[string]float64{"age": 30},
			expect: false,
		},
		{
			xml:    `<SimplePredicate field="age" operator="greaterThan" value="30"/>`,
			lua:    `v.age > 30`,
			input:  map[string]float64{"age": 31},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="age" operator="greaterOrEqual" value="30"/>`,
			lua:    `v.age >= 30`,
			input:  map[string]float64{"age": 30},
			expect: true,
		},
	}

	for _, tt := range td {
		t.Run(tt.xml, func(t *testing.T) {

			// Generate the script
			var out schema.Predicate
			w, r := writerFor(tt.xml, &out)
			assert.NoError(t, w.SimplePredicate(out.SimplePredicate))
			assert.EqualValues(t, tt.lua, r())

			// Run the generated script
			s := makeScript("return %v", r())
			v, err := s.Run(context.Background(), tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, bool(v.(lua.Bool)))
		})
	}
}
