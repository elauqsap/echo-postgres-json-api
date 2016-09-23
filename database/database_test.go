package database

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type (
	ModelTest struct {
		Title string
		Func  func(*Store) func()
	}
	TestConfig struct {
		Config `json:"database"`
	}
)

var Conf TestConfig
var Data *Store

func TestDatabase(t *testing.T) {
	Convey("The Database Should", t, func() {
		Convey("Be Configurable From A JSON File", func() {
			data, err := ioutil.ReadFile("../configs/example.config.json")
			So(err, ShouldBeNil)
			So(json.Unmarshal(data, &Conf), ShouldBeNil)
			So(Conf, ShouldNotBeEmpty)
			Data, err = Conf.NewStore()
			So(err, ShouldBeNil)
			So(Data, ShouldNotBeNil)
		})
		Convey("Have Migrations For The Schema", func() {
			So(Data.Migrate(false), ShouldBeNil)
		})
	})
	var modelTests = []ModelTest{UserTest}
	for _, model := range modelTests {
		Convey("The Database "+model.Title, t, model.Func(Data))
	}
}
