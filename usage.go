package deeplgo

import "net/http"

type Usage struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}

func (c *Client) GetUsage() (*Usage, error) {
	url := c.baseURL + usageEndpoint
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res := Usage{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
