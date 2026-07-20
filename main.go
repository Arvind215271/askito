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
	apiSignal "github.com/Arvind215271/askito/internal/api/signal"
	apiSubtitle "github.com/Arvind215271/askito/internal/api/subtitle"
	apiTranscript "github.com/Arvind215271/askito/internal/api/transcript"
	apiVideo "github.com/Arvind215271/askito/internal/api/video"

	// cache
	"github.com/Arvind215271/askito/internal/cache"

	// youtube
	"github.com/Arvind215271/askito/internal/youtube/metadata"
	youtubeapi "github.com/Arvind215271/askito/internal/youtube/metadata/youtube_api"
	ytdlpmetadata "github.com/Arvind215271/askito/internal/youtube/metadata/ytdlp"
	"github.com/Arvind215271/askito/internal/youtube/metadata/ytdlp/python"
	ytSubtitle "github.com/Arvind215271/askito/internal/youtube/subtitle"

	// transcript
	ytTranscript "github.com/Arvind215271/askito/internal/youtube/transcript"

	// signal
	ytSignal "github.com/Arvind215271/askito/internal/youtube/signal"

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
		logger.Warn(
			"failed to create youtube client",
			"error",
			err,
		)
		youtubeClient = nil
	}

	youtubeProvider := youtubeapi.NewProvider(
		youtubeClient,
		logger,
	)

	// cache manager
	cacheManager := cache.NewManager(config.YtdlpCache, logger)

	// cleanup.
	cacheManager.Cleanup()

	pythonPool, err := python.NewSinglePool(config.PythonWorkers, logger, cacheManager)
	pythonPool.WarmUp(ctx)

	if err != nil {
		logger.Fatal("failed to create python pool", "error", err)
	}

	ytdlpMetadataClient := ytdlpmetadata.NewClient(pythonPool, logger)

	ytdlpMetadataProvider := ytdlpmetadata.NewProvider(ytdlpMetadataClient, logger)

	// validate ytdlp ig? IDK...

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

	subtitleService := ytSubtitle.NewSubtitleService(cacheManager, logger, pythonPool)
	transcriptService := ytTranscript.NewService()
	signalService := ytSignal.NewSignalService()

	// handlers
	videoHandler := apiVideo.NewHandler(youtubeService)
	subtitleHandler := apiSubtitle.NewHandler(youtubeService, subtitleService)
	transcriptHandler := apiTranscript.NewHandler(youtubeService, subtitleService, transcriptService)
	signalHandler := apiSignal.NewHandler(youtubeService, subtitleService, transcriptService, signalService)

	// routes
	apiVideo.RegisterVideoRoutes(e.Group("/videos"), videoHandler)
	apiSubtitle.RegisterSubtitleRoutes(e.Group("/subtitles"), subtitleHandler)
	e.POST("/transcripts", transcriptHandler.GetTranscript)
	e.POST("/signals", signalHandler.GetVideoSignals)

	// description
	descriptionService := description.NewService()

	// pipeline
	pipelineService := pipeline.NewService(youtubeService, descriptionService, subtitleService, transcriptService, signalService, logger, 2*config.PythonWorkers)

	// export

	exportService := exportservice.NewService()

	exportService.RegisterExporter(
		exportservice.FormatJSON,
		exportservice.NewJSONExporter(),
	)

	// Export handler
	exportHandler := export.NewHandler(youtubeService, pipelineService, exportService)
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
