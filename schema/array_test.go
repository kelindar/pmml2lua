package schema

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArray_Ints(t *testing.T) {
	input := `<Array n="3" type="int">1 22 3</Array>`

	var out Array
	assert.NoError(t, xml.Unmarshal([]byte(input), &out))
	assert.EqualValues(t, Array{
		Length: 3,
		Type:   "int",
		Values: "1 22 3",
	}, out)

	v, err := out.Ints()
	assert.NoError(t, err)
	assert.Equal(t, []int{
		1, 22, 3,
	}, v)

	j, err := out.CSV()
	assert.NoError(t, err)
	assert.Equal(t, `1,22,3`, string(j))
}

func TestArray_Floats(t *testing.T) {
	input := `<Array n="3" type="real">1.5 22.1 3.95</Array>`

	var out Array
	assert.NoError(t, xml.Unmarshal([]byte(input), &out))
	assert.EqualValues(t, Array{
		Length: 3,
		Type:   "real",
		Values: "1.5 22.1 3.95",
	}, out)

	v, err := out.Floats()
	assert.NoError(t, err)
	assert.Equal(t, []float64{
		1.5, 22.1, 3.95,
	}, v)

	j, err := out.CSV()
	assert.NoError(t, err)
	assert.Equal(t, `1.5,22.1,3.95`, string(j))
}

func TestArray_Strings(t *testing.T) {
	input := `<Array n="3" type="string">ab  "a b"   "with \"quotes\" "</Array>`

	var out Array
	assert.NoError(t, xml.Unmarshal([]byte(input), &out))
	assert.EqualValues(t, Array{
		Length: 3,
		Type:   "string",
		Values: `ab  "a b"   "with \"quotes\" "`,
	}, out)

	v, err := out.Strings()
	assert.NoError(t, err)
	assert.Equal(t, []string{
		`ab`, `a b`, `with "quotes" `,
	}, v)

	j, err := out.CSV()
	assert.NoError(t, err)
	assert.Equal(t, `"ab","a b","with \"quotes\" "`, string(j))
}
