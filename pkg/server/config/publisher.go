package config

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func NewEndpointConnection() *subapps.HeLiCoPtErConnection {
	var connection *subapps.HeLiCoPtErConnection

	conn, err := subapps.PublishApp(
		viper.GetString("base_url"),
		"matrix",
		models.SubAppExposedPage{
			Icon:  "mdi-store",
			Name:  "matrix",
			Title: "Matrix",
			Path:  "/",
			Meta: map[string]any{
				"gatekeeper": map[string]any{
					"must": true,
				},
			},
		},
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to publish app into endpoint.")
	} else {
		connection = conn
		log.Info().Msg("Successfully published service into endpoints!")
	}

	return connection
}
