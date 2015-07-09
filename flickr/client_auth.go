package flickr

import (
	"encoding/json"
	"errors"
)

type FrobResp struct {
	Frob Content
}

func (c *Client) GetAuthURL() (frob string, url string, err error) {
	b, err := c.Do("flickr.auth.getFrob", nil)
	if err != nil {
		return
	}
	resp := FrobResp{}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return
	}
	url = URL(c.apiKey, c.secret, authEndpoint, Args{"perms": "delete", "frob": string(resp.Frob)})
	return
}

type TokenResp struct {
	Auth struct {
		Token, Perms Content
		User         User
	}
}

func (c *Client) getUser(newToken bool, val string) (*User, error) {
	if val == "" {
		return nil, errors.New("Missing value")
	}
	var b []byte
	var err error
	if newToken {
		b, err = c.Do("flickr.auth.getToken", Args{"frob": val})
	} else {
		b, err = c.Do("flickr.auth.checkToken", Args{"auth_token": val})
	}
	if err != nil {
		return nil, err
	}
	resp := TokenResp{}
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	user := &resp.Auth.User
	user.c = c
	user.Token = string(resp.Auth.Token)
	user.Perms = string(resp.Auth.Perms)
	return user, nil
}

func (c *Client) UserFromFrob(frob string) (*User, error) {
	return c.getUser(true, frob)
}

func (c *Client) UserFromToken(token string) (*User, error) {
	return c.getUser(false, token)
}
