package schema

import (
	"encoding/xml"
	"testing"

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

	var out CompoundPredicate
	assert.NoError(t, xml.Unmarshal([]byte(input), &out))
	assert.EqualValues(t, CompoundPredicate{
		Operator: "or",
		Predicates: []Predicate{{
			CompoundPredicate: &CompoundPredicate{
				Operator: "and",
				Predicates: []Predicate{{
					SimplePredicate: &SimplePredicate{
						Field: "temperature", Operator: "lessThan", Value: "90",
					},
				}, {
					SimplePredicate: &SimplePredicate{
						Field: "temperature", Operator: "greaterThan", Value: "50",
					},
				}},
			},
		}, {
			SimplePredicate: &SimplePredicate{
				Field: "humidity", Operator: "greaterOrEqual", Value: "80",
			},
		}},
	}, out)
}
