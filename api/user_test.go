package api

import (
	"net/http"

	"github.com/labstack/echo"
)

var UserAPI = []HandlerTest{
	{
		"Create A New User (Valid User)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.POST,
				content:      `{"first":"abc","last":"xyz","role":"user"}`,
				handler:      h.User.CreateUser,
				expectedCode: http.StatusOK,
				expectedBody: `{"code":1000,"message":"user was successfully created"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Create A New User (Invalid User)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.POST,
				content:      `{"last":"xyz","role":"user"}`,
				handler:      h.User.CreateUser,
				expectedCode: http.StatusInternalServerError,
				expectedBody: `{"code":1011,"message":"user could not be created"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Create A New User (Invalid Input)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.POST,
				content:      `invalid_content`,
				handler:      h.User.CreateUser,
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"code":1010,"message":"invalid user or json format in request body"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Read A User (Valid User)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.GET,
				content:      `{"id":1}`,
				handler:      h.User.ReadUser,
				expectedCode: http.StatusOK,
				expectedBody: `{"code":1001,"message":{"api_key":"WDpaAirzlzWCfuuMlexarniCdKIPeocr","first":"abc","id":1,"last":"xyz","role":"admin"}}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Read A User (Invalid User)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.GET,
				content:      `{"id":1000}`,
				handler:      h.User.ReadUser,
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"code":1012,"message":"user id does not exist"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Read A User (Invalid Input)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.GET,
				content:      `invalid_content`,
				handler:      h.User.ReadUser,
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"code":1010,"message":"invalid user or json format in request body"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Update A User (Valid Update)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.PUT,
				content:      `{"first":"xyz","id":1,"last":"abc","role":"manager"}`,
				handler:      h.User.UpdateUser,
				expectedCode: http.StatusOK,
				expectedBody: `{"code":1002,"message":"user successfully updated"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Update A User (Invalid Update)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.PUT,
				content:      `{"first":"xxyyzz","id":1000,"last":"abc","role":"manager"}`,
				handler:      h.User.UpdateUser,
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"code":1013,"message":"user could not be updated"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Update A User (Invalid Input)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.PUT,
				content:      `invalid_content`,
				handler:      h.User.UpdateUser,
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"code":1010,"message":"invalid user or json format in request body"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Delete A User (Valid User)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.DELETE,
				content:      `{"id":2}`,
				handler:      h.User.DeleteUser,
				expectedCode: http.StatusOK,
				expectedBody: `{"code":1003,"message":"user successfully deleted"}`,
			}
			return ExpectedResponse(test)
		},
	},
	{
		"Delete A User (Invalid Input)",
		func(e *echo.Echo, h *Handlers) func() {
			test := TestCase{
				e:            e,
				url:          `/api/v1/user`,
				method:       echo.DELETE,
				content:      `invalid_content`,
				handler:      h.User.DeleteUser,
				expectedCode: http.StatusBadRequest,
				expectedBody: `{"code":1010,"message":"invalid user or json format in request body"}`,
			}
			return ExpectedResponse(test)
		},
	},
}
