package deeplgo

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test HTTPClient function ProcessError with request return 404
// Function must return an error RequestResourceNotFound
func Test_HTTPClient_ProcessErrorBasic(t *testing.T) {
	body := io.NopCloser(bytes.NewReader([]byte(``)))

	hc := NewHTTPClient("NO_API_KEY")

	err := hc.ProcessError(&http.Response{
		StatusCode: 404,
		Body:       body,
	})

	assert.ErrorIs(t, err, errRequestResourceNotFound)
}

// Test HTTPClient function ProcessError with data in response (var json)
// include message only.
// Function must return an error RequestResourceNotFound with message
func Test_HTTPClient_ProcessErrorWithMessage(t *testing.T) {
	json := `{"message":"Message test"}`
	body := io.NopCloser(bytes.NewReader([]byte(json)))

	hc := NewHTTPClient("NO_API_KEY")

	err := hc.ProcessError(&http.Response{
		StatusCode: 404,
		Body:       body,
	})

	errTarget := errRequestResourceNotFound.Error() + ", message : Message test"
	assert.EqualError(t, err, errTarget)
}

// Test HTTPClient function ProcessError with data in response (var json)
// include message and detail.
// Function must return an error ResourceNotFound with message and detail
func Test_HTTPClient_ProcessErrorWithMessageAndDetail(t *testing.T) {
	json := `{"message":"Message test","detail":"Detail test"}`
	body := io.NopCloser(bytes.NewReader([]byte(json)))

	hc := NewHTTPClient("NO_API_KEY")

	err := hc.ProcessError(&http.Response{
		StatusCode: 404,
		Body:       body,
	})

	errTarget := errRequestResourceNotFound.Error() + ", message : Message test, detail : Detail test"
	assert.EqualError(t, err, errTarget)
}

// Test HTTPClient function ProcessError were http return code not in err list.
// Function must return an error unknow error with status code
func Test_HTTPClient_ProcessErrorWithUnknownError(t *testing.T) {
	body := io.NopCloser(bytes.NewReader([]byte(``)))

	hc := NewHTTPClient("NO_API_KEY")

	err := hc.ProcessError(&http.Response{
		StatusCode: 418,
		Body:       body,
	})

	assert.EqualError(t, err, "unknown error, status code: 418")
}

type TestStruct struct {
	Page  int      `json:"page" validate:"required"`
	Count int      `json:"count" validate:"required"`
	Data  []string `json:"data" validate:"required,dive,min=0"`
}

type TestArray []struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// Test HTTPClient function SendRequest verify headers Authorization
func Test_HTTPClient_SendRequestHeaderAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Header.Get("Authorization"), "DeepL-Auth-Key AN_API_KEY")
			w.Write([]byte(`{"page":1,"count":2,"data":["Apple","Pear"]}`))
		},
	))

	defer server.Close()

	hc := NewHTTPClient("AN_API_KEY")

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.Nil(t, err)

	res := TestStruct{}
	err = hc.SendRequest(req, &res)
	assert.Nil(t, err)
}

// Test HTTPClient function SendRequest with Json object.
// Function must return no error and decode response has correct struct
func Test_HTTPClient_SendRequestWithJsonObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":2,"data":["Apple","Pear"]}`))
		},
	))

	defer server.Close()

	hc := NewHTTPClient("NO_API_KEY")

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.Nil(t, err)

	res := TestStruct{}
	err = hc.SendRequest(req, &res)
	assert.Nil(t, err)
	assert.Equal(t, TestStruct{
		Page:  1,
		Count: 2,
		Data:  []string{"Apple", "Pear"},
	}, res)
}

// Test HTTPClient function SendRequest with Json Array.
// Function must return no error and decode response has correct struct
func Test_HTTPClient_SendRequestWithJsonArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[
				{"title": "Apple", "description": "Forbidden fruit"},
				{"title": "Pear", "description": "Not meat part"}
			]`))
		},
	))

	defer server.Close()

	hc := NewHTTPClient("NO_API_KEY")

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.Nil(t, err)

	res := TestArray{}
	err = hc.SendRequest(req, &res)
	assert.Nil(t, err)
	assert.Equal(t, TestArray{
		{Title: "Apple", Description: "Forbidden fruit"},
		{Title: "Pear", Description: "Not meat part"},
	}, res)
}

// Test HTTPClient function SendRequest with malformatted Json object.
// Function must return an error with unexpected EOF
func Test_HTTPClient_SendRequestErrDecode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":2,"data":["Apple"]`))
		},
	))

	defer server.Close()

	hc := NewHTTPClient("NO_API_KEY")

	req, err := http.NewRequest("GET", server.URL, nil)

	assert.Nil(t, err)

	res := TestStruct{}
	err = hc.SendRequest(req, &res)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// Test HTTPClient function SendRequest with Json object not corresponding
// with struct declaration.
// Function must return an error about field validation data required
func Test_HTTPClient_SendRequestErrValidator(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":2,"donnees":"Apple"}`))
		},
	))

	defer server.Close()

	hc := NewHTTPClient("NO_API_KEY")

	req, err := http.NewRequest("GET", server.URL, nil)

	assert.Nil(t, err)

	res := TestStruct{}
	err = hc.SendRequest(req, &res)

	assert.NotNil(t, err)

	errTarget := errors.New("Key: 'TestStruct.Data' Error:Field validation for 'Data' failed on the 'required' tag")
	assert.Error(t, err, errTarget)
}

// Test HTTPClient function GET with Json object
// Function must return no error and decode response has correct struct
func Test_HTTPClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":1,"data":["Apple"]}`))
		},
	))

	defer server.Close()
	hc := NewHTTPClient("NO_API_KEY")

	res := TestStruct{}
	err := hc.Get(server.URL, []QueryParameter{}, &res)

	assert.Nil(t, err)
	assert.Equal(t, TestStruct{
		Page:  1,
		Count: 1,
		Data:  []string{"Apple"},
	}, res)
}

func Test_HTTPClient_GetQueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.URL.RawQuery, "keyTest=valueTest")
			w.Write([]byte(`{"page":1,"count":1,"data":["Apple"]}`))
		},
	))

	defer server.Close()
	hc := NewHTTPClient("NO_API_KEY")

	res := TestStruct{}
	qp := []QueryParameter{
		{key: "keyTest", value: "valueTest"},
	}
	err := hc.Get(server.URL, qp, &res)

	assert.Nil(t, err)
}

// Test HTTPClient function GET on bad url
// Function must return error
func Test_HTTPClient_GetErrorRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":1,"data":["Apple"]}`))
		},
	))

	defer server.Close()
	hc := NewHTTPClient("NO_API_KEY")

	res := TestStruct{}
	err := hc.Get(server.URL+"yolo", []QueryParameter{}, &res)

	assert.NotNil(t, err)
}

// Test HTTPClient function GET with bad json return
// Function must return error return by SendRequest
func Test_HTTPClient_GetErrorNested(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":2,"data":["Apple"]`))
		},
	))

	defer server.Close()

	hc := NewHTTPClient("NO_API_KEY")

	res := TestStruct{}
	err := hc.Get(server.URL, []QueryParameter{}, &res)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// Test HTTPClient function POST with Json object
// Function must return no error and decode response has correct struct
func Test_HTTPClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":1,"data":["Apple"]}`))
		},
	))

	defer server.Close()
	hc := NewHTTPClient("NO_API_KEY")

	res := TestStruct{}
	err := hc.Post(server.URL, nil, &res)

	assert.Nil(t, err)
	assert.Equal(t, TestStruct{
		Page:  1,
		Count: 1,
		Data:  []string{"Apple"},
	}, res)
}

// Test HTTPClient function POST on bad url
// Function must return error
func Test_HTTPClient_PostErrorRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":1,"data":["Apple"]}`))
		},
	))

	defer server.Close()
	hc := NewHTTPClient("NO_API_KEY")

	res := TestStruct{}
	err := hc.Post(server.URL+"yolo", nil, &res)

	assert.NotNil(t, err)
}

// Test HTTPClient function POST with bad json return
// Function must return error return by SendRequest
func Test_HTTPClient_PostErrorNested(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"page":1,"count":2,"data":["Apple"]`))
		},
	))

	defer server.Close()

	hc := NewHTTPClient("NO_API_KEY")

	res := TestStruct{}
	err := hc.Post(server.URL, nil, &res)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
