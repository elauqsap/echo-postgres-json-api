package api

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/elauqsap/echo-postgres-json-api/database"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// Config for the server and database
	Config struct {
		Server struct {
			Bind    string `json:"bind"`
			Cert    string `json:"cert"`
			Key     string `json:"key"`
			LogPath string `json:"log,omitempty"`
		} `json:"server"`
		Database database.Config `json:"database"`
	}
	// Data ...
	Data struct {
		*database.Store
	}
	// Handlers ...
	Handlers struct {
		User UserCRUD
	}
)

// NewConfig loads the configurations into a structure
// to be used as a global operator at runtime
func NewConfig(path string) (conf *Config) {
	data, _ := ioutil.ReadFile(path)
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil
	}
	return conf
}

// NewServer returns a configured api server instance
func (c *Config) NewServer() (*echo.Echo, error) {
	e := echo.New()
	if logFile, err := os.OpenFile(c.Server.LogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660); err != nil {
		e.Use(middleware.Logger())
	} else {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: logFile}))
	}
	e.Use(middleware.Recover())
	api := e.Group("/api/v1")
	if store, err := c.Database.NewStore(); err == nil {
		data := &Data{store}
		handlers := &Handlers{User: data}
		api.POST("/user", handlers.User.CreateUser)
		api.GET("/user", handlers.User.ReadUser)
		api.PUT("/user", handlers.User.UpdateUser)
		api.DELETE("/user", handlers.User.DeleteUser)
	} else {
		return nil, err
	}
	return e, nil
}
