package server

import (
	"github.com/tanopwan/oauth-farm/common"
	"github.com/tanopwan/oauth-farm/common/session"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Server server
type Server struct {
	sessionService *session.Service
}

// New return server object
func New() *Server {
	return &Server{
		sessionService: session.NewService(common.NewRedisCache()),
	}
}

var defaultBackend = "http://localhost:3000"
var apiBackend = "http://localhost:8082"

// Start gateway server
func (s *Server) Start() {
	addr := ":8888"
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.handleHTTP(w, r)
		}),
	}

	log.Printf("Starting https server at %s\n", addr)
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}

func handleError(w http.ResponseWriter, status int, message string) {
	log.Printf("handleError: %s\n", message)
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func (s *Server) handleHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request URI: %s\n", req.RequestURI)
	fragments := strings.Split(strings.TrimPrefix(req.RequestURI, "/"), "/")
	log.Printf("Fragments 0: %s\n", fragments[0])
	var target string
	if fragments[0] == "api" {
		target = apiBackend
		s.writeAuthHeader(w, req)
	} else {
		target = defaultBackend
	}

	log.Printf("target: %s\n", target)
	s.serveReverseProxy(target, w, req)
}

func (s *Server) writeAuthHeader(w http.ResponseWriter, req *http.Request) {
	sess, err := req.Cookie("session")
	if err != nil {
		handleError(w, http.StatusUnauthorized, "session cookie is not found")
	}
	sessModel, err := s.sessionService.ValidateSession(sess.Value)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "session validate failed")
	}

	log.Printf("write user header: %s\n", sessModel.Value)
	req.Header.Set("X-User-Id", sessModel.Value)
}

func (s *Server) serveReverseProxy(target string, w http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(w, req)
}
