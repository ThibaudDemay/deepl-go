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

func (c *Client) GetLanguages(target LanguageType) (*Languages, error) {
	url := c.baseURL + fmt.Sprintf(languagesEndpoint, target)

	res := Languages{}
	if err := c.httpClient.Get(url, &res); err != nil {
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
