package api

import "github.com/labstack/echo"

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
	return nil
}

// ReadUser ...
func (d *Data) ReadUser() error {
	return nil
}

// UpdateUser ...
func (d *Data) UpdateUser() error {
	return nil
}

// DeleteUser ...
func (d *Data) DeleteUser() error {
	return nil
}
