package database

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"database/sql/driver"
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
	// CRUD interface implies the functions that need to be implemented
	// by a model in order to be used by the injected method receivers
	CRUD interface {
		Create() Statement
		Read() Statement
		Update(interface{}) Statement
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

// NewStore returns a configured *Store instance
func (c Config) NewStore() (*Store, error) {
	pg, err := Reverse(c.Auth)
	if err != nil {
		return nil, err
	}
	source := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s connect_timeout=%d", pg.A, pg.B, c.Name, c.Host, c.Port, c.SSL, c.ConnectTimeout)
	db, err := sql.Open("postgres", source)
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

// EWT is execute with transaction
func (s *Store) EWT(st Statement) error {
	return s.Transact(func(tx *sql.Tx) error {
		res, err := tx.Exec(st.Query, st.Args...)
		if aff, _ := res.RowsAffected(); aff < 1 {
			return errors.New("no change during execution")
		}
		return err
	})
}

// QWT is query with transaction
func (s *Store) QWT(st Statement, v interface{}) error {
	return s.Transact(func(tx *sql.Tx) error {
		return tx.QueryRow(st.Query, st.Args...).Scan(v)
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
