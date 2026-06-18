package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"

	// config
	"github.com/Arvind215271/askito/internal/config"

	// logger
	"github.com/Arvind215271/askito/internal/logger"

	// api
	"github.com/Arvind215271/askito/internal/api"

	// youtube
	"github.com/Arvind215271/askito/internal/youtube"
	youtubeapi "github.com/Arvind215271/askito/internal/youtube/youtube_api"

	// transcript
	"github.com/Arvind215271/askito/internal/youtube/transcript"
	ytdlp "github.com/Arvind215271/askito/internal/youtube/transcript/providers/ytdlp"

	// export
	"github.com/Arvind215271/askito/internal/youtube/export"

	// debug
	"github.com/Arvind215271/askito/debug"
)
func main() {

	// get the config
	config := config.Load()

	// get the logger
	logger := logger.New(
		config.Env,
	)

	logger.Info(
		config.Env,
	)

	// get the error handler
	errorHandler := api.NewErrorHandler(
		logger,
	)

	// create echo instance
	e := echo.New()

	// ping route
	e.GET(
		"/ping",
		ping,
	)

	// context to be used by youtube API
	ctx := context.Background()

	// youtube

	youtubeClient, err := youtubeapi.NewClient(
		ctx,
		config.YouTubeAPIKey,
	)
	if err != nil {

		logger.Fatal(
			"failed to create youtube client",
			"error",
			err,
		)
	}
	

	youtubeProvider := youtubeapi.NewProvider(
		youtubeClient,
	)

	youtubeService := youtube.NewService(
		youtubeProvider,
	)

	// transcript

	ytdlpClient := &ytdlp.Client{}

	if err := ytdlpClient.ValidateYTDLP(); err != nil {

		logger.Warn(
			"yt-dlp unavailable",
			"error",
			err,
		)
	}

	transcriptProvider := ytdlp.NewProvider(
		ytdlpClient,
	)

	transcriptService := transcript.NewService(
		transcriptProvider,
	)

	// export

	exportService := export.NewService()

	exportService.RegisterExporter(
		export.FormatJSON,
		export.NewJSONExporter(),
	)

	// only run debug in development

	if config.Env == "development" {

		debug.DebugInput(
			ctx,
			logger,
			youtubeService,
			transcriptService,
			exportService,
		)
	}

	e.HTTPErrorHandler = errorHandler.Handle

	// start the server

	if err := e.Start(":8080"); err != nil {

		fmt.Println(
			"FAILED TO START THE SERVER",
			"ERROR:",
			err,
		)
	}
}



func ping(
	c *echo.Context,
) error {

	return c.JSON(
		http.StatusOK,
		"pong",
	)
} 