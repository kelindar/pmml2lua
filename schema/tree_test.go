package schema

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	input := `<Node id="1" score="will play" recordCount="100" defaultChild="2">
	<True/>
	<ScoreDistribution value="will play" recordCount="60" confidence="0.6"/>
	<ScoreDistribution value="may play" recordCount="30" confidence="0.3"/>
	<ScoreDistribution value="no play" recordCount="10" confidence="0.1"/>
	<Node id="2" score="will play" recordCount="50" defaultChild="3">
	  <True/>
	  <Node id="3" score="will play" recordCount="40">
	  	<True/>
	  </Node>
	  <Node id="4" score="no play" recordCount="10">
	  	<True/>
	  </Node>
	</Node>
  </Node>`

	var out Node
	assert.NoError(t, xml.Unmarshal([]byte(input), &out))
	assert.EqualValues(t, Node{
		ID:           "1",
		Score:        "will play",
		RecordCount:  100,
		DefaultChild: "2",
		Predicate:    &Predicate{True: new(True)},
		Distributions: []ScoreDistribution{
			{Value: "will play", RecordCount: 60, Confidence: 0.6},
			{Value: "may play", RecordCount: 30, Confidence: 0.3},
			{Value: "no play", RecordCount: 10, Confidence: 0.1},
		},
		Nodes: []Node{
			{
				ID:           "2",
				Score:        "will play",
				RecordCount:  50,
				DefaultChild: "3",
				Predicate:    &Predicate{True: new(True)},
				Nodes: []Node{
					{
						ID:          "3",
						Score:       "will play",
						RecordCount: 40,
						Predicate:   &Predicate{True: new(True)},
					},
					{
						ID:          "4",
						Score:       "no play",
						RecordCount: 10,
						Predicate:   &Predicate{True: new(True)},
					},
				},
			},
		},
	}, out)
}
