package deeplgo

import (
	"fmt"
)

type Languages []struct {
	Language          string `json:"language" validate:"required"`
	Name              string `json:"name" validate:"required"`
	SupportsFormality bool   `json:"supports_formality"`
}

type LanguageType string

const (
	Source LanguageType = "source"
	Target LanguageType = "target"
)

func (c *Client) GetLanguages(target LanguageType, option ...string) (*Languages, error) {
	url := c.baseURL + fmt.Sprintf(languagesEndpoint, target)

	res := Languages{}
	if err := c.httpClient.Get(url, []QueryParameter{}, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetSourceLanguages(options ...string) (*Languages, error) {
	return c.GetLanguages(Source, options...)
}

func (c *Client) GetTargetLanguages(options ...string) (*Languages, error) {
	return c.GetLanguages(Target, options...)
}
