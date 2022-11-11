package deeplgo

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test endpoint /glossary-language-pairs with invalid API KEY
// API return Forbidden error and no data
func TestGlossaryLanguagePairsErrInvalidAPIKey(t *testing.T) {
	// t.Parallel()

	c := NewClient("NO_API_KEY")

	res, err := c.GetGlossaryLanguagePairs()

	assert.Nil(t, res)
	assert.ErrorContains(t, err, errAuthorizationFailed.Error())
}

// Test endpoint /glossary-language-pairs
// API return data without error
func TestGlossaryLanguagePairs(t *testing.T) {
	// t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetGlossaryLanguagePairs()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)

	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.IsType(t, GlossaryLanguagePairs{}, *res)
}
