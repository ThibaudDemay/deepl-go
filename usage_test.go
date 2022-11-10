package deeplgo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test endpoint /usage with invalid API KEY
// API return Forbidden error and no data
func TestUsageErrInvalidAPIKey(t *testing.T) {
	t.Parallel()

	c := NewClient("NO_API_KEY")

	res, err := c.GetUsage()

	assert.Nil(t, res)
	assert.ErrorIs(t, err, errForbidden)
}

// Test endpoint /usage
// API return data without error
func TestUsage(t *testing.T) {
	t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetUsage()

	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.IsType(t, Usage{}, *res)
}
