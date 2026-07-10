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
	"github.com/Arvind215271/askito/internal/api/export"
	"github.com/Arvind215271/askito/internal/api/video"

	// youtube
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	youtubeapi "github.com/Arvind215271/askito/internal/youtube/metadata/youtube_api"
	ytdlpmetadata "github.com/Arvind215271/askito/internal/youtube/metadata/ytdlp"
	"github.com/Arvind215271/askito/internal/youtube/subtitle"

	// transcript
	"github.com/Arvind215271/askito/internal/youtube/transcript"

	// signal
	"github.com/Arvind215271/askito/internal/youtube/signal"

	// export
	exportservice "github.com/Arvind215271/askito/internal/youtube/export"

	// description
	"github.com/Arvind215271/askito/internal/youtube/description"

	// debug
	"github.com/Arvind215271/askito/debug"

	//pipeline
	"github.com/Arvind215271/askito/internal/youtube/pipeline"
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
		logger,
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
		logger,
	)

	ytdlpMetadataClient := ytdlpmetadata.NewClient(config.YtdlpCache, logger)
	ytdlpMetadataProvider := ytdlpmetadata.NewProvider(ytdlpMetadataClient, logger)

	// run cleanup on startup
	if err := ytdlpMetadataClient.Cleanup(); err != nil {
		logger.Error("failed to perform ytdlp cache cleanup", "error", err)
	} else {
		logger.Info("ytdlp cache cleanup completed successfully")
	}

	youtubeService := metadata.NewService(
		youtubeProvider,
		ytdlpMetadataProvider,
	)

	subtitleService := subtitle.NewSubtitleService()
	transcriptService := transcript.NewService()
	signalService := signal.NewSignalService()

	// video handler
	videoHandler := video.NewHandler(youtubeService, subtitleService, transcriptService, signalService)
	
	video.RegisterVideoRoutes(e.Group("/videos"), videoHandler)
	video.RegisterSubtitleRoutes(e.Group("/subtitles"), videoHandler)

    e.POST("/transcripts", videoHandler.GetTranscript)

	// description
	descriptionService := description.NewService()

	// pipeline
	pipelineService := pipeline.NewService(youtubeService, descriptionService, subtitleService, transcriptService, signalService)

	// export

	exportService := exportservice.NewService()

	exportService.RegisterExporter(
		exportservice.FormatJSON,
		exportservice.NewJSONExporter(),
	)

    // Export handler
    exportHandler := export.NewHandler(pipelineService, exportService)
    export.RegisterRoutes(e.Group("/export"), exportHandler)

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

	if err := e.Start(":" + config.Port); err != nil {
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
