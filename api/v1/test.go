package v1

import (
	"github.com/SamsonGedefa/simulator/main.go/api"
	"github.com/labstack/echo/v4"
)

type Test struct {
	api.Handler
}

func NewTest() Test {
	return Test{
		Handler: api.New(),
	}
}

func (t Test) Test(c echo.Context) error {

	return c.JSON(200, "Hello World")

}
