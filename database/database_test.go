package database

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type (
	ModelTest struct {
		Title string
		Func  func(*Store) func()
	}
	TestConfig struct {
		Config `json:"databse"`
	}
)

var Conf TestConfig
var Data *Store

func TestDatabase(t *testing.T) {
	Convey("", t, func() {

	})
}
