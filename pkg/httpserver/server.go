package httpserver

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType,
			echo.HeaderAccept, echo.HeaderAuthorization,
		},
	}))
	// e.Use(middleware.CORS())

	e.HTTPErrorHandler = defaultHTTPErrorHandler

	return e
}

func defaultHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	c.Logger().Errorf("request error: %v", err)

	code := http.StatusInternalServerError
	message := any(http.StatusText(http.StatusInternalServerError))

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		message = he.Message
	}

	var be *echo.BindingError
	if errors.As(err, &be) {
		code = be.Code
		message = be.Message
	}

	if m, ok := message.(string); ok {
		message = map[string]any{"message": m}
	}

	// Send response
	if c.Request().Method == http.MethodHead { // Issue #608
		err = c.NoContent(code)
	} else {
		err = c.JSON(code, message)
	}
	if err != nil {
		c.Logger().Errorf("failed writing error response: %v", err)
	}
}

// # api:
// #   build:
// #     dockerfile: Dockerfile
// #     context: .
// #     target: dev
// #   depends_on:
// #     - db
// #   ports:
// #     - "8080:8080"
// #   volumes:
// #     - ./:/app
