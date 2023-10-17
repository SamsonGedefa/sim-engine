package httpserver

import (
	"encoding/json"
	"fmt"
	"os"

	v1 "github.com/SamsonGedefa/simulator/main.go/api/v1"
	"github.com/labstack/echo/v4"
)

const (
	version = "v1"
)

type router struct {
	*echo.Echo
	version *echo.Group
}

func NewRouterV1(e *echo.Echo) router {
	prefix := fmt.Sprintf("/api/%s", version)
	version := e.Group(prefix)

	return router{
		Echo:    e,
		version: version,
	}
}

func TestRoute(r router, handler v1.Test) {
	m := r.version.Group("/test")
	m.GET("/", handler.Test)

}

func WriteRoutesToFile(r router) error {
	data, err := json.MarshalIndent(r.Routes(), "", "  ")
	if err != nil {
		return err
	}
	os.WriteFile("routes.json", data, 0644)

	return nil
}
