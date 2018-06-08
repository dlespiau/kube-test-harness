package harness

import (
	"testing"
)

func TestToSnake(t *testing.T) {
	cases := [][]string{
		[]string{"testCase", "test-case"},
		[]string{"TestCase", "test-case"},
		[]string{"Test Case", "test-case"},
		[]string{" Test Case", "test-case"},
		[]string{"Test Case ", "test-case"},
		[]string{" Test Case ", "test-case"},
		[]string{"test", "test"},
		[]string{"test-case", "test-case"},
		[]string{"Test", "test"},
		[]string{"", ""},
		[]string{"ManyManyWords", "many-many-words"},
		[]string{"manyManyWords", "many-many-words"},
		[]string{"AnyKind of-string", "any-kind-of-string"},
		[]string{"numbers2and55with000", "numbers-2-and-55-with-000"},
		[]string{"JSONData", "json-data"},
		[]string{"userID", "user-id"},
		[]string{"AAAbbb", "aa-abbb"},
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
