package pmml2lua

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/kelindar/lua"
)

// WriterFor creates a new writer for an input + schema combination
func writerFor(input string, schema interface{}) (*Writer, func() string) {
	if err := xml.Unmarshal([]byte(input), schema); err != nil {
		panic(err)
	}

	out := new(strings.Builder)
	writer := New(out)

	// Return the writer and a result function
	return writer, func() string {
		writer.Flush()
		return out.String()
	}
}

// MakeScript makes a script for testing
func makeScript(format string, args ...interface{}) *lua.Script {
	s, err := lua.FromString("test.lua", fmt.Sprintf(`
		function main(v)
			%v
		end
	`, fmt.Sprintf(format, args...)))
	if err != nil {
		panic(err)
	}

	return s
}
