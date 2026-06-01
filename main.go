package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v5"

	// config
	"github.com/Arvind215271/askito/internal/config"

	// logger
	"github.com/Arvind215271/askito/internal/logger"

	// api
	"github.com/Arvind215271/askito/internal/api"
)


func main() {
	// get the config
	config := config.Load()

	// get the logger
	logger := logger.New(config.Env)

	// get the error handler
	errorHandler := api.NewErrorHandler(logger)


	// create echo instance
	e := echo.New()

	// ping route.
	e.GET("/ping", ping)

	// use the error handler instead of internal error handler
	e.HTTPErrorHandler = errorHandler.Handle


	// start the server
	if err := e.Start(":8080"); err != nil {
		fmt.Println("FAILED TO START THE SERVER", "ERROR: ",err)
	}
		


}


func ping(c *echo.Context) error {
	return 	c.JSON(http.StatusOK, "pong")

}