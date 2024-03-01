package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	// Run ReceiveChangeStream to test the function
	// no return value, just test if it runs without runtime errors
	ReceiveChangeStream("database", "collection")

	// Run ReceiveBytes to test the function
	// no return value, just test if it runs without runtime errors
	ReceiveBytes("database", "collection", 256)

	// Run HandleChangeEventFailed to test the function
	// no return value, just test if it runs without runtime errors
	HandleChangeEventFailed("database", "collection")

	// Run HandleChangeEventSuccess to test the function
	// no return value, just test if it runs without runtime errors
	HandleChangeEventSuccess("database", "collection")

	// Check if the Collectors function returns a non-empty slice
	cols := Collectors()
	assert.NotEqual(t, len(cols), 0)
}
