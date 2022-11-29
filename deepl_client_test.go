package deeplgo_test

import (
	"errors"
	"net/http"
	"os"
	"testing"

	deeplgo "github.com/ThibaudDemay/deepl-go"
	"github.com/stretchr/testify/suite"
)

// Private from deeplgo, expose in blackbox testing for now
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

type TestSuite struct {
	suite.Suite
	deeplClient *deeplgo.Client
}

func (suite *TestSuite) SetupTest() {
	suite.deeplClient = deeplgo.NewClient(getEnv("DEEPL_TEST_API_KEY", "DEEPL_NO_ENV_API_KEY"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (suite *TestSuite) Test_Client_GetSetter() {
	apiKey := "API_4_TESTING"
	suite.deeplClient.SetApiKey(apiKey)
	suite.Equal(suite.deeplClient.GetApiKey(), apiKey)

	baseURL := "http://dumb.testing.tld/"
	suite.deeplClient.SetBaseUrl(baseURL)
	suite.Equal(suite.deeplClient.GetBaseUrl(), baseURL)
}

func (suite *TestSuite) Test_Client_ApiKeyError() {
	suite.deeplClient.SetApiKey("NO_API_KEY")
	res, err := suite.deeplClient.GetUsage()

	suite.Nil(res)
	suite.ErrorContains(err, errForbidden.Error())
}

// Test endpoint /usage
// API return data without error
func (suite *TestSuite) Test_Client_Usage() {
	res, err := suite.deeplClient.GetUsage()

	if err != nil {
		suite.Error(err)
	}

	suite.Nil(err)
	suite.NotEmpty(res)
	suite.IsType(deeplgo.Usage{}, *res)
}

// Test endpoint /languages with type to unknown
// API return data with default type set as 'source' without error
func (suite *TestSuite) Test_Languages_Unknown() {
	res, err := suite.deeplClient.GetLanguages("unknown")

	suite.Nil(err)
	suite.NotEmpty(res)
	suite.IsType(deeplgo.Languages{}, *res)
}

// Test endpoint /languages with type to source
// API return data with type source without error
func (suite *TestSuite) Test_Languages_Source() {
	res, err := suite.deeplClient.GetSourceLanguages()

	suite.Nil(err)
	suite.NotEmpty(res)
	suite.IsType(deeplgo.Languages{}, *res)
}

// Test endpoint /languages with type to target
// API return data with type target without error
func (suite *TestSuite) Test_Languages_Target() {
	res, err := suite.deeplClient.GetTargetLanguages()

	suite.Nil(err)
	suite.NotEmpty(res)
	suite.IsType(deeplgo.Languages{}, *res)
}

// Test endpoint /glossary-language-pairs
// API return data without error
func (suite *TestSuite) Test_GlossaryLanguagePairs() {
	res, err := suite.deeplClient.GetGlossaryLanguagePairs()

	if err != nil {
		suite.Error(err)
	}

	suite.Nil(err)
	suite.NotEmpty(res)
	suite.IsType(deeplgo.GlossaryLanguagePairs{}, *res)
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
