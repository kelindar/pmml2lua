package pmml2lua

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kelindar/lua"
)

// scopeFor creates a new writer for an input + schema combination
func scopeFor(input string, schema interface{}) (*Scope, *Scope, func() string) {
	content := []byte(input)
	if strings.HasSuffix(input, ".xml") {
		b, err := ioutil.ReadFile(input)
		if err != nil {
			panic(err)
		}
		content = b
	}

	// Unmarshal the XML
	if err := xml.Unmarshal(content, schema); err != nil {
		panic(err)
	}

	// Create an outer scope with a global space and a main func
	main := NewScope()
	body := main.Function("main", "v")
	global := NewScope().With(
		Append(`local tree = require("tree")`),
	)

	return body, global, func() string {
		code1, err := global.Compile()
		if err != nil {
			panic(err)
		}

		code2, err := main.Compile()
		if err != nil {
			panic(err)
		}

		return string(code1) + "\n\n" + string(code2)
	}
}

// MakeScript makes a script for testing
func makeScript(code string) *lua.Script {
	f, _ := os.Open("tree.lua")
	moduleCode, err := lua.FromReader("tree.lua", f)
	module := &lua.ScriptModule{
		Script:  moduleCode,
		Name:    "tree",
		Version: "1.0.0",
	}
	if err != nil {
		panic(err)
	}

	s, err := lua.FromString("test.lua", code, module)
	if err != nil {
		panic(err)
	}
	return s
}
