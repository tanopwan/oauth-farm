package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const googleAPIURL = "https://www.googleapis.com"
const timeout = time.Duration(10 * time.Second)

// GoogleClient client
type GoogleClient struct {
	ClientID     string
	ClientSecret string
}

// NewGoogleClient return *GoogleClient
func NewGoogleClient(clientID, clientSecret string) *GoogleClient {
	return &GoogleClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// GoogleTokenV4 response
type GoogleTokenV4 struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

// TokenV4 dao
// resource /oauth2/v4/token
func (c *GoogleClient) TokenV4(code string) (*GoogleTokenV4, error) {
	resource := "/oauth2/v4/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("redirect_uri", "postmessage")
	data.Set("grant_type", "authorization_code")

	u, err := url.ParseRequestURI(googleAPIURL)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to parse request uri: %s", err.Error()))
	}
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{Timeout: timeout}
	r, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to create request body: %s", err.Error()))
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.Wrap(ErrExternalDAO, fmt.Sprintf("failed to read response: %s", err.Error()))
	}
	defer resp.Body.Close()

	respBody := struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		Scope            string `json:"scope"`
		TokenType        string `json:"token_type"`
		IDToken          string `json:"id_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to decode response body: %s", err.Error()))
	}

	if respBody.Error != "" {
		return nil, errors.Wrap(ErrExternalDAO, fmt.Sprintf("error from api: %s, %s", respBody.Error, respBody.ErrorDescription))
	}

	response := GoogleTokenV4{
		AccessToken:  respBody.AccessToken,
		ExpiresIn:    respBody.ExpiresIn,
		RefreshToken: respBody.RefreshToken,
		Scope:        respBody.Scope,
		TokenType:    respBody.TokenType,
		IDToken:      respBody.IDToken,
	}

	return &response, nil
}

// GoogleUserInfoV1 response
type GoogleUserInfoV1 struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

// UserInfoV1 dao
// resource /oauth2/v1/userinfo
func (c *GoogleClient) UserInfoV1(accessToken string) (*GoogleUserInfoV1, error) {
	resource := "/oauth2/v1/userinfo"

	u, err := url.ParseRequestURI(googleAPIURL)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to parse request uri: %s", err.Error()))
	}
	u.Path = resource
	query := url.Values{}
	query.Add("alt", "json")
	query.Add("access_token", accessToken)
	u.RawQuery = query.Encode()
	urlStr := u.String()

	fmt.Printf("urlStr: %s\n", urlStr)

	client := &http.Client{Timeout: timeout}
	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to create request body: %s", err.Error()))
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.Wrap(ErrExternalDAO, fmt.Sprintf("failed to read response: %s", err.Error()))
	}
	defer resp.Body.Close()

	userInfoResp := struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Link          string `json:"link"`
		Picture       string `json:"picture"`
		Gender        string `json:"gender"`
		Locale        string `json:"locale"`
		Error         struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"error"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&userInfoResp)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to decode response body: %s", err.Error()))
	}

	if userInfoResp.Error.Status != "" {
		return nil, errors.Wrap(ErrExternalDAO, fmt.Sprintf("error from api: %s, %s", userInfoResp.Error.Status, userInfoResp.Error.Message))
	}

	response := GoogleUserInfoV1{
		ID:            userInfoResp.ID,
		Email:         userInfoResp.Email,
		VerifiedEmail: userInfoResp.VerifiedEmail,
		Name:          userInfoResp.Name,
		GivenName:     userInfoResp.GivenName,
		FamilyName:    userInfoResp.FamilyName,
		Link:          userInfoResp.Link,
		Picture:       userInfoResp.Picture,
		Gender:        userInfoResp.Gender,
		Locale:        userInfoResp.Locale,
	}

	return &response, nil
}

// GoogleTokenInfoV3 response
type GoogleTokenInfoV3 struct {
	Iss           string    `json:"iss"`
	Azp           string    `json:"azp"`
	Aud           string    `json:"aud"`
	Sub           string    `json:"sub"`
	Email         string    `json:"email"`
	EmailVerified string    `json:"email_verified"`
	AtHash        string    `json:"at_hash"`
	Name          string    `json:"name"`
	Picture       string    `json:"picture"`
	GivenName     string    `json:"given_name"`
	FamilyName    string    `json:"family_name"`
	Locale        string    `json:"locale"`
	Iat           string    `json:"iat"`
	Exp           time.Time `json:"exp"`
	Jti           string    `json:"jti"`
	Alg           string    `json:"alg"`
	Kid           string    `json:"kid"`
	Typ           string    `json:"typ"`
}

// TokenInfoForTest dao
// resouce /oauth2/v3/tokeninfo
/*
For debugging purposes, you can use Google’s tokeninfo endpoint. Suppose your ID token’s value is XYZ123.
Then you would dereference the URI https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=XYZ123.
If the token is valid, the response would be its decoded JSON form.

This involves an HTTP round trip, introducing latency and the potential for network breakage.
The tokeninfo endpoint is useful for debugging but for production purposes, retrieve Google’s public keys
from the keys endpoint and perform the validation locally. You should retrieve the keys URI from
the Discovery document using the key jwks_uri.

Since Google changes its public keys only infrequently (on the order of once per day), you can cache them and,
in the vast majority of cases, perform local validation much more efficiently than by using
the tokeninfo endpoint. This requires retrieving and parsing certificates, and making
the appropriate crypto calls to check the signature. Fortunately, there are well-debugged libraries available
in a wide variety of languages to accomplish this (see jwt.io).
*/
func (c *GoogleClient) TokenInfoForTest(IDToken string) (*GoogleTokenInfoV3, error) {
	resource := "/oauth2/v3/tokeninfo"

	u, err := url.ParseRequestURI(googleAPIURL)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to parse request uri: %s", err.Error()))
	}
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{Timeout: timeout}
	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to create request body: %s", err.Error()))
	}
	q := r.URL.Query()
	q.Add("id_token", IDToken)
	r.URL.RawQuery = q.Encode()

	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.Wrap(ErrExternalDAO, fmt.Sprintf("failed to read response: %s", err.Error()))
	}
	defer resp.Body.Close()

	respBody := struct {
		Iss              string `json:"iss"`
		Azp              string `json:"azp"`
		Aud              string `json:"aud"`
		Sub              string `json:"sub"`
		Email            string `json:"email"`
		EmailVerified    string `json:"email_verified"`
		AtHash           string `json:"at_hash"`
		Name             string `json:"name"`
		Picture          string `json:"picture"`
		GivenName        string `json:"given_name"`
		FamilyName       string `json:"family_name"`
		Locale           string `json:"locale"`
		Iat              string `json:"iat"`
		Exp              string `json:"exp"`
		Jti              string `json:"jti"`
		Alg              string `json:"alg"`
		Kid              string `json:"kid"`
		Typ              string `json:"typ"`
		ErrorDescription string `json:"error_description"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to decode response body: %s", err.Error()))
	}

	if respBody.ErrorDescription != "" {
		return nil, errors.Wrap(ErrExternalDAO, fmt.Sprintf("error from api error_description: '%s'", respBody.ErrorDescription))
	}

	response := GoogleTokenInfoV3{
		Iss: respBody.Iss,
		Aud: respBody.Aud,
		// Exp: respBody.Exp,
	}

	return &response, nil
}
