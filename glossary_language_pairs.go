package deeplgo

import "net/http"

type GlossaryLanguagePairs struct {
	SupportedLanguages []struct {
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
	} `json:"supported_languages"`
}

func (c *Client) GetGlossaryLanguagePairs() (*GlossaryLanguagePairs, error) {
	url := c.baseURL + glossaryLanguagePairsEndpoint
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res := GlossaryLanguagePairs{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
