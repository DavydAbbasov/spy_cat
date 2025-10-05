package app

import (
	"net/http"

	pinghandler "github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers"
	cathandlers "github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers/cat"
	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers/swagger"
	logmiddleware "github.com/DavydAbbasov/spy-cat/internal/controllers/http/middleware"
	validator "github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"

	catservice "github.com/DavydAbbasov/spy-cat/internal/service/cat_service"

	"github.com/gin-gonic/gin"
)

func NewRouter(catSvc catservice.CatService) http.Handler {

	router := gin.Default()
	validator := validator.NewValidator()

	// middleware
	router.Use(logmiddleware.RequestResponseLogger())

	// handlers
	catHandler := cathandlers.NewCatHandler(catSvc, validator)

	// cats
	router.POST("/cats/create", catHandler.CreateCat())
	router.GET("/cats/:id", catHandler.GetCat())
	router.GET("/cats", catHandler.GetCats())
	// missions

	// swagger
	router.GET("/swagger/*any", swagger.Swagger())

	// ping
	router.GET("/ping", pinghandler.Ping())

	return router

}
