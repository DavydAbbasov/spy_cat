package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	dto "github.com/DavydAbbasov/spy-cat/internal/controllers/http/dto/mission"
	httperror "github.com/DavydAbbasov/spy-cat/internal/controllers/http/helpers"
	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"
	"github.com/DavydAbbasov/spy-cat/internal/domain"
	serviceerrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
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
			case errors.Is(err, serviceerrors.ErrInvalidCreateMission):
				httperror.RespondError(c, http.StatusBadRequest, "invalid_mission", "mission fields are invalid")
				return
			case errors.Is(err, serviceerrors.ErrMissionAlreadyExists):
				httperror.RespondError(c, http.StatusConflict, "already_exists", "mission with same title already exists")
				return
			case errors.Is(err, serviceerrors.ErrExternalService):
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

// @Summary Assign mission to a cat (or unassign with null)
// @Tags missions
// @Description Used to link or unlink a mission with a cat
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param AssignMissionRequest body dto.AssignMissionRequest true "Assignment payload"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /missions/{id}/assign [patch]
func (h *MissionHandler) AssignMission() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		missionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || missionID <= 0 {
			httperror.RespondError(c, http.StatusBadRequest, "invalid_id", "mission id must be positive integer")
			return
		}

		req, err := validator.DecodeJSON[dto.AssignMissionRequest](h.validator, c.Request)
		if err != nil {
			if errors.Is(err, validator.ErrHandlerValidationFailed) {
				httperror.RespondError(c, http.StatusBadRequest, "invalid_body", err.Error())
				return
			}
			httperror.RespondError(c, http.StatusBadRequest, "invalid_json", "invalid json body")
			return
		}

		err = h.missionSvc.AssignCat(ctx, missionID, req.CatID)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrMissionNotFound):
				httperror.RespondError(c, http.StatusNotFound, "mission_not_found", "mission not found")
				return
			case errors.Is(err, serviceerrors.ErrCatNotFound):
				httperror.RespondError(c, http.StatusNotFound, "cat_not_found", "cat not found")
				return
			default:
				log.Error().Err(err).Msg("assign mission failed")
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
				return
			}
		}

		c.Status(http.StatusNoContent)

	}
}

// Get a single mission
// @Summary      Get a single mission
// @Description  The ability to receive information about a single mission
// @Tags         missions
// @Param        id   path int true "ID Mission"
// @Produce      json
// @Success      200 {object} dto.MissionResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /mission/{id} [get]
func (h *MissionHandler) GetMission() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idStr := c.Param("id")
		if idStr == "" {
			httperror.RespondError(c, http.StatusBadRequest, "invalid_path", "id is required")
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			httperror.RespondError(c, http.StatusBadRequest, "invalid_path", "id must be a positive integer")
			return
		}

		mission, goals, err := h.missionSvc.GetMission(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, serviceerrors.ErrMissionNotFound):
				httperror.RespondError(c, http.StatusNotFound, "not_found", "mission not found ")
			default:
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			}
			return
		}
		c.JSON(http.StatusOK, dto.ToMissionResponse(mission, goals))
	}

}

// @Summary List missions
// @Tags missions
// @Produce json
// @Param status query string false "planned|active|completed"
// @Param catId  query int    false "Cat ID"
// @Param q      query string false "search by title"
// @Param limit  query int    false "limit (1..200)"
// @Param offset query int    false "offset"
// @Success 200 {object} dto.GetMissionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /missions [get]
func (h *MissionHandler) GetMissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var q dto.GetMissionsQuery

		if err := c.ShouldBindQuery(&q); err != nil {
			httperror.RespondError(c, http.StatusBadRequest, "invalid_query", err.Error())
			return
		}

		var f domain.MissionFilter
		if q.Status != nil && *q.Status != "" {
			st := domain.MissionStatus(*q.Status)
			f.Status = &st
		}
		f.CatID = q.CatID
		if q.Q != nil {
			if t := strings.TrimSpace(*q.Q); t != "" {
				f.Q = &t
			}
		}
		f.Limit = q.Limit
		f.Offset = q.Offset

		items, total, err := h.missionSvc.List(c.Request.Context(), f)
		if err != nil {
			httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			return
		}
		resp := dto.ToGetMissionsResponse(items, q.Limit, q.Offset, total)
		c.JSON(http.StatusOK, resp)
	}
}
