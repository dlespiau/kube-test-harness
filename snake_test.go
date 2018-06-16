package harness

import (
	"testing"
)

func TestToSnake(t *testing.T) {
	cases := [][]string{
		{"testCase", "test-case"},
		{"TestCase", "test-case"},
		{"Test Case", "test-case"},
		{" Test Case", "test-case"},
		{"Test Case ", "test-case"},
		{" Test Case ", "test-case"},
		{"test", "test"},
		{"test-case", "test-case"},
		{"Test", "test"},
		{"", ""},
		{"ManyManyWords", "many-many-words"},
		{"manyManyWords", "many-many-words"},
		{"AnyKind of-string", "any-kind-of-string"},
		{"numbers2and55with000", "numbers-2-and-55-with-000"},
		{"JSONData", "json-data"},
		{"userID", "user-id"},
		{"AAAbbb", "aa-abbb"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := toSnake(in)
		if result != out {
			t.Error("'" + result + "' != '" + out + "'")
		}
	}
}
