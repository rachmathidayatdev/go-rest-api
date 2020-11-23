package base

import "github.com/labstack/echo"

// BaseCRUDController handle common create read update delete
type BaseCRUDController interface {
	Get(c echo.Context) error
	List(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}
