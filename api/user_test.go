package api

import "github.com/labstack/echo"

var UserAPI = []HandlerTest{
	{
		"New User",
		func(e *echo.Echo) func() {
			test := TestCase{}
			return ExpectedResponse(test)
		},
	},
}
