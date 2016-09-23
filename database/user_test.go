package database

import (
	. "github.com/smartystreets/goconvey/convey"
)

var UserTest = ModelTest{
	Title: "User Model Should",
	Func: func(store *Store) func() {
		return func() {
			Convey("Implement The CRUD Interface", func() {
				So(&User{}, ShouldImplement, (*CRUD)(nil))
			})
		}
	},
}
