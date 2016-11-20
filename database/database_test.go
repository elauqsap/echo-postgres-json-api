package database

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type (
	modelTest struct {
		Title string
		Func  func(*Store) func()
	}
	testConfig struct {
		Config `json:"database"`
	}
)

var conf testConfig
var store *Store

func TestDatabase(t *testing.T) {
	Convey("The Database Should", t, func() {
		Convey("Be Configurable From A JSON File", func() {
			data, err := ioutil.ReadFile("../configs/example.config.json")
			So(err, ShouldBeNil)
			So(json.Unmarshal(data, &conf), ShouldBeNil)
			So(conf, ShouldNotBeEmpty)
			store, err = conf.New()
			So(err, ShouldBeNil)
			So(data, ShouldNotBeNil)
		})
		Convey("Have Migrations For The Schema", func() {
			So(store.SetSchema(), ShouldBeNil)
		})
	})
	var modelTests = []modelTest{userTest}
	for _, model := range modelTests {
		Convey("The Database "+model.Title, t, model.Func(store))
	}
}
