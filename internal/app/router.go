package app

import (
	"net/http"

	pinghandler "github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers"
	cathandlers "github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers/cat"
	missionhandlers "github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers/mission"

	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/handlers/swagger"
	logmiddleware "github.com/DavydAbbasov/spy-cat/internal/controllers/http/middleware"
	validator "github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"

	catservice "github.com/DavydAbbasov/spy-cat/internal/service/cat_service"
	missionservice "github.com/DavydAbbasov/spy-cat/internal/service/mission_service"

	"github.com/gin-gonic/gin"
)

func NewRouter(catSvc catservice.CatService, missionSvc missionservice.MissionService) http.Handler {

	router := gin.Default()
	validator := validator.NewValidator()

	// middleware
	router.Use(logmiddleware.RequestResponseLogger())

	// handlers
	catHandler := cathandlers.NewCatHandler(catSvc, validator)
	missionHandler := missionhandlers.NewMissionHandler(missionSvc, validator)

	// cats
	router.POST("/cats/create", catHandler.CreateCat())
	router.GET("/cats/:id", catHandler.GetCat())
	router.GET("/cats", catHandler.GetCats())
	router.DELETE("/cats/:id", catHandler.DeleteCat())
	router.PATCH("/cats/:id/salary", catHandler.UpdateSalary())

	// missions
	router.POST("/missions", missionHandler.CreateMission())
	router.PATCH("/missions/:id/assign", missionHandler.AssignMission())
	router.GET("/mission/:id", missionHandler.GetMission())
	router.GET("/missions", missionHandler.GetMissions())
	router.PATCH("/missions/:id/status", missionHandler.UpdateMissionStatus())
	router.POST("/missions/:id/goals", missionHandler.AddGoal())

	// swagger
	router.GET("/swagger/*any", swagger.Swagger())

	// ping
	router.GET("/ping", pinghandler.Ping())

	return router

}
