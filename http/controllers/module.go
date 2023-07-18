package controllers

import (
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"go.uber.org/fx"
)

type HttpController interface {
	Map(router *ctx.App)
}

func AsController(f any) any {
	return fx.Annotate(f, fx.As(new(HttpController)), fx.ResultTags(`group:"http"`))
}

func Module() fx.Option {
	return fx.Module("http.controllers",
		fx.Provide(
			AsController(NewStatusController),
			AsController(NewExploreController),
			AsController(NewLibraryController),
			AsController(NewAppController),
			AsController(NewReleaseController),
			AsController(NewPostController),
			AsController(NewCloudSaveController),
		),
	)
}
