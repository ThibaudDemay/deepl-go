package deeplgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	errBadRequest               = errors.New("Bad request. Please check error message and your parameters.")                                      // 400
	errAuthorizationFailed      = errors.New("Authorization failed. Please supply a valid `DeepL-Auth-Key` via the `Authorization` header.")      // 401
	errForbidden                = errors.New("Forbidden. The access to the requested resource is denied, because of insufficient access rights.") // 403
	errRequestResourceNotFound  = errors.New("The requested resource could not be found.")                                                        // 404
	errRequestSizeExceedsLimit  = errors.New("The request size exceeds the limit.")                                                               // 413
	errAcceptHeaderNotSupported = errors.New("The requested entries format specified in the `Accept` header is not supported.")                   // 415
	errTooManyRequests          = errors.New("Too many requests. Please wait and resend your request.")                                           // 429 && 529
	errQuotaExceeded            = errors.New("Quota exceeded. The character limit has been reached.")                                             // 456
	errInternalError            = errors.New("Internal error.")                                                                                   // 500
	errResourceUnavailable      = errors.New("Resource currently unavailable. Try again later.")                                                  // 503
)

var errCodes = map[int]error{
	http.StatusBadRequest:            errBadRequest,
	http.StatusUnauthorized:          errAuthorizationFailed,
	http.StatusForbidden:             errForbidden,
	http.StatusNotFound:              errRequestResourceNotFound,
	http.StatusRequestEntityTooLarge: errRequestSizeExceedsLimit,
	http.StatusUnsupportedMediaType:  errAcceptHeaderNotSupported,
	http.StatusTooManyRequests:       errTooManyRequests,
	456:                              errQuotaExceeded,
	http.StatusInternalServerError:   errInternalError,
	http.StatusServiceUnavailable:    errResourceUnavailable,
	529:                              errTooManyRequests,
}

var (
	baseProUrl                    = "https://api.deepl.com/v2"
	baseFreeUrl                   = "https://api-free.deepl.com/v2"
	usageEndpoint                 = "/usage"
	languagesEndpoint             = "/languages?target=%s"
	glossaryLanguagePairsEndpoint = "/glossary-language-pairs"
	translateEndpoint             = "/translate"
	glossariesEndpoint            = "/glossaries"
	glossaryEndpoint              = "/glossaries/%d"
	glossaryEntriesEndpoint       = "/glossaries/%d/entries"
	documentEndpoint              = "/document"
	documentStatusEndpoint        = "/document/%d"
	documentResultEndpoint        = "/document/%d/result"
)

type Client struct {
	baseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(APIKey string) *Client {
	baseURL := baseProUrl
	if strings.HasSuffix(APIKey, ":fx") {
		baseURL = baseFreeUrl
	}
	return &Client{
		baseURL: baseURL,
		APIKey:  APIKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) sendRequest(req *http.Request, dataInterface interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", c.APIKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if err, ok := errCodes[res.StatusCode]; ok {
			return err
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&dataInterface); err != nil {
		return err
	}

	return nil
}
