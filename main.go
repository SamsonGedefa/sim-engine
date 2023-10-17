package main

import (
	v1 "github.com/SamsonGedefa/simulator/main.go/api/v1"
	"github.com/SamsonGedefa/simulator/main.go/pkg/config"
	"github.com/SamsonGedefa/simulator/main.go/pkg/database"
	"github.com/SamsonGedefa/simulator/main.go/pkg/httpserver"
	"go.uber.org/fx"
)

func main() {

	fx.New(
		fx.Provide(
			config.Load,
			database.NewPGX,
			httpserver.New,
		),

		v1.Module,
		httpserver.Module,
		fx.Invoke(httpserver.Invoke),
	).Run()
}
