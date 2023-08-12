package hypertext

import (
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/controllers"
	"code.smartsheep.studio/atom/matrix/pkg/server/hypertext/middleware"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("http",
		middleware.Module(),
		controllers.Module(),

		fx.Invoke(fx.Annotate(MapControllers, fx.ParamTags(`group:"hypertext_controllers"`))),
		fx.Provide(NewHttpServer),
	)
}
