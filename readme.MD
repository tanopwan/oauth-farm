## Authentication Server
### Debug & Run
```
env \
  GOOGLE_CLIENT_ID=[Client_ID] \
  GOOGLE_CLIENT_SECRET=[Client_Secret] \
  go run main.go server.go
```

### Build & Run
```
./build.sh

docker run \
  -e GOOGLE_CLIENT_ID=[Client_ID] \
  -e GOOGLE_CLIENT_SECRET=[Client_Secret] \
  -p 3000:3000 \
  oauth-farm/authentication-server
````