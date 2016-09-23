package database

import (
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type (
	// User ...
	User struct {
		ID    int    `json:"id"`
		First string `json:"first"`
		Last  string `json:"last"`
		Role  string `json:"role"`
		Key   string `json:"api_key"`
	}
)

// Create ...
func (u *User) Create() string {
	u.GenerateKey(32)
	return fmt.Sprintf("INSERT INTO app.users (first,last,role,key) VALUES ('%s','%s','%s','%s')", u.First, u.Last, u.Role, u.Key)
}

// Read ...
func (u *User) Read() string {
	return fmt.Sprintf("SELECT ROW_TO_JSON(user) FROM (SELECT (id,first,last,role) FROM app.users WHERE id='%d') as user", u.ID)
}

// Update ...
func (u *User) Update(x interface{}) string {
	// We don't have anything to merge so ignore x
	return fmt.Sprintf("UPDATE app.users SET first='%s',last='%s',role='%s' WHERE id='%d'", u.First, u.Last, u.Role, u.ID)
}

// Delete ...
func (u *User) Delete() string {
	return fmt.Sprintf("DELETE FROM app.users WHERE id='%d'", u.ID)
}

// GenerateKey ...
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
