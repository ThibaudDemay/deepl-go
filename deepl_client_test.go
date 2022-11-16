package deeplgo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyJsonReturn struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestSendRequest(t *testing.T) {
	c := NewClient("NO_API_KEY")

	url := "https://dummyjson.com/http/404/spider_bob"
	req, err := http.NewRequest("GET", url, nil)

	assert.Nil(t, err)

	res := DummyJsonReturn{}
	err = c.httpClient.SendRequest(req, &res)

	assert.ErrorContains(t, err, errRequestResourceNotFound.Error())
}
