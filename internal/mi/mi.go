package mi

import (
	"errors"
	"net/http"
)

type Client struct {
	cli     *http.Client
	baseURL string
	headers http.Header
}

func New(token string, userAgent string) *Client {
	headers := http.Header{}
	headers.Set("Authorization", token)
	headers.Set("User-Agent", userAgent)

	cli := &http.Client{}

	return &Client{
		cli:     cli,
		baseURL: "https://mi.within.website",
		headers: headers,
	}
}

func (c *Client) Refresh() error {
	if c == nil {
		return nil
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/blog/refresh", nil)
	if err != nil {
		return err
	}

	for k, v := range c.headers {
		req.Header[k] = v
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("non-200 status code")
	}

	return nil
}
