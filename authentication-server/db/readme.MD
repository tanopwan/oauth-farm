### How to use goose db migrations ###
```
cd oauth-farm/authentication-server
goose -dir "db/migrations/" postgres "user=tanopwan dbname=oauthfarm sslmode=disable" status
```