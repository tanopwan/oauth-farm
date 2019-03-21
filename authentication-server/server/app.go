package server

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/tanopwan/oauth-farm/authentication-server/oauth"
	"github.com/tanopwan/oauth-farm/authentication-server/openid"
	"github.com/tanopwan/oauth-farm/authentication-server/user/repository/postgres"
	user "github.com/tanopwan/oauth-farm/authentication-server/user/service"
	"github.com/tanopwan/oauth-farm/common"
	"github.com/tanopwan/oauth-farm/common/session"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// App server
type App struct {
	userService    *user.Service
	sessionService *session.Service
}

// New App
func New(db *sql.DB) *App {
	return &App{
		userService:    user.NewService(postgres.NewRepository(db)),
		sessionService: session.NewService(common.NewRedisCache()),
	}
}

// GoogleLoginCode read authorization_code to get access_token and use access_token to get user info
// then validate user in our system
func (a *App) GoogleLoginCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/octet-stream") {
			return errors.New("invalid Content-Type header")
		}
		defer c.Request().Body.Close()
		bb, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: read request body error")
		}
		code := string(bb)
		clientID := os.Getenv("GOOGLE_CLIENT_ID")
		clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
		client := oauth.NewGoogleClient(clientID, clientSecret)

		token, err := client.TokenV4(code)
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: call TokenV4 error")
		}

		userInfo, err := client.UserInfoV1(token.AccessToken)
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: call UserInfoV1 error")
		}
		fmt.Printf("result: %+v\n", userInfo)

		user, err := a.userService.GetActiveUserByEmail(userInfo.Email)
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: user is not found or not active")
		}

		sess, err := a.sessionService.CreateSession(fmt.Sprintf("%d", user.ID))
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: create session error")
		}

		cookie := http.Cookie{
			Name:     "session",
			Value:    sess.Hash,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600,
			Domain:   "localhost",
			Secure:   true,
		}

		http.SetCookie(c.Response().Writer, &cookie)
		return c.String(http.StatusOK, "/users/new")
	}
}

// GoogleLoginToken read access_token and use access_token to get user info
// then validate user in our system
func (a *App) GoogleLoginToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/octet-stream") {
			return errors.New("invalid Content-Type header")
		}
		defer c.Request().Body.Close()
		bb, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: read request body error")
		}
		accessToken := string(bb)
		clientID := os.Getenv("GOOGLE_CLIENT_ID")
		clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
		client := oauth.NewGoogleClient(clientID, clientSecret)

		fmt.Printf("accessToken: %s\n", accessToken)
		userInfo, err := client.UserInfoV1(accessToken)
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: call UserInfoV1 error")
		}
		fmt.Printf("result: %+v\n", userInfo)

		user, err := a.userService.GetActiveUserByEmail(userInfo.Email)
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: user is not found or not active")
		}

		sess, err := a.sessionService.CreateSession(fmt.Sprintf("%d", user.ID))
		if err != nil {
			return errors.Wrap(err, "ctrlr googlelogin: create session error")
		}

		cookie := http.Cookie{
			Name:     "session",
			Value:    sess.Hash,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600,
			Domain:   "localhost",
			Secure:   true,
		}

		http.SetCookie(c.Response().Writer, &cookie)
		return c.String(http.StatusOK, "/users/new")
	}
}

// GoogleLoginOpenID read id_token and validate
// then validate user in our system
func (a *App) GoogleLoginOpenID() echo.HandlerFunc {
	n := "googleloginopenid"
	return func(c echo.Context) error {
		c.Logger().Debugf("GoogleLoginOpenID")
		idToken := c.Request().FormValue("data")
		claims, err := openid.TokenInfoForProd(idToken)
		if err != nil {
			return returnError(http.StatusUnauthorized, n, errors.Wrap(err, "validate id_token for prod error"))
		}

		if err := claims.Valid(); err != nil {
			return returnError(http.StatusUnauthorized, n, errors.Wrap(err, "claims is not valid"))
		}

		if !claims.VerifyExpiresAt(makeTimestamp(), true) {
			return returnError(http.StatusUnauthorized, n, errors.New("claims is expired"))
		}

		caud := os.Getenv("GOOGLE_CLIENT_ID")
		if aud, ok := (*claims)["aud"].(string); !ok || aud != caud {
			return returnError(http.StatusUnauthorized, n, errors.New("claims' aud is invalid found "+aud))
		}

		if iss, ok := (*claims)["iss"].(string); !ok || (iss != ciss1 && iss != ciss2) {
			return returnError(http.StatusUnauthorized, n, errors.New("claims' iss is invalid found "+iss))
		}

		email, ok := (*claims)["email"].(string)
		if !ok {
			return returnError(http.StatusUnauthorized, n, errors.New("claims' email is invalid"))
		}

		user, err := a.userService.GetActiveUserByEmail(email)
		if err != nil {
			return returnError(http.StatusUnauthorized, n, errors.Wrap(err, "user is not found or not active"))
		}

		sess, err := a.sessionService.CreateSession(fmt.Sprintf("%d", user.ID))
		if err != nil {
			return returnError(http.StatusInternalServerError, n, errors.Wrap(err, "create session error"))
		}

		cookie := http.Cookie{
			Name:     cCookieName,
			Value:    sess.Hash,
			Path:     cCookiePath,
			HttpOnly: true,
			MaxAge:   cCookieMaxAge,
			Domain:   cCookieDomain,
			Secure:   true,
		}

		http.SetCookie(c.Response().Writer, &cookie)
		return c.Redirect(http.StatusFound, "/")
	}
}

func (a *App) createEmptySession(c echo.Context) (*session.Model, error) {
	sess, err := a.sessionService.CreateSession("")
	if err != nil {
		return nil, errors.Wrap(err, "create session error")
	}

	cookie := http.Cookie{
		Name:     cCookieName,
		Value:    sess.Hash,
		Path:     cCookiePath,
		HttpOnly: true,
		MaxAge:   cCookieMaxAge,
		Domain:   cCookieDomain,
		Secure:   true,
	}

	http.SetCookie(c.Response().Writer, &cookie)
	return sess, nil
}

func (a *App) extractUserFromCookie(c echo.Context) (*session.Model, *user.Model, error) {
	cookie, err := c.Cookie(cCookieName)
	if err != nil || cookie == nil {
		// Cookie is not found
		sess, err := a.createEmptySession(c)
		if err != nil {
			return nil, nil, errors.Wrap(err, "create session error")
		}

		return sess, nil, nil
	}

	// Validate valid session from cookie
	userID, err := a.sessionService.ValidateSession(cookie.Value)
	if err != nil {
		// Failed to validate session
		sess, err := a.createEmptySession(c)
		if err != nil {
			return nil, nil, errors.Wrap(err, "create session error")
		}

		return sess, nil, nil
	}

	if userID == "" {
		// First time entering web-site or session is already expired
		sess, err := a.createEmptySession(c)
		if err != nil {
			return nil, nil, errors.Wrap(err, "create session error")
		}

		return sess, nil, nil
	}

	if userID == "empty" {
		c.Logger().Debugf("found empty session")
		return &session.Model{
			Hash:  cookie.Value,
			Value: userID,
		}, nil, nil
	}

	c.Logger().Debugf("found session with user %s", userID)
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.Logger().Errorf("parse user id error: %s", err.Error())
		return &session.Model{
			Hash:  cookie.Value,
			Value: userID,
		}, nil, nil
	}
	user, err := a.userService.GetUserByID(int(userIDInt))
	if err != nil {
		c.Logger().Errorf("get user error: %s", err.Error())
	} else if user == nil || user.ID == 0 {
		c.Logger().Debugf("user is not found")
	} else {
		return &session.Model{
			Hash:  cookie.Value,
			Value: userID,
		}, user, nil
	}
	return &session.Model{
		Hash:  cookie.Value,
		Value: userID,
	}, nil, nil
}

// RenderHTML render MPA
func (a *App) RenderHTML(c echo.Context) error {
	c.Logger().Debugf("RenderHTML")
	n := "renderhtml"

	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	data := struct {
		ClientID   string
		IsLoggedIn bool
		Title      string
		UserID     int
		Username   string
	}{
		ClientID:   clientID,
		IsLoggedIn: false,
		Title:      "Title",
	}

	_, user, err := a.extractUserFromCookie(c)
	if err != nil {
		return returnError(http.StatusInternalServerError, n, errors.Wrap(err, "extractUserFromCookie error"))
	}

	if user != nil {
		data.IsLoggedIn = true
		data.UserID = user.ID
		data.Username = user.Username
	}

	return c.Render(http.StatusOK, "html", data)
}

// Logout remove active session
func (a *App) Logout() echo.HandlerFunc {
	n := "logout"
	return func(c echo.Context) error {
		{
			cookie, err := c.Cookie(cCookieName)
			if err == nil && cookie != nil {
				a.sessionService.RemoveSession(cookie.Value)
			}
		}

		sess, err := a.sessionService.CreateSession("")
		if err != nil {
			return returnError(http.StatusInternalServerError, n, errors.Wrap(err, "create session error"))
		}

		cookie := http.Cookie{
			Name:     cCookieName,
			Value:    sess.Hash,
			Path:     cCookiePath,
			HttpOnly: true,
			MaxAge:   cCookieMaxAge,
			Domain:   cCookieDomain,
			Secure:   true,
		}

		http.SetCookie(c.Response().Writer, &cookie)
		return c.Redirect(http.StatusFound, "/")
	}
}

// ResolveSession from cookie to user
func (a *App) ResolveSession() echo.HandlerFunc {
	n := "resolvesession"
	return func(c echo.Context) error {
		cookie, err := c.Request().Cookie(cCookieName)
		if err != nil {
			return returnError(http.StatusUnauthorized, n, errors.Wrap(err, "get cookie error"))
		}

		session := cookie.Value
		userID, err := a.sessionService.ValidateSession(session)
		if err != nil {
			return returnError(http.StatusUnauthorized, n, errors.Wrap(err, "validate session error"))
		}

		c.Response().Header().Set("X-User-Id", userID)
		return c.Redirect(http.StatusFound, "/")
	}
}
