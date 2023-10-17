package httpserver

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewRouterV1,
	),

	fx.Invoke(
		TestRoute,
		WriteRoutesToFile,
	),
)
