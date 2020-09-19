package api

import "github.com/uscott/go-api-deribit/inout"

// Auth requests authorization
func (c *Client) Auth(key string, secret string) (err error) {
	params := inout.ClientCredentials{
		GrantType:    "client_credentials",
		ClientID:     key,
		ClientSecret: secret,
	}
	var result inout.AuthOut
	err = c.Call("public/auth", params, &result)
	if err != nil {
		return
	}
	c.auth.token = result.AccessToken
	c.auth.refresh = result.RefreshToken
	return
}

// Logout logs out
func (c *Client) Logout() (err error) {
	var result = struct {
	}{}
	err = c.Call("public/auth", nil, &result)
	return
}
