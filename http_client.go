package deeplgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

type HTTPClient struct {
	client   *http.Client
	validate *validator.Validate
	apiKey   string
}

type QueryParameter struct {
	key   string
	value string
}

func NewHTTPClient(apiKey string) *HTTPClient {
	return &HTTPClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: time.Minute,
		},
		validate: validator.New(),
	}
}

// ProcessError check if HTTP status code is in list of status code known then
// if request return data try to parse it. And if HTTP status is unknown return
// status code not managed.
func (hc *HTTPClient) ProcessError(resp *http.Response) error {
	if err, ok := errCodes[resp.StatusCode]; ok {
		var errResp ErrorMessage
		if errDec := json.NewDecoder(resp.Body).Decode(&errResp); errDec == nil {
			errText := err.Error()
			if errResp.Message != nil {
				errText += ", message : " + *errResp.Message
			}
			if errResp.Detail != nil {
				errText += ", detail : " + *errResp.Detail
			}
			return errors.New(errText)
		}
		return err
	}

	return fmt.Errorf("unknown error, status code: %d", resp.StatusCode)
}

// SendRequest Take http.Request to send it and manage errors from several
// source and set API Key in header `Authorization` for current request.
func (hc *HTTPClient) SendRequest(req *http.Request, dataInterface interface{}) error {
	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", hc.apiKey))

	var err error
	resp, _ := hc.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return hc.ProcessError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(dataInterface); err != nil {
		return err
	}

	// In case of data is JSON Array of JSON Object directly  they each have
	// their own way of being verified
	if reflect.TypeOf(dataInterface).Elem().Kind() == reflect.Slice {
		err = hc.validate.Var(dataInterface, "dive")
	} else if reflect.TypeOf(dataInterface).Elem().Kind() == reflect.Struct {
		err = hc.validate.Struct(dataInterface)
	}

	if err != nil {
		return err
	}

	return nil
}

// Get Wrap request creation and SendRequest call in HTTP GET context
func (hc *HTTPClient) Get(url string, queryParameters []QueryParameter, dataInterface interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if len(queryParameters) > 0 {
		currentQ := req.URL.Query()
		for _, q := range queryParameters {

			currentQ.Add(q.key, q.value)
		}
		req.URL.RawQuery = currentQ.Encode()
	}

	if err != nil {
		return err
	} else if err := hc.SendRequest(req, dataInterface); err != nil {
		return err
	}

	return nil
}

// Post Wrap request creation and SendRequest call in HTTP POST context
func (hc *HTTPClient) Post(url string, data io.Reader, dataInterface interface{}) error {
	req, err := http.NewRequest(http.MethodPost, url, data)

	if err != nil {
		return err
	} else if err := hc.SendRequest(req, dataInterface); err != nil {
		return err
	}

	return nil
}
