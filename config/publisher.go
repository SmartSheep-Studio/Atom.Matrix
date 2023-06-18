package config

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	tmodels "repo.smartsheep.studio/atom/nucleus/datasource/models"
	"repo.smartsheep.studio/atom/nucleus/toolbox"
)

func NewEndpointConnection(cycle fx.Lifecycle) *toolbox.ExternalServiceConnection {
	connection := &toolbox.ExternalServiceConnection{}

	cycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			conn, err := toolbox.PublishService(viper.GetString("general.endpoints"), viper.GetString("general.mount_key"), toolbox.ExternalServiceRequest{
				Name:        "LineupMarketplace",
				InstanceID:  viper.GetString("general.instance_id"),
				PackageID:   "repo.smartsheep.studio/atom/lineupmarketplace",
				Description: "A developer-friendly, geek-friendly store for apps and games.",
				Address:     viper.GetString("general.base_url"),
				Options: tmodels.ExternalServiceOptions{
					Pages:        []tmodels.ExternalPage{},
					Requirements: []string{"oauth"},
					Properties: fiber.Map{
						"oauth.urls":      []string{viper.GetString("general.base_url")},
						"oauth.callbacks": []string{fmt.Sprintf("%s/api/auth/callback", viper.GetString("general.base_url"))},
					},
				},
			})

			if err != nil {
				return err
			} else {
				connection.Configuration = conn.Configuration
				connection.Service = conn.Service
				connection.Additional = conn.Additional

				log.Info().Fields(connection.Service).Msg("Successfully published service into endpoints!")
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return toolbox.DepublishService(viper.GetString("general.endpoints"), viper.GetString("general.mount_key"), connection.Service.Secret)
		},
	})

	return connection
}
