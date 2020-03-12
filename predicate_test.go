package pmml2lua

import (
	"context"
	"testing"

	"github.com/kelindar/lua"
	"github.com/kelindar/pmml2lua/schema"
	"github.com/stretchr/testify/assert"
)

// Benchmark_Predicate/compound-8         	  139858	      8533 ns/op	    3152 B/op	      72 allocs/op
func Benchmark_Predicate(b *testing.B) {
	input :=
		`<CompoundPredicate booleanOperator="or">
	<CompoundPredicate booleanOperator="and">
		<SimplePredicate field="temperature" operator="lessThan" value="90"/>
		<SimplePredicate field="temperature" operator="greaterThan" value="50"/>
	</CompoundPredicate>
	<SimplePredicate field="humidity" operator="greaterOrEqual" value="80"/>
	<SimpleSetPredicate field="humidity" booleanOperator="isNotIn">
		<Array n="5" type="int">1 2 3 4 5</Array>
	</SimpleSetPredicate>
</CompoundPredicate>`

	var out schema.Predicate
	body, global, code := scopeFor(input, &out)
	body.With(
		NewStatement().Return().CompoundPredicate(*out.CompoundPredicate, global),
	)

	s := makeScript(code())
	b.Run("compound", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Run(context.Background(), map[string]int{
				"temperature": 60,
			})
		}
	})
}

func TestSimpleSetPredicate(t *testing.T) {
	input :=
		`<SimpleSetPredicate field="value" booleanOperator="isNotIn">
			<Array n="10" type="int">1 2 3 4 5 10 11 12 13 15</Array>
		</SimpleSetPredicate>`

	var out schema.Predicate
	body, global, code := scopeFor(input, &out)
	body.With(
		NewStatement().Return().SimpleSetPredicate(*out.SimpleSetPredicate, global),
	)

	assert.Contains(t, code(),
		`tree.IsNotIn(v.value, {1,2,3,4,5,10,11,12,13,15; n=10})`,
	)

	s := makeScript(code())
	v, err := s.Run(context.Background(), map[string]int{
		"value": 7,
	})
	assert.NoError(t, err)
	assert.Equal(t, "true", v.String())
}

func TestSurrogate(t *testing.T) {
	input :=
		`<CompoundPredicate booleanOperator="surrogate">
		<CompoundPredicate booleanOperator="and">
		  <SimplePredicate field="temperature" operator="lessThan" value="90"/>
		  <SimplePredicate field="temperature" operator="greaterThan" value="50"/>
		</CompoundPredicate>
		<SimplePredicate field="humidity" operator="greaterOrEqual" value="80"/>
		<False/>
	  </CompoundPredicate>`

	var out schema.Predicate
	body, global, code := scopeFor(input, &out)
	body.With(
		NewStatement().Return().CompoundPredicate(*out.CompoundPredicate, global),
	)

	assert.Contains(t, code(),
		`tree.Surrogate({tree.And({v.temperature and v.temperature < 90, v.temperature and v.temperature > 50; n=2}), v.humidity and v.humidity >= 80, false; n=3})`,
	)

	s := makeScript(code())
	v, err := s.Run(context.Background(), map[string]int{
		"humidity": 60,
	})
	assert.NoError(t, err)
	assert.Equal(t, "false", v.String())
}

func TestCompoundPredicate(t *testing.T) {
	input :=
		`<CompoundPredicate booleanOperator="or">
		<CompoundPredicate booleanOperator="and">
			<SimplePredicate field="temperature" operator="lessThan" value="90"/>
			<SimplePredicate field="temperature" operator="greaterThan" value="50"/>
		</CompoundPredicate>
		<SimplePredicate field="humidity" operator="greaterOrEqual" value="80"/>
		<SimpleSetPredicate field="humidity" booleanOperator="isNotIn">
			<Array n="5" type="int">1 2 3 4 5</Array>
		</SimpleSetPredicate>
	</CompoundPredicate>`

	var out schema.Predicate
	body, global, code := scopeFor(input, &out)
	body.With(
		NewStatement().Return().CompoundPredicate(*out.CompoundPredicate, global),
	)

	assert.Contains(t, code(),
		`tree.Or({tree.And({v.temperature and v.temperature < 90, v.temperature and v.temperature > 50; n=2}), v.humidity and v.humidity >= 80, tree.IsNotIn(v.humidity, {1,2,3,4,5; n=5}); n=3})`,
	)

	s := makeScript(code())
	v, err := s.Run(context.Background(), map[string]int{
		"temperature": 60,
	})
	assert.NoError(t, err)
	assert.Equal(t, "true", v.String())
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
		{
			xml:    `<SimplePredicate field="age" operator="isMissing" />`,
			lua:    `v.age == nil `,
			input:  map[string]int{},
			expect: true,
		},
		{
			xml:    `<SimplePredicate field="age" operator="isNotMissing" />`,
			lua:    `v.age ~= nil `,
			input:  map[string]int{"age": 12},
			expect: true,
		},
	}

	for _, tt := range td {
		t.Run(tt.xml, func(t *testing.T) {

			// Generate the script
			var out schema.Predicate
			body, _, code := scopeFor(tt.xml, &out)
			body.With(
				NewStatement().Return().SimplePredicate(*out.SimplePredicate),
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
