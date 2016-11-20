package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elauqsap/echo-postgres-json-api/reverse"

	// add postgres functionality to sql lib
	_ "github.com/lib/pq"
)

type (
	// Config for postgres db
	Config struct {
		Protected      reverse.Protected `json:"protected"`
		Name           string            `json:"name"`
		Host           string            `json:"host"`
		Port           int               `json:"port"`
		SSL            string            `json:"ssl_mode"`
		ConnectTimeout int               `json:"connect_timeout"`
	}
	// CRUD interface implies the functions that need to be implemented
	// by a model in order to be used by the injected method receivers
	CRUD interface {
		Create() Statement
		Read() Statement
		Delete() Statement
	}
	// PropertyMap allows to map and store JSON data with Postgres
	PropertyMap map[string]interface{}
	// Statement holds a SQL query and the params for exec/query
	Statement struct {
		Query string
		Args  []interface{}
	}
	// Store embeds an instance of a sql.DB so we can
	// inject new method receivers
	Store struct {
		*sql.DB
	}
)

// UpBuilder ...
func UpBuilder(table string, where string, value interface{}, pairs map[string]interface{}) Statement {
	// "UPDATE app.users SET login = $1, password = $2, jwt = $3 WHERE id = $4",
	var args []interface{}
	query := fmt.Sprintf("UPDATE app.%s SET", table)
	i := 1
	for k, v := range pairs {
		if i == len(pairs) {
			query = query + fmt.Sprintf(" %s = $%d", k, i)
			args = append(args, v)
		} else {
			query = query + fmt.Sprintf(" %s = $%d,", k, i)
			args = append(args, v)
		}
		i++
	}
	query = query + fmt.Sprintf(" WHERE %s = $%d", where, i)
	args = append(args, value)
	return Statement{query, args}
}

// New returns a configured *Store instance
func (c Config) New() (*Store, error) {
	creds, err := c.Protected.Reverse()
	if err != nil {
		return nil, err
	}
	source := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s connect_timeout=%d", creds.User, creds.Pass, c.Name, c.Host, c.Port, c.SSL, c.ConnectTimeout)
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	creds = nil
	return &Store{db}, nil
}

// Execute is execute with transaction
func (s *Store) Execute(st Statement) error {
	res, err := s.Exec(st.Query, st.Args...)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff < 1 {
		return errors.New("no change during execution")
	}
	return err
}

// Query is query with transaction
func (s *Store) Query(st Statement, v interface{}) error {
	return s.QueryRow(st.Query, st.Args...).Scan(v)
}

// Value map to a sql driver value
func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Scan map from sql return to a PropertyMap
func (p *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

//ToNullString invalidates a sql.NullString if empty, validates if not empty
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
