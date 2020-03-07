package schema

import (
	"encoding/xml"
	"fmt"
)

// Predicate ...
type Predicate struct {
	SimplePredicate    *SimplePredicate
	CompoundPredicate  *CompoundPredicate
	SimpleSetPredicate *SimpleSetPredicate
	True               *True
	False              *False
}

// UnmarshalXML ...
func (p *Predicate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "SimplePredicate":
		p.SimplePredicate = new(SimplePredicate)
		return d.DecodeElement(p.SimplePredicate, &start)
	case "CompoundPredicate":
		p.CompoundPredicate = new(CompoundPredicate)
		return d.DecodeElement(p.CompoundPredicate, &start)
	case "SimpleSetPredicate":
		p.SimpleSetPredicate = new(SimpleSetPredicate)
		return d.DecodeElement(p.SimpleSetPredicate, &start)
	case "True":
		p.True = new(True)
		return d.DecodeElement(p.True, &start)
	case "False":
		p.False = new(False)
		return d.DecodeElement(p.False, &start)
	default:
		return fmt.Errorf("unsupported predicate type")
	}
}

// CompoundPredicate ...
type CompoundPredicate struct {
	Operator   string
	Predicates []Predicate
}

// UnmarshalXML ...
func (p *CompoundPredicate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "booleanOperator" {
			p.Operator = attr.Value
		}
	}

	var done bool
	for !done {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch t := t.(type) {
		case xml.StartElement:
			var predicate Predicate
			if err := predicate.UnmarshalXML(d, t); err != nil {
				return err
			}
			p.Predicates = append(p.Predicates, predicate)
		case xml.EndElement:
			done = true
		}
	}

	return nil
}

// SimpleSetPredicate ...
type SimpleSetPredicate struct {
	Field     string      `xml:"field,attr"`
	Operator  string      `xml:"booleanOperator,attr"`
	Extension []Extension `xml:"Extension"`
	Array     *Array      `xml:"Array"`
}

// True ...
type True struct {
	Extension []Extension `xml:"Extension"`
}

// False ...
type False struct {
	Extension []Extension `xml:"Extension"`
}

// Extension ...
type Extension struct {
	Extender string `xml:"extender,attr,omitempty"`
	Name     string `xml:"name,attr,omitempty"`
	Value    Value  `xml:"value,attr,omitempty"`
}

// SimplePredicate ...
type SimplePredicate struct {
	Field     string      `xml:"field,attr"`
	Operator  string      `xml:"operator,attr"`
	Value     Value       `xml:"value,attr,omitempty"`
	Extension []Extension `xml:"Extension"`
}

// Array ...
type Array struct {
	Length int       `xml:"n,attr,omitempty"`
	Type   ArrayType `xml:"type,attr"`
}

// Value ...
type Value string

// ArrayType ...
type ArrayType string
