package pmml2lua

import (
	"testing"

	"github.com/kelindar/pmml2lua/schema"
	"github.com/stretchr/testify/assert"
)

func TestDecisionTree1(t *testing.T) {

	var out schema.DecisionTree
	body, global, code := scopeFor("fixtures/tree1.xml", &out)
	global.DecisionTree(out, global)
	body.With(
		NewStatement().Return().Call(out.ModelName, "v"),
	)

	assert.Contains(t, code(), `xxxx`)
}

func TestNode(t *testing.T) {
	input :=
		`<Node id="1" score="will play" recordCount="100" defaultChild="2">
		<True/>
		<ScoreDistribution value="will play" recordCount="60" confidence="0.6"/>
		<ScoreDistribution value="may play" recordCount="30" confidence="0.3"/>
		<ScoreDistribution value="no play" recordCount="10" confidence="0.1"/>
		<Node id="2" score="will play" recordCount="50" defaultChild="3">
		  <SimplePredicate field="outlook" operator="equal" value="sunny"/>
		  <ScoreDistribution value="will play" recordCount="40" confidence="0.8"/>
		  <ScoreDistribution value="may play" recordCount="2" confidence="0.04"/>
		  <ScoreDistribution value="no play" recordCount="8" confidence="0.16"/>
		</Node>
	  </Node>`

	var out schema.Node
	body, global, code := scopeFor(input, &out)
	body.Node(out, schema.DecisionTree{}, global)

	assert.Contains(t, code(), `local x = (v.outlook and v.outlook == 'sunny'); if x then`)
}
