package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/elauqsap/echo-postgres-json-api/database"
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
