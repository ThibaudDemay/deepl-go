package deeplgo

import (
	"fmt"
	"os"
	"testing"
)

func TestSourceLanguages(t *testing.T) {
	t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetSourceLanguages()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}

func TestTargetLanguages(t *testing.T) {
	t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetSourceLanguages()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
