package pmml2lua

import (
	"github.com/kelindar/pmml2lua/schema"
)

// DecisionTree generates the LUA code for the element.
func (s *Scope) DecisionTree(v schema.DecisionTree, global *Scope) *Scope {
	return s.Function(v.ModelName, "v").
		With(
			Append("model = model or {}"),
			Append("model.%s = model.%s or ", v.ModelName, v.ModelName),
			NewScope().Node(v.Node, v, global),
			Append("return tree.NewTree('%s', model.%s).next(v)", v.MissingValueStrategy, v.ModelName),
		)
}

// Node generates the LUA code for the element.
func (s *Scope) Node(v schema.Node, tree schema.DecisionTree, global *Scope) *Scope {
	hasChildren := len(v.Nodes) > 0
	hasMissingStrategy := (tree.MissingValueStrategy != "none" && tree.MissingValueStrategy != "")
	masMissingEstimate := tree.MissingValueStrategy == "weightedConfidence" ||
		tree.MissingValueStrategy == "aggregateNodes" ||
		tree.MissingValueStrategy == "defaultChild"

	// Generate the run content
	inner := NewScope().With(
		Append("t.last = '%v'", v.Score),
		AppendIf(!hasChildren, "return true"), // Stop
	)

	// Missing values aggregation
	missing := NewScope()
	for _, d := range v.Distributions {
		missing.WithIf(masMissingEstimate,
			Append(`t.miss('%s', %d, %v)`, d.Value, d.RecordCount, d.Confidence),
		)
	}

	missing.WithIf(tree.MissingValueStrategy == "defaultChild" && v.DefaultChild != "",
		Append(`return n.children["%s"].eval(t, n, v)`, v.DefaultChild),
	)

	missing.WithIf(tree.MissingValueStrategy == "nullPrediction",
		Append("t.last = nil"), // Set last to UNKNOWN
		Append("return true"),  // Stop
	)

	missing.WithIf(tree.MissingValueStrategy == "lastPrediction",
		Append("return true"), // Stop
	)

	// Generate children nodes
	children := NewScope()
	for _, child := range v.Nodes {
		children.Node(child, tree, global)
	}

	return s.With(
		Append("tree.NewNode('%v', function(t, n, v) ", v.ID),
		NewScope().With(
			Append("local x = (").Predicate(v.Predicate, global).Append("); if x then"),
			inner,
			AppendIf(hasMissingStrategy, "elseif x == nil then"),
			missing,
			Append("end"),
			Append("return false"), // Continue
		),
		AppendIf(!hasChildren, "end)"),
		AppendIf(hasChildren, "end,"),
		children,
		AppendIf(hasChildren, ")"),
	)
}
