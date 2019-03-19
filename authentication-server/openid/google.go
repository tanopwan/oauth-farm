package openid

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
	"github.com/tanopwan/oauth-farm/common"

	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"time"
)

var cc = common.NewMemoryCache()

const timeout = time.Duration(10 * time.Second)

func getJWTKeys(token *jwt.Token) (interface{}, error) {
	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	if cc != nil {
		value := cc.get(keyID)
		if value != nil {
			fmt.Printf("cache found google_public\n")
			return value, nil
		}
	}

	fmt.Printf("cache not found google_public\n")

	resource := "/.well-known/openid-configuration"
	client := &http.Client{Timeout: timeout}

	var jwksURI string
	{
		u, err := url.ParseRequestURI("https://accounts.google.com")
		if err != nil {
			return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to parse request uri: %s", err.Error()))
		}
		u.Path = resource
		urlStr := u.String()

		r, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to create request body: %s", err.Error()))
		}

		resp, err := client.Do(r)
		if err != nil {
			return nil, errors.Wrap(ErrExternalDAO, fmt.Sprintf("failed to read response: %s", err.Error()))
		}
		defer resp.Body.Close()

		respBody := struct {
			Issuer                            string   `json:"issuer"`
			AuthorizationEndpoint             string   `json:"authorization_endpoint"`
			TokenEndpoint                     string   `json:"token_endpoint"`
			UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
			RevocationEndpoint                string   `json:"revocation_endpoint"`
			JwksURI                           string   `json:"jwks_uri"`
			ResponseTypesSupported            []string `json:"response_types_supported"`
			SubjectTypesSupported             []string `json:"subject_types_supported"`
			IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
			ScopesSupported                   []string `json:"scopes_supported"`
			TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
			ClaimsSupported                   []string `json:"claims_supported"`
			CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
		}{}

		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return nil, errors.Wrap(ErrInternalServerError, fmt.Sprintf("failed to decode response body: %s", err.Error()))
		}

		jwksURI = respBody.JwksURI
	}

	// we want to verify a JWT
	set, err := jwk.FetchHTTP(jwksURI)
	if err != nil {
		return nil, err
	}

	if key := set.LookupKeyID(keyID); len(key) == 1 {
		m, err := key[0].Materialize()
		if err != nil {
			return nil, err
		}
		cc.setExpire(keyID, m, time.Hour*24)
		return m, nil
	}

	return nil, errors.New("unable to find key")
}

// TokenInfoForProd dao
func TokenInfoForProd(IDToken string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(IDToken, getJWTKeys)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return &claims, nil
}
