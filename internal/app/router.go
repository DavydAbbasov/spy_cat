package app

import (
	pinghandler "github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers"
	cathandlers "github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers/cat"
	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers/swagger"
	logmiddleware "github.com/DavydAbbasov/spy-cat/internal/controllers/http/middleware"
	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"

	catservice "github.com/DavydAbbasov/spy-cat/internal/service/cat_service"

	"github.com/gin-gonic/gin"
)

func NewRouter(catSvc catservice.CatService) *gin.Engine {
	router := gin.Default()
	validator := validator.NewValidator()

	// middleware
	router.Use(logmiddleware.RequestResponseLogger())

	// handlers
	catHandler := cathandlers.NewCatHandler(catSvc, validator)

	// cats
	router.POST("/cats/create", catHandler.CreateCat())
	router.GET("/cats/:id", catHandler.GetCat())
	// missions

	// swagger
	router.GET("/", swagger.Swagger())
	// ping
	router.GET("/ping", pinghandler.Ping())

	return router

}
