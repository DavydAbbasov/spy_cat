package handler

import (
	"errors"
	"fmt"
	"net/http"

	dto "github.com/DavydAbbasov/spy-cat/internal/controllers/http/dto/mission"
	httperror "github.com/DavydAbbasov/spy-cat/internal/controllers/http/helpers"
	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"
	serviceserrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
	"github.com/rs/zerolog/log"

	missionservice "github.com/DavydAbbasov/spy-cat/internal/service/mission_service"
	"github.com/gin-gonic/gin"
)

type MissionHandler struct {
	missionSvc missionservice.MissionService
	validator  *validator.Validator
}

func NewMissionHandler(missionSvc missionservice.MissionService, validator *validator.Validator) *MissionHandler {
	return &MissionHandler{
		missionSvc: missionSvc,
		validator:  validator,
	}
}

// @Summary Create a new mission
// @Tags missions
// @Description Used to create a new mission
// @Accept json
// @Produce json
// @Param CreateMissionRequest body dto.CreateMissionRequest true "Request to create a mission"
// @Success 201 {object} dto.CreateMissionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 502 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /missions [post]
func (h *MissionHandler) CreateMission() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		req, err := validator.DecodeJSON[dto.CreateMissionRequest](h.validator, c.Request)
		if err != nil {
			if errors.Is(err, validator.ErrHandlerValidationFailed) {
				httperror.RespondError(c, http.StatusBadRequest, "invalid_body", err.Error())
				return
			}
			httperror.RespondError(c, http.StatusBadRequest, "invalid_json", "invalid json body")
			return
		}

		//dto -> domain
		m := dto.ToCreateMissionParams(*req)

		mission, err := h.missionSvc.CreateMission(ctx, m)
		if err != nil {
			switch {
			case errors.Is(err, serviceserrors.ErrInvalidCreateMission):
				httperror.RespondError(c, http.StatusBadRequest, "invalid_mission", "mission fields are invalid")
				return
			case errors.Is(err, serviceserrors.ErrMissionAlreadyExists):
				httperror.RespondError(c, http.StatusConflict, "already_exists", "mission with same title already exists")
				return
			case errors.Is(err, serviceserrors.ErrExternalService):
				httperror.RespondError(c, http.StatusBadGateway, "external_unavailable", "external dependency unavailable")
				return
			default:
				log.Error().Err(err).Msg("failed to create mission")
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
				return
			}
		}

		c.Header("Location", fmt.Sprintf("/missions/%d", mission.ID))
		c.JSON(http.StatusCreated, dto.CreateMissionResponse{ID: mission.ID})
	}
}
