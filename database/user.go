package database

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
		ID       int    `json:"id,omitempty"`
		Login    string `json:"login,omitempty"`
		Password []byte `json:"password,omitempty"`
		JWT      string `json:"jwt,omitempty"`
	}
)

// Create builds the Statement to insert a user into the database
func (u *User) Create() Statement {
	return Statement{
		"INSERT INTO app.users (login,password,jwt) VALUES ($1,$2,$3)",
		[]interface{}{u.Login, u.Password, ToNullString(u.JWT)},
	}
}

// 		"SELECT ROW_TO_JSON(u) FROM (SELECT id, login, encode(password, 'escape') as password, jwt FROM app.users WHERE id = $1) u",

// Read creates the Statement to read a user from the database
func (u *User) Read() Statement {
	return Statement{
		"SELECT ROW_TO_JSON(u) FROM (SELECT id, login, jwt FROM app.users WHERE id = $1) u",
		[]interface{}{u.ID},
	}
}

// Update creates the Statement to update a user in the database
func (u *User) Update(v interface{}) Statement {
	// no merging needed to ignore v
	return Statement{
		"UPDATE app.users SET login = $1, password = $2, jwt = $3 WHERE id = $4",
		[]interface{}{u.Login, u.Password, ToNullString(u.JWT), u.ID},
	}
}

// Delete creates the Statement to delete a user from the database
func (u *User) Delete() Statement {
	return Statement{
		"DELETE FROM app.users WHERE id = $1",
		[]interface{}{u.ID},
	}
}
