package deeplgo

type Usage struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}

func (c *Client) GetUsage() (*Usage, error) {
	url := c.baseURL + usageEndpoint

	res := Usage{}
	if err := c.httpClient.Get(url, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
