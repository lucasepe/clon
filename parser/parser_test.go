package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func oneFieldObjGen(field string, value interface{}) Generator {
	return NewObj().Add(field, Value{value: value})
}

func TestParseFieldGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input:    `a=^API_VERSION`,
			expected: oneFieldObjGen("a", "v1"),
		},
		{
			input:    `a=b`,
			expected: oneFieldObjGen("a", "b"),
		},
		{
			input:    `a=42`,
			expected: oneFieldObjGen("a", "42"),
		},
		{
			input:    `a=true`,
			expected: oneFieldObjGen("a", "true"),
		},
		{
			input:    `a=false`,
			expected: oneFieldObjGen("a", "false"),
		},
		{
			input:    `a=null`,
			expected: oneFieldObjGen("a", "null"),
		},

		{
			input:    `a=:42`,
			expected: oneFieldObjGen("a", int64(42)),
		},
		{
			input:    `a=:true`,
			expected: oneFieldObjGen("a", true),
		},
		{
			input:    `a=:false`,
			expected: oneFieldObjGen("a", false),
		},
		{
			input:    `a=:null`,
			expected: oneFieldObjGen("a", nil),
		},
	}

	os.Setenv("API_VERSION", "v1")

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := ParseString(cas.input)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestParseObjectGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input: `a={b=c}`,
			expected: Obj{
				fields: map[string]Generator{
					"a": Obj{
						fields: map[string]Generator{
							"b": Value{value: "c"},
						},
					},
				},
			},
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := ParseString(cas.input)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestParseDotObjectGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input: `a."b.b".c=d`,
			expected: Obj{
				fields: map[string]Generator{
					"a": Obj{
						fields: map[string]Generator{
							"b.b": Obj{
								fields: map[string]Generator{
									"c": Value{value: "d"},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `parent.child1=value1 parent.child2=value2`,
			expected: Obj{
				fields: map[string]Generator{
					"parent": Obj{
						fields: map[string]Generator{
							"child1": Value{value: "value1"},
							"child2": Value{value: "value2"},
						},
					},
				},
			},
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := ParseString(cas.input)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestComplexParse(t *testing.T) {
	expected := Obj{
		fields: map[string]Generator{
			"id":      Value{value: "42"},
			"enabled": Value{value: true},
			"caller": Obj{
				fields: map[string]Generator{
					"gender": Obj{
						fields: map[string]Generator{
							"code": Value{value: int64(1)},
						},
					},
				},
			},
			"customer": Obj{
				fields: map[string]Generator{
					"name": Value{value: "Geralt"},
					"age":  Value{value: "86"},
					"address": Obj{
						fields: map[string]Generator{
							"zip": Value{value: "75018"},
						},
					},
				},
			},
		},
	}

	ast, err := ParseString(`id = 42 caller.gender.code = :1  customer={name = "Geralt" age  = 86 address.zip = 75018 } enabled = :true`)

	require.NoError(t, err)
	require.Equal(t, []Generator{expected}, ast)
}
