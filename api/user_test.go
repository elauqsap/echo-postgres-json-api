package api

import (
	"net/http"

	"github.com/labstack/echo"
)

var UserAPI = []HandlerTest{
	{
		"New User",
		func(e *echo.Echo, h Handlers) func() {
			test := TestCase{
				e:            e,
				url:          ``,
				method:       ``,
				content:      ``,
				key:          ``,
				handler:      h.User.CreateUser,
				expectedCode: http.StatusOK,
				expectedBody: ``,
			}
			return ExpectedResponse(test)
		},
	},
}
