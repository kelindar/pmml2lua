package pmml2lua

import (
	"encoding/xml"
	"os"

	"github.com/kelindar/lua"
)

// scopeFor creates a new writer for an input + schema combination
func scopeFor(input string, schema interface{}) (*Scope, *Scope, func() string) {
	if err := xml.Unmarshal([]byte(input), schema); err != nil {
		panic(err)
	}

	// Create an outer scope with a global space and a main func
	outer := NewScope()
	global := outer.Scope().With(
		NewStatement().Append(`local eval = require("eval")`),
	)
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
	f, _ := os.Open("eval.lua")
	moduleCode, _ := lua.FromReader("eval.lua", f)
	module := &lua.ScriptModule{
		Script:  moduleCode,
		Name:    "eval",
		Version: "1.0.0",
	}

	s, err := lua.FromString("test.lua", code, module)
	if err != nil {
		panic(err)
	}
	return s
}
