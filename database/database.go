package database

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	// add postgres functionality to sql lib
	_ "github.com/lib/pq"
)

type (
	// Auth for postgres db
	Auth struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	// Config for postgres db
	Config struct {
		Auth
		Name           string `json:"name"`
		Host           string `json:"host"`
		Port           int    `json:"port"`
		SSL            string `json:"ssl_mode"`
		ConnectTimeout int    `json:"connect_timeout"`
	}
	// Store embeds an instance of a sql.DB so we can
	// inject new method receivers
	Store struct {
		*sql.DB
	}
	// PropertyMap allows to map and store JSON data with Postgres
	PropertyMap map[string]interface{}
	// CRUD interface implies the functions that need to be implemented
	// by a model in order to be used by the injected method receivers
	CRUD interface {
		Create() string
		Read() string
		Update(interface{}) string
		Delete() string
	}
)

// NewStore returns a configured *Store instance
func (c Config) NewStore() (*Store, error) {
	pg, err := Reverse(c.Auth)
	if err != nil {
		return nil, err
	}
	source := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s connect_timeout=%d", pg.A, pg.B, c.Name, c.Host, c.Port, c.SSL, c.ConnectTimeout)
	db, err := sql.Open("Postgres", source)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db}, nil
}

// ExecTransact performs a sql.Exec with a transaction
func (s *Store) ExecTransact(query string) error {
	return s.Transact(func(tx *sql.Tx) error {
		_, err := tx.Exec(query)
		return err
	})
}

// SingleRowTransact performs a sql.QueryRow with a transaction
func (s *Store) SingleRowTransact(query string, v interface{}) error {
	return s.Transact(func(tx *sql.Tx) error {
		return tx.QueryRow(query).Scan(v)
	})
}

// Transact allows for code reuse of a sql transaction
func (s *Store) Transact(txFunc func(*sql.Tx) error) (err error) {
	tx, err := s.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	return txFunc(tx)
}

// Upsert performs an insert/upsert of the model into the
// database depending on if a record is found
func (s *Store) Upsert(c CRUD, v interface{}) (err error) {
	if err = s.SingleRowTransact(c.Read(), v); err != nil {
		// Record was not found, create it
		err = s.ExecTransact(c.Create())
	} else {
		// Record was found, merge and update it
		var data []byte
		if data, err = json.Marshal(v); err == nil {
			err = s.ExecTransact(c.Update(data))
		}
	}
	return err
}

// Reverse performs reverse encryption of the user:pass
func Reverse(keys Auth) (*Auth, error) {
	c, err := base64.StdEncoding.DecodeString(keys.A)
	if err != nil {
		return nil, err
	}

	d, err := base64.StdEncoding.DecodeString(keys.B)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c)
	if err != nil {
		return nil, err
	}

	if len(d) < aes.BlockSize {
		return nil, errors.New("config a is too short")
	}

	iv := d[:aes.BlockSize]
	d = d[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(d, d)
	if ret := strings.Split(string(d), ":"); len(ret) == 2 {
		return &Auth{A: ret[0], B: ret[1]}, nil
	}
	return nil, errors.New("invalid structure to reverse")
}
