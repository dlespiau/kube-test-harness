package harness

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRemoveNamespace(t *testing.T) {
	test := &Test{}
	assert.Equal(t, len(test.namespaces), 0)

	test.removeNamespace("foobar")
	assert.Equal(t, len(test.namespaces), 0)

	test.addNamespace("ns1")
	assert.Equal(t, len(test.namespaces), 1)
	assert.Equal(t, test.namespaces[0], "ns1")

	test.addNamespace("ns2")
	assert.Equal(t, len(test.namespaces), 2)

	test.removeNamespace("ns1")
	assert.Equal(t, len(test.namespaces), 1)
	assert.Equal(t, test.namespaces[0], "ns2")

	test.removeNamespace("ns1")
	assert.Equal(t, len(test.namespaces), 1)
	assert.Equal(t, test.namespaces[0], "ns2")

	test.removeNamespace("ns2")
	assert.Equal(t, len(test.namespaces), 0)
}
