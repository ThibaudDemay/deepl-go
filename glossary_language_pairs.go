package deeplgo

type GlossaryLanguagePairs struct {
	SupportedLanguages []struct {
		SourceLang string `json:"source_lang" validate:"required"`
		TargetLang string `json:"target_lang" validate:"required"`
	} `json:"supported_languages" validate:"required"`
}

func (c *Client) GetGlossaryLanguagePairs() (*GlossaryLanguagePairs, error) {
	url := c.baseURL + glossaryLanguagePairsEndpoint

	res := GlossaryLanguagePairs{}
	if err := c.httpClient.Get(url, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
