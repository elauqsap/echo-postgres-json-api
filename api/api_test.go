package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
		request      request
		handler      func(c echo.Context) error
		expectedCode int
		expectedBody string
	}
	HandlerTest struct {
		Title string
		Func  func(*echo.Echo, *Handlers) func()
	}
)

type request struct {
	url         string
	method      string
	content     string
	paramNames  []string
	paramValues []string
}

// config contains the global settings for the server and database
type config struct {
	Server   Config          `json:"server"`
	Database database.Config `json:"database"`
}

var conf config
var store *database.Store

func TestAPI(t *testing.T) {
	Convey("The API Should", t, func() {
		Convey("Be Configurable From A JSON File", func() {
			data, err := ioutil.ReadFile("../configs/example.config.json")
			So(err, ShouldBeNil)
			So(json.Unmarshal(data, &conf), ShouldBeNil)
			So(conf, ShouldNotBeEmpty)
			store, err = conf.Database.New()
			So(err, ShouldBeNil)
			So(store, ShouldNotBeNil)
			So(store.SetSchema(), ShouldBeNil)
		})
		data := &Data{store, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)}
		handlers := &Handlers{User: data}
		Convey("Seed Initial Test Users", func() {
			user1 := &database.User{
				Login:    "user@test.com",
				Password: []byte("$2a$10$WqfTERA.eAu.0e5hXp2A7.7g84/qr7VA3qWjooT/Y3eGRIt0xYz9S"),
			}
			st := database.Statement{
				Query: "INSERT INTO app.users (login,password,jwt) VALUES ($1,$2,$3)",
				Args:  []interface{}{user1.Login, user1.Password, database.ToNullString(user1.JWT)},
			}
			So(store.Execute(st), ShouldBeNil)
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
		req, err = http.NewRequest(test.request.method, test.request.url, strings.NewReader(test.request.content))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		So(err, ShouldBeNil)
		Convey("Accept Them via The Model Handler", func() {
			rec := httptest.NewRecorder()
			context := test.e.NewContext(standard.NewRequest(req, test.e.Logger()), standard.NewResponse(rec, test.e.Logger()))
			context.SetParamNames(test.request.paramNames...)
			context.SetParamValues(test.request.paramValues...)
			err = test.handler(context)
			So(err, ShouldBeNil)
			Convey("And Receive Expected Response & Body", func() {
				So(rec.Code, ShouldEqual, test.expectedCode)
				So(rec.Body.String(), ShouldEqual, test.expectedBody)
			})
		})
	}
}
