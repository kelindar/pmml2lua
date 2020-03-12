package schema

import (
	"encoding/xml"
	"io"
	"strconv"
)

// DecisionTree ...
type DecisionTree struct {
	ModelName            string      `xml:"modelName,attr,omitempty"`
	FunctionName         string      `xml:"functionName,attr"`
	AlgorithmName        string      `xml:"algorithmName,attr,omitempty"`
	MissingValueStrategy string      `xml:"missingValueStrategy,attr,omitempty"`
	MissingValuePenalty  float64     `xml:"missingValuePenalty,attr,omitempty"`
	NoTrueChildStrategy  string      `xml:"noTrueChildStrategy,attr,omitempty"`
	Extension            []Extension `xml:"Extension"`
	Node                 Node        `xml:"Node"`

	//SplitCharacteristicAttr  interface{}           `xml:"splitCharacteristic,attr,omitempty"`
	//Output                   *Output               `xml:"Output"`
	//ModelStats               *ModelStats           `xml:"ModelStats"`
	//Targets                  *Targets              `xml:"Targets"`
	//LocalTransformations     *LocalTransformations `xml:"LocalTransformations"`
	//ResultField              []*ResultField        `xml:"ResultField"`
}

// Node ...
type Node struct {
	ID            string              `xml:"id,attr,omitempty"`
	Score         string              `xml:"score,attr,omitempty"`
	RecordCount   int64               `xml:"recordCount,attr,omitempty"`
	DefaultChild  string              `xml:"defaultChild,attr,omitempty"`
	Extension     []Extension         `xml:"Extension"`
	Distributions []ScoreDistribution `xml:"ScoreDistribution"`
	Nodes         []Node              `xml:"Node"`
	Predicate     *Predicate
	//EmbeddedModel     *EmbeddedModel
	//Partition         *Partition           `xml:"Partition"`
}

// UnmarshalXML ...
func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	for _, v := range start.Attr {
		switch v.Name.Local {
		case "id":
			n.ID = v.Value
		case "score":
			n.Score = v.Value
		case "defaultChild":
			n.DefaultChild = v.Value
		case "recordCount":
			if n.RecordCount, err = strconv.ParseInt(v.Value, 10, 64); err != nil {
				return err
			}
		}
	}

	var t xml.Token
	for t, err = d.Token(); t != nil && err != io.EOF; t, err = d.Token() {
		if err != nil {
			return err
		}

		if el, ok := t.(xml.StartElement); ok {
			switch el.Name.Local {
			case "Node":
				var child Node
				if err := d.DecodeElement(&child, &el); err != nil {
					return err
				}
				n.Nodes = append(n.Nodes, child)

			case "ScoreDistribution":
				var dist ScoreDistribution
				if err := d.DecodeElement(&dist, &el); err != nil {
					return err
				}
				n.Distributions = append(n.Distributions, dist)

			default:
				n.Predicate = new(Predicate)
				if err := d.DecodeElement(n.Predicate, &el); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// ScoreDistribution ...
type ScoreDistribution struct {
	Extension   []Extension `xml:"Extension"`
	Value       string      `xml:"value,attr"`
	RecordCount int64       `xml:"recordCount,attr"`
	Confidence  float64     `xml:"confidence,attr,omitempty"`
	Probability float64     `xml:"probability,attr,omitempty"`
}
