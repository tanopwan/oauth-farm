package main

import (
	"context"
	"github.com/labstack/echo"
	appdb "github.com/tanopwan/oauth-farm/authentication-server/db"
	"github.com/tanopwan/oauth-farm/authentication-server/server"
	"html/template"
)

// Server server
type Server struct {
	e     *echo.Echo
	a     *server.App
	close func()
}

// New server
func New() *Server {
	db, c, err := appdb.DialPostgresDB()
	if err != nil {
		panic(err)
	}

	return &Server{
		e:     echo.New(),
		a:     server.New(db),
		close: c,
	}
}

func (s *Server) start() {
	s.e.Logger.Info("starting server ...")
	s.e.Start(":3000")
}

func (s *Server) startTLS() {
	s.e.Logger.Info("starting server ...")
	s.e.StartTLS(":3000", "cert.pem", "key.pem")
}

func (s *Server) shutdown(ctx context.Context) error {
	s.close()
	return s.e.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	s.e.Static("/", "public")
	s.e.Logger.SetLevel(1)
	// s.e.Use(middleware.Logger())
	g := s.e.Group("/login")
	g.POST("/auth/google", s.a.GoogleLoginCode())
	g.POST("/token/google", s.a.GoogleLoginToken())

	s.e.POST("/openid/login/google", s.a.GoogleLoginOpenID())
	s.e.POST("/logout", s.a.Logout())

	s.e.Renderer = server.NewTemplateRenderer(template.Must(template.ParseGlob("templates/*.tmpl")))
	s.e.GET("/", s.a.RenderHTML)
}
