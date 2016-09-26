package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elauqsap/echo-postgres-json-api/database"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	. "github.com/smartystreets/goconvey/convey"
)

type (
	TestCase struct {
		e            *echo.Echo
		url          string
		method       string
		content      string
		key          string
		handler      func(c echo.Context) error
		expectedCode int
		expectedBody string
	}
	HandlerTest struct {
		Title string
		Func  func(*echo.Echo, *Handlers) func()
	}
)

var Conf Config
var Store *database.Store

func TestAPI(t *testing.T) {
	Convey("The API Should", t, func() {
		Convey("Be Configurable From A JSON File", func() {
			data, err := ioutil.ReadFile("../configs/example.config.json")
			So(err, ShouldBeNil)
			So(json.Unmarshal(data, &Conf), ShouldBeNil)
			So(Conf, ShouldNotBeEmpty)
			Store, err = Conf.Database.NewStore()
			So(err, ShouldBeNil)
			So(Store, ShouldNotBeNil)
			So(Store.Migrate(false), ShouldBeNil)
		})
		data := &Data{Store}
		handlers := &Handlers{User: data}
		Convey("Seed Initial Test Users", func() {
			user1 := &database.User{
				First: "abc",
				Last:  "xyz",
				Role:  "admin",
				Key:   "WDpaAirzlzWCfuuMlexarniCdKIPeocr",
			}
			So(Store.ExecTransact(fmt.Sprintf("INSERT INTO app.users (first,last,role,api_key) VALUES ('%s','%s','%s','%s')", user1.First, user1.Last, user1.Role, user1.Key)), ShouldBeNil)
		})
		for _, api := range [][]HandlerTest{UserAPI} {
			for _, test := range api {
				Convey(test.Title, test.Func(echo.New(), handlers))
			}
		}
	})
}

func ExpectedResponse(test TestCase) func() {
	var req *http.Request
	var err error
	return func() {
		req, err = http.NewRequest(test.method, test.url, strings.NewReader(test.content))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		So(err, ShouldBeNil)
		Convey("Accept Them via The Model Handler", func() {
			rec := httptest.NewRecorder()
			context := test.e.NewContext(standard.NewRequest(req, test.e.Logger()), standard.NewResponse(rec, test.e.Logger()))
			err = test.handler(context)
			So(err, ShouldBeNil)
			Convey("And Receive Expected Response & Body", func() {
				So(rec.Code, ShouldEqual, test.expectedCode)
				So(rec.Body.String(), ShouldEqual, test.expectedBody)
			})
		})
	}
}
