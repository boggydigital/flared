package cf_api

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

func (c *Client) newAuthRequest(
	method string,
	u *url.URL,
	bts ...byte) (*http.Request, error) {

	var reader io.Reader
	if len(bts) > 0 {
		reader = bytes.NewReader(bts)
	}

	req, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add(authorizationHeader, bearerPrefix+c.token)

	return req, err
}
