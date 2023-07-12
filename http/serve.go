package http

import (
	ctx "code.smartsheep.studio/atom/neutron/http/context"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/matrix/http/middleware"
	"code.smartsheep.studio/atom/matrix/renderer"
	"code.smartsheep.studio/atom/neutron/toolbox"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2/middleware/filesystem"
	flog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var conn *toolbox.ExternalServiceConnection
var server *ctx.App

func NewHttpServer(cycle fx.Lifecycle, cors middleware.CorsHandler, conf *viper.Viper, c *toolbox.ExternalServiceConnection) *ctx.App {
	conn = c

	// Create app
	server = &ctx.App{P: fiber.New(fiber.Config{
		Prefork:               viper.GetBool("http.advanced.prefork"),
		CaseSensitive:         false,
		StrictRouting:         false,
		DisableStartupMessage: true,
		ServerHeader:          "Matrix",
		AppName:               "Matrix v2.0",
		BodyLimit:             viper.GetInt("http.max_body_size"),
	}),
	}

	// Apply global middleware
	server.P.Use(flog.New(flog.Config{
		Format: "${status} | ${latency} | ${method} ${path} ${body}\n",
		Output: log.Logger,
	}))
	server.P.Use(cors())

	cycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.P.Listen(conf.GetString("http.listen_addr"))
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to start http http")
				}
			}()

			return nil
		},
	})

	return server
}

func MapControllers(controllers []HttpController, server *ctx.App) {
	for _, controller := range controllers {
		controller.Map(server)
	}

	// Fallback not found api to nucleus
	server.All("/api/*", func(c *ctx.Ctx) error {
		uri := fmt.Sprintf("%s?%s", c.P.Request().URI().Path(), c.P.Request().URI().QueryArgs().String())
		return c.P.Redirect(conn.GetEndpointPath(uri), fiber.StatusFound)
	})

	// Serve static files
	server.P.Use("/", filesystem.New(filesystem.Config{
		Root:         renderer.GetHttpFS(),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))
}
