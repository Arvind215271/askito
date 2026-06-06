package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v5"
	

	"context"


	// config
	"github.com/Arvind215271/askito/internal/config"

	// logger
	"github.com/Arvind215271/askito/internal/logger"

	// api
	"github.com/Arvind215271/askito/internal/api"

	// youtube
	"github.com/Arvind215271/askito/internal/youtube"
	
	
)


func main() {
	// get the config
	config := config.Load()

	// get the logger
	logger := logger.New(config.Env)

	logger.Info(config.Env)

	// get the error handler
	errorHandler := api.NewErrorHandler(logger)


	// create echo instance
	e := echo.New()

	// ping route.
	e.GET("/ping", ping)


	// context to be used by youtube API
	ctx := context.Background()

	youtubeClient, err := youtube.NewClient(
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

	// only run debug in development
	if config.Env == "development" {
		youtube.DebugYouTube(
			ctx,
			logger,
			youtubeClient,
		)
	}
	
	e.HTTPErrorHandler = errorHandler.Handle


	// start the server
	if err := e.Start(":8080"); err != nil {
		fmt.Println("FAILED TO START THE SERVER", "ERROR: ",err)
	}
}


func ping(c *echo.Context) error {
	return 	c.JSON(http.StatusOK, "pong")

}