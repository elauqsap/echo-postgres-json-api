package database

import (
	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

var user = &User{
	ID:    1,
	First: "First",
	Last:  "Last",
	Role:  USER,
}

var UserTest = ModelTest{
	Title: "User Model Should",
	Func: func(store *Store) func() {
		return func() {
			Convey("Implement The CRUD Interface", func() {
				So(&User{}, ShouldImplement, (*CRUD)(nil))
				Convey("A User Can Be Created", func() {
					So(store.Upsert(user, new(PropertyMap)), ShouldBeNil)
					read := new(User)
					pm := new(PropertyMap)
					So(store.SingleRowTransact(user.Read(), pm), ShouldBeNil)
					data, err := json.Marshal(pm)
					So(err, ShouldBeNil)
					So(json.Unmarshal(data, read), ShouldBeNil)
					So(read, ShouldResemble, user)
				})
				Convey("A User Can Be Read", func() {
					pm := new(PropertyMap)
					read := new(User)
					So(store.SingleRowTransact(user.Read(), pm), ShouldBeNil)
					data, err := json.Marshal(pm)
					So(err, ShouldBeNil)
					So(json.Unmarshal(data, read), ShouldBeNil)
					So(read, ShouldResemble, user)
				})
				Convey("A User Can Be Updated", func() {
					user.Role = ADMIN
					So(store.Upsert(user, new(PropertyMap)), ShouldBeNil)
					read := new(User)
					pm := new(PropertyMap)
					So(store.SingleRowTransact(user.Read(), pm), ShouldBeNil)
					data, err := json.Marshal(pm)
					So(err, ShouldBeNil)
					So(json.Unmarshal(data, read), ShouldBeNil)
					So(read, ShouldResemble, user)
				})
				Convey("A User Can Be Deleted", func() {
					pm := new(PropertyMap)
					So(store.ExecTransact(user.Delete()), ShouldBeNil)
					So(store.SingleRowTransact(user.Read(), pm).Error(), ShouldEqual, "sql: no rows in result set")
				})
			})
		}
	},
}
