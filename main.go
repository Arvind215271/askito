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
	apiPlaylist "github.com/Arvind215271/askito/internal/api/playlist"
	apiSubtitle "github.com/Arvind215271/askito/internal/api/subtitle"
	apiTranscript "github.com/Arvind215271/askito/internal/api/transcript"

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

	//resource
	"github.com/Arvind215271/askito/internal/youtube/resource"
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

	// description
	descriptionService := description.NewService()

	// pipeline
	pipelineService := pipeline.NewService(youtubeService, descriptionService, subtitleService, transcriptService, signalService, logger, 2*config.PythonWorkers)

	// export service
	exportService := exportservice.NewService()

	exportService.RegisterExporter(
		exportservice.FormatJSON,
		&exportservice.JSONExporter{Pretty: true},
	)
	exportService.RegisterExporter(
		exportservice.FormatCSV,
		&exportservice.CSVExporter{},
	)
	exportService.RegisterExporter(
		exportservice.FormatMarkdown,
		&exportservice.MarkdownExporter{},
	)
	exportService.RegisterExporter(
		exportservice.FormatExcel,
		&exportservice.ExcelExporter{},
	)
	exportService.RegisterExporter(
		exportservice.FormatYAML,
		&exportservice.YAMLExporter{},
	)
	exportService.RegisterExporter(
		exportservice.FormatXML,
		&exportservice.XMLExporter{},
	)

	// resource service
	resourceService := resource.NewService(youtubeService, pipelineService, logger, 2*config.PythonWorkers)

	// handlers & routes
	exportHandler := export.NewHandler(resourceService, exportService)
	export.RegisterRoutes(e.Group("/export"), exportHandler)

	subtitleHandler := apiSubtitle.NewHandler(resourceService, subtitleService, exportService)
	apiSubtitle.RegisterSubtitleRoutes(e.Group("/subtitle"), subtitleHandler)

	transcriptHandler := apiTranscript.NewHandler(resourceService)
	e.POST("/transcript", transcriptHandler.GetTranscript)

	playlistHandler := apiPlaylist.NewHandler(youtubeService)
	apiPlaylist.RegisterPlaylistRoutes(e.Group("/playlist"), playlistHandler)

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
