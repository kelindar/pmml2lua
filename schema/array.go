package schema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
)

// ArrayType ...
type ArrayType string

// Array ...
type Array struct {
	Length int       `xml:"n,attr,omitempty"`
	Type   ArrayType `xml:"type,attr"`
	Values string    `xml:",chardata"`
}

// Ints converts the values to the integer slice.
func (a Array) Ints() ([]int, error) {
	arr := strings.Split(string(a.Values), " ")
	out := make([]int, 0, len(arr))
	for _, s := range arr {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, int(v))
	}
	return out, nil
}

// Floats converts the values to the float64 slice.
func (a Array) Floats() ([]float64, error) {
	arr := strings.Split(string(a.Values), " ")
	out := make([]float64, 0, len(arr))
	for _, s := range arr {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

// Strings converts the values to the string slice.
func (a Array) Strings() ([]string, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(a.Values))

	quoted := regexp.MustCompile(`^"(.*)"$`)
	v := make([]string, 0, a.Length)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		t := s.TokenText()
		b := []byte(t)
		if quoted.Match(b) {
			if err := json.Unmarshal([]byte(t), &t); err != nil {
				return nil, err
			}
		}

		v = append(v, t)
	}

	return v, nil
}

// CSV returns the comma separated values of the array
func (a Array) CSV() (string, error) {
	var arr interface{}
	var err error

	switch a.Type {
	case "int":
		arr, err = a.Ints()
	case "real":
		arr, err = a.Floats()
	case "string":
		arr, err = a.Strings()
	default:
		err = fmt.Errorf("unsupported array type %v", a.Type)
	}

	if err != nil {
		return "", err
	}

	b, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}

	return regexp.MustCompile(`^\[(.*)\]$`).ReplaceAllString(string(b), `$1`), nil
}
