package deeplgo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test endpoint /languages with invalid API KEY
// API return Forbidden error and no data
func TestLanguagesInvalidAPIKey(t *testing.T) {
	// t.Parallel()

	c := NewClient("INVALID_API_KEY")

	res, err := c.GetSourceLanguages()

	assert.Nil(t, res)
	assert.ErrorContains(t, err, errForbidden.Error())
}

// Test endpoint /languages with type to unknown
// API return data with default type set as 'source' without error
func TestLanguagesUnknown(t *testing.T) {
	// t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetLanguages("unknown")

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.IsType(t, Languages{}, *res)
}

// Test endpoint /languages with type to source
// API return data with type source without error
func TestSourceLanguages(t *testing.T) {
	// t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetSourceLanguages()

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.IsType(t, Languages{}, *res)
}

// Test endpoint /languages with type to target
// API return data with type target without error
func TestTargetLanguages(t *testing.T) {
	// t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetTargetLanguages()

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.IsType(t, Languages{}, *res)
}
