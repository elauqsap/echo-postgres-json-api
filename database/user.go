package database

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const (
	// USER level role
	USER = `user`
	// MANAGER level role
	MANAGER = `manager`
	// ADMIN level role
	ADMIN = `admin`
)

type (
	// User models the app.users table
	User struct {
		ID    int    `json:"id"`
		First string `json:"first"`
		Last  string `json:"last"`
		Role  string `json:"role"`
		Key   string `json:"api_key"`
	}
)

// Create builds the Statement to insert a user into the database
func (u *User) Create() Statement {
	u.GenerateKey(32)
	if len(u.Role) <= 0 {
		u.Role = USER
	}
	return Statement{
		"INSERT INTO app.users (first,last,role,api_key) VALUES ($1,$2,$3,$4)",
		[]interface{}{u.First, u.Last, u.Role, u.Key},
	}
}

// Read creates the Statement to read a user from the database
func (u *User) Read() Statement {
	return Statement{
		"SELECT ROW_TO_JSON(u) FROM (SELECT * FROM app.users WHERE id = $1) AS u",
		[]interface{}{u.ID},
	}
}

// Update creates the Statement to update a user in the database
func (u *User) Update(v interface{}) Statement {
	// no merging needed to ignore v
	return Statement{
		"UPDATE app.users SET first = $1,last = $2,role = $3 WHERE id = $4",
		[]interface{}{u.First, u.Last, u.Role, u.ID},
	}
}

// Delete creates the Statement to delete a user from the database
func (u *User) Delete() Statement {
	return Statement{
		"DELETE FROM app.users WHERE id = $1",
		[]interface{}{u.ID},
	}
}

// GenerateKey creates a unique api key for each user
func (u *User) GenerateKey(n int) {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	u.Key = string(b)
}
