package deeplgo

import (
	"errors"
	"net/http"
	"strings"
)

var (
	errBadRequest               = errors.New("bad request")                                                                                      // 400
	errAuthorizationFailed      = errors.New("authorization failed. please supply a valid `DeepL-Auth-Key` via the `Authorization` header")      // 401
	errForbidden                = errors.New("forbidden. the access to the requested resource is denied, because of insufficient access rights") // 403
	errRequestResourceNotFound  = errors.New("the requested resource could not be found")                                                        // 404
	errRequestSizeExceedsLimit  = errors.New("the request size exceeds the limit")                                                               // 413
	errAcceptHeaderNotSupported = errors.New("the requested entries format specified in the `Accept` header is not supported")                   // 415
	errTooManyRequests          = errors.New("too many requests. please wait and resend your request")                                           // 429 && 529
	errQuotaExceeded            = errors.New("quota exceeded. the character limit has been reached")                                             // 456
	errInternalError            = errors.New("internal error")                                                                                   // 500
	errResourceUnavailable      = errors.New("resource currently unavailable. try again later")                                                  // 503
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

type ErrorMessage struct {
	Message *string `json:"message,omitempty"`
	Detail  *string `json:"detail,omitempty"`
}

type Client struct {
	baseURL    string
	httpClient *HTTPClient
}

func NewClient(apiKey string) *Client {
	httpClient := NewHTTPClient(apiKey)

	baseURL := baseProUrl
	if strings.HasSuffix(apiKey, ":fx") {
		baseURL = baseFreeUrl
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}
