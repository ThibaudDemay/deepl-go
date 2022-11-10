package deeplgo

import (
	"fmt"
	"net/http"
)

type Languages []struct {
	Language          string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality bool   `json:"supports_formality"`
}

type LanguageType string

const (
	Source LanguageType = "source"
	Target LanguageType = "target"
)

func (c *Client) GetLanguages(target LanguageType) (*Languages, error) {
	url := c.baseURL + fmt.Sprintf(languagesEndpoint, target)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res := Languages{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetSourceLanguages() (*Languages, error) {
	return c.GetLanguages(Source)
}

func (c *Client) GetTargetLanguages() (*Languages, error) {
	return c.GetLanguages(Target)
}
