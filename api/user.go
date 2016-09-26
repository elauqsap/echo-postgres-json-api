package api

import (
	"net/http"

	"github.com/elauqsap/echo-postgres-json-api/database"
	"github.com/labstack/echo"
)

type (
	// UserCRUD ...
	UserCRUD interface {
		CreateUser(c echo.Context) error
		ReadUser(c echo.Context) error
		UpdateUser(c echo.Context) error
		DeleteUser(c echo.Context) error
	}
)

// CreateUser ...
func (d *Data) CreateUser(c echo.Context) error {
	bind := struct {
		First string `json:"first"`
		Last  string `json:"last"`
		Role  string `json:"role"`
	}{}
	if err := c.Bind(&bind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1010,
			"message": "invalid user or json format in request body",
		})
	}
	user := &database.User{First: bind.First, Last: bind.Last, Role: bind.Role}
	if err := d.ExecTransact(user.Create()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    1011,
			"message": "user could not be created",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1000,
		"message": "user was successfully created",
	})
}

// ReadUser ...
func (d *Data) ReadUser(c echo.Context) error {
	pm := new(database.PropertyMap)
	bind := struct {
		ID int `json:"id"`
	}{}
	if err := c.Bind(&bind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1010,
			"message": "invalid user or json format in request body",
		})
	}
	user := &database.User{ID: bind.ID}
	if err := d.SingleRowTransact(user.Read(), pm); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1012,
			"message": "user id does not exist",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1001,
		"message": pm,
	})
}

// UpdateUser ...
func (d *Data) UpdateUser(c echo.Context) error {
	bind := struct {
		ID    int    `json:"id"`
		First string `json:"first"`
		Last  string `json:"last"`
		Role  string `json:"role"`
	}{}
	if err := c.Bind(&bind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1010,
			"message": "invalid user or json format in request body",
		})
	}
	pm := new(database.PropertyMap)
	user := &database.User{ID: bind.ID, First: bind.First, Last: bind.Last, Role: bind.Role}
	if err := d.Merge(user, pm); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1013,
			"message": "user could not be updated",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1002,
		"message": "user successfully updated",
	})
}

// DeleteUser ...
func (d *Data) DeleteUser(c echo.Context) error {
	bind := struct {
		ID int `json:"id"`
	}{}
	if err := c.Bind(&bind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1010,
			"message": "invalid user or json format in request body",
		})
	}
	user := &database.User{ID: bind.ID}
	if err := d.ExecTransact(user.Delete()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    1014,
			"message": "user could not be deleted",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1003,
		"message": "user successfully deleted",
	})
}
