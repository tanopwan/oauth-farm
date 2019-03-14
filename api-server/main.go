package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	e.GET("/api/products/product1", func(c echo.Context) error {
		product := struct {
			ID   string
			Name string
		}{
			ID:   "product1",
			Name: "product A",
		}
		return c.JSON(http.StatusOK, product)
	})

	e.GET("/api/products/product2", func(c echo.Context) error {
		xUserID := c.Request().Header.Get("X-User-Id")
		log.Printf("Headers: %+v\n", c.Request().Header)
		if xUserID != "0" {
			response := struct {
				Message string
			}{
				Message: "Forbidden user",
			}
			return c.JSON(http.StatusForbidden, response)
		}
		product := struct {
			ID   string
			Name string
		}{
			ID:   "product2",
			Name: "product B",
		}
		return c.JSON(http.StatusOK, product)
	})

	// Start server
	go func() {
		if err := e.Start(":8082"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
