package config

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

func NewEndpointConnection() *subapps.HeLiCoPtErConnection {
	var connection *subapps.HeLiCoPtErConnection

	conn, err := subapps.PublishApp(
		fmt.Sprintf(
			"BEDROCK_ENDPOINT_URL=http://127.0.0.1:%s",
			strings.SplitN(viper.GetString("hypertext.bind_addr"), ":", 2)[1],
		),
		"matrix",
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to publish app into endpoint.")
	} else {
		connection = conn
		log.Info().Msg("Successfully published service into endpoints!")
	}

	return connection
}
