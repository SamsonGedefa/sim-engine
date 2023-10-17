package v1

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewTest,
	),
	fx.Invoke(),
)