package feature_test

import (
	"github.com/labstack/echo"
	"github.com/tanopwan/oauth-farm/authentication-server/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGoogleLoginWithOpenIDSuccess(t *testing.T) {
	app := server.New()

	payload := ""
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/openid/login/google", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, "application/octet-stream; charset=utf-8")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := app.GoogleLoginOpenID()

	err := h(c)
	if err != nil {
		t.Logf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	if rec.Code != 200 {
		t.Errorf("expect 200, actual: %d", rec.Code)
	}

	t.Logf("response body: %s", rec.Body.String())
}
