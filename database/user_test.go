package database

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	. "github.com/smartystreets/goconvey/convey"
)

var users = []*User{
	&User{
		ID:       1,
		Login:    "user@test.com",
		Password: []byte("$2a$10$WqfTERA.eAu.0e5hXp2A7.7g84/qr7VA3qWjooT/Y3eGRIt0xYz9S"),
	},
	&User{
		ID:       2,
		Login:    "user1@test.com",
		Password: []byte("$2a$10$WqfTERA.eAu.0e5hXp2A7.7g84/qr7VA3qWjooT/Y3eGRIt0xYz9S"),
	},
}

var userTest = modelTest{
	Title: "User Model Should",
	Func: func(store *Store) func() {
		return func() {
			Convey("Implement The CRUD Interface", func() {
				So(&User{}, ShouldImplement, (*CRUD)(nil))
				Convey("A User Can Be Created", func() {
					for _, user := range users {
						So(store.Execute(user.Create()), ShouldBeNil)
					}
					read := new(User)
					pm := new(PropertyMap)
					So(store.Query(users[0].Read(), pm), ShouldBeNil)
					data, err := json.Marshal(pm)
					So(err, ShouldBeNil)
					So(json.Unmarshal(data, read), ShouldBeNil)
					So(read, ShouldResemble, &User{users[0].ID, users[0].Login, nil, users[0].JWT})
				})
				Convey("A User Login is Unique", func() {
					user := users[0].Create()
					_, err := store.Exec(user.Query, user.Args...)
					So(err, ShouldNotBeNil)
				})
				Convey("A User Can Be Read", func() {
					pm := new(PropertyMap)
					read := new(User)
					So(store.Query(users[0].Read(), pm), ShouldBeNil)
					data, err := json.Marshal(pm)
					So(err, ShouldBeNil)
					So(json.Unmarshal(data, read), ShouldBeNil)
					So(read, ShouldResemble, &User{users[0].ID, users[0].Login, nil, users[0].JWT})
				})
				Convey("A User Can Be Updated", func() {
					pairs := make(map[string]interface{})
					pairs["login"] = "updated@test.com"
					pairs["password"], _ = bcrypt.GenerateFromPassword([]byte("updated"), bcrypt.DefaultCost)
					pairs["jwt"] = "placeholder"
					st := UpBuilder("users", "id", 2, pairs)
					So(store.Execute(st), ShouldBeNil)
					read := new(User)
					pm := new(PropertyMap)
					So(store.Query(users[1].Read(), pm), ShouldBeNil)
					data, err := json.Marshal(pm)
					So(err, ShouldBeNil)
					So(json.Unmarshal(data, read), ShouldBeNil)
					So(read, ShouldResemble, &User{users[1].ID, "updated@test.com", nil, "placeholder"})
				})
				Convey("A User Can Be Deleted", func() {
					pm := new(PropertyMap)
					So(store.Execute(users[0].Delete()), ShouldBeNil)
					So(store.Query(users[0].Read(), pm).Error(), ShouldEqual, "sql: no rows in result set")
				})
			})
		}
	},
}
