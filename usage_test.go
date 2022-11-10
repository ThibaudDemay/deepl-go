package deeplgo

import (
	"fmt"
	"os"
	"testing"
)

func TestUsage(t *testing.T) {
	t.Parallel()

	c := NewClient(os.Getenv("DEEPL_TEST_API_KEY"))

	res, err := c.GetUsage()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
