package feature_test

import (
	appdb "github.com/tanopwan/oauth-farm/authentication-server/db"
	"github.com/tanopwan/oauth-farm/authentication-server/server"
	"testing"
)

var app *server.App

func TestMain(t *testing.T) {
	db, close, err := appdb.DialPostgresDB()
	if err != nil {
		panic(err)
	}
	defer close()

	app = server.New(db)
}
