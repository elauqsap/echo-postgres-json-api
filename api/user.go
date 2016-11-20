package api

import (
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

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
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&bind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1010,
			"message": "invalid user or json format in request body",
		})
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(bind.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1010,
			"message": "invalid user or json format in request body",
		})
	}
	user := &database.User{Login: bind.Login, Password: hashed}
	if err := d.Execute(user.Create()); err != nil {
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
	id, _ := strconv.Atoi(c.Param("id"))
	user := &database.User{ID: id}
	if err := d.Query(user.Read(), pm); err != nil {
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
		Login    string `json:"login,omitempty"`
		Password string `json:"password,omitempty"`
		JWT      string `json:"jwt,omitempty"`
	}{}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&bind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1010,
			"message": "invalid user or json format in request body",
		})
	}
	pairs := make(map[string]interface{})
	if len(bind.Login) > 0 {
		pairs["login"] = bind.Login
	}
	if len(bind.Password) > 0 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(bind.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    1010,
				"message": "invalid user or json format in request body",
			})
		}
		pairs["password"] = hashed
	}
	if len(bind.JWT) > 0 {
		pairs["jwt"] = bind.JWT
	}
	st := database.UpBuilder("users", "id", id, pairs)
	if _, err := d.Exec(st.Query, st.Args...); err != nil {
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
	// bind := struct {
	// 	ID int `json:"id"`
	// }{}
	// if err := c.Bind(&bind); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"code":    1010,
	// 		"message": "invalid user or json format in request body",
	// 	})
	// }
	// user := &database.User{ID: bind.ID}
	// if err := d.Execute(user.Delete()); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 		"code":    1014,
	// 		"message": "user could not be deleted",
	// 	})
	// }
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    1003,
		"message": "user successfully deleted",
	})
}
