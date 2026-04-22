// Package oauth implements Google OAuth 2.0 login (web flow).
// Interactive user login only — Drive uploads use a separate offline refresh_token flow.
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	authEndpoint     = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint    = "https://oauth2.googleapis.com/token"
	userinfoEndpoint = "https://openidconnect.googleapis.com/v1/userinfo"
)

// Config holds OAuth client credentials + redirect URI.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// GoogleUserinfo is the subset of OIDC userinfo we consume.
type GoogleUserinfo struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// AuthorizeURL builds the Google consent page URL for the given CSRF state.
// access_type=online → no refresh token (we don't need one; JWT session is our own).
func (c *Config) AuthorizeURL(state string) string {
	q := url.Values{
		"client_id":     {c.ClientID},
		"redirect_uri":  {c.RedirectURI},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"access_type":   {"online"},
		"prompt":        {"select_account"},
	}
	return authEndpoint + "?" + q.Encode()
}

// ExchangeCode swaps an authorization code for an access token.
func (c *Config) ExchangeCode(ctx context.Context, code string) (string, error) {
	body := url.Values{
		"code":          {code},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"redirect_uri":  {c.RedirectURI},
		"grant_type":    {"authorization_code"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint,
		strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out struct {
		AccessToken string `json:"access_token"`
		ErrorDesc   string `json:"error_description"`
		ErrorKind   string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.AccessToken == "" {
		if out.ErrorDesc != "" {
			return "", fmt.Errorf("google oauth: %s", out.ErrorDesc)
		}
		return "", fmt.Errorf("google oauth: %s", out.ErrorKind)
	}
	return out.AccessToken, nil
}

// FetchUserinfo queries the OIDC userinfo endpoint with a valid access token.
func (c *Config) FetchUserinfo(ctx context.Context, accessToken string) (*GoogleUserinfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userinfoEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo status=%d", resp.StatusCode)
	}
	var u GoogleUserinfo
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
