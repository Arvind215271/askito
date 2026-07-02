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
	"github.com/Arvind215271/askito/internal/api/video"

	// youtube
	"github.com/Arvind215271/askito/internal/youtube"
	youtubeapi "github.com/Arvind215271/askito/internal/youtube/metadata/youtube_api"
	ytdlpmetadata "github.com/Arvind215271/askito/internal/youtube/metadata/ytdlp"

	// transcript
	"github.com/Arvind215271/askito/internal/youtube/transcript"
	ytdlptranscript "github.com/Arvind215271/askito/internal/youtube/transcript/providers/ytdlp"

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

	ytdlpMetadataClient := ytdlpmetadata.NewClient()
	ytdlpMetadataProvider := ytdlpmetadata.NewProvider(ytdlpMetadataClient)

	youtubeService := youtube.NewService(
		youtubeProvider,
		ytdlpMetadataProvider,
	)

	// video handler
	videoHandler := video.NewHandler(youtubeService)
	e.GET("/videos/:id", videoHandler.GetVideoByID)

	// transcript

	ytdlpTranscriptClient := &ytdlptranscript.Client{}

	if err := ytdlpTranscriptClient.ValidateYTDLP(); err != nil {

		logger.Warn(
			"yt-dlp unavailable",
			"error",
			err,
		)
	}

	transcriptProvider := ytdlptranscript.NewProvider(
		ytdlpTranscriptClient,
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
	if config.Env == "dev" {

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
