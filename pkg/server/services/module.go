package services

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("services",
		fx.Provide(NewExploreService),
	)
}
