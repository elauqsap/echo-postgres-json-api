package api

import (
	"log"
	"os"

	"github.com/elauqsap/echo-postgres-json-api/database"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// Config ...
	Config struct {
		Bind  string `json:"bind"`
		Cert  string `json:"cert"`
		Key   string `json:"key"`
		Token string `json:"token"`
	}
	// Data embeds the backend so we can implement pointer methods for the API
	Data struct {
		*database.Store
		*log.Logger
	}
	// Handlers are the HTTP handlers to be used by the API router
	Handlers struct {
		User UserCRUD
	}
)

// New returns a configured api server instance
func (c *Config) New(sfile *os.File, dfile *os.File, store *database.Store) (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: sfile}))
	e.Use(middleware.Recover())
	api := e.Group("/api/v1")
	// api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: []byte(c.Token),
	// }))
	data := &Data{store, log.New(dfile, "", log.Ldate|log.Ltime|log.Lshortfile)}
	handlers := &Handlers{User: data}
	api.POST("/user", handlers.User.CreateUser)
	api.GET("/user/:id", handlers.User.ReadUser)
	api.PUT("/user/:id", handlers.User.UpdateUser)
	api.DELETE("/user/:id", handlers.User.DeleteUser)
	return e, nil
}
