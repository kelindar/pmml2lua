package pmml2lua

import (
	"encoding/xml"

	"github.com/kelindar/lua"
)

// scopeFor creates a new writer for an input + schema combination
func scopeFor(input string, schema interface{}) (*Scope, *Scope, func() string) {
	if err := xml.Unmarshal([]byte(input), schema); err != nil {
		panic(err)
	}

	// Create an outer scope with a global space and a main func
	outer := NewScope()
	global := outer.Scope()
	body := outer.Function("main", "v")

	return body, global, func() string {
		code, err := outer.Compile()
		if err != nil {
			panic(err)
		}

		return string(code)
	}
}

// MakeScript makes a script for testing
func makeScript(code string) *lua.Script {
	s, err := lua.FromString("test.lua", code)
	if err != nil {
		panic(err)
	}
	return s
}
