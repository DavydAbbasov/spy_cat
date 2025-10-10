package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	serviceserrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
	"github.com/gin-gonic/gin"

	dto "github.com/DavydAbbasov/spy-cat/internal/controllers/http/dto/cat"
	httperror "github.com/DavydAbbasov/spy-cat/internal/controllers/http/helpers"
	catservice "github.com/DavydAbbasov/spy-cat/internal/service/cat_service"

	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"
)

type CatHandler struct {
	svc       catservice.CatService
	validator *validator.Validator
}

func NewCatHandler(svc catservice.CatService, validator *validator.Validator) *CatHandler {
	return &CatHandler{
		svc:       svc,
		validator: validator,
	}
}

// @Summary Create a new spy cat
// @Tags cats
// @Description Used to create a new spy cat
// @Accept json
// @Produce json
// @Param CreateCatRequest body  dto.CreateCatRequest true "Request to create a cat"
// @Success 201 {object} dto.CreateCatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cats/create [post]
func (h *CatHandler) CreateCat() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		req, err := validator.DecodeJSON[dto.CreateCatRequest](h.validator, c.Request)
		if err != nil {
			if errors.Is(err, validator.ErrHandlerValidationFailed) {
				httperror.RespondError(c, http.StatusBadRequest, "invalid_body", err.Error())
				return
			}
			httperror.RespondError(c, http.StatusBadRequest, "invalid_json", "invalid json body")
			return
		}

		cat := dto.ToNewCatDomain(*req)

		id, err := h.svc.CreateCat(ctx, &cat)
		if err != nil {
			switch {
			case errors.Is(err, serviceserrors.ErrBreedInvalid):
				httperror.RespondError(c, http.StatusBadRequest, "invalid_breed", "breed is not allowed")
				return
			case errors.Is(err, serviceserrors.ErrExternalService):
				httperror.RespondError(c, http.StatusBadGateway, "external_unavailable", "breed validation service unavailable")
				return
			default:
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
				return
			}
		}
		c.Header("Location", fmt.Sprintf("/cats/%d", id))
		c.JSON(http.StatusCreated, dto.CreateCatResponse{ID: id})
	}
}

// Get a single spy cat
// @Summary      Get a single spy cat
// @Description  The ability to receive information about a single cat
// @Tags         cats
// @Param        id   path int true "ID Cat"
// @Produce      json
// @Success      200 {object} dto.CatResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /cats/{id} [get]
func (h *CatHandler) GetCat() gin.HandlerFunc {
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

		cat, err := h.svc.GetCat(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, serviceserrors.ErrCatNotFound):
				httperror.RespondError(c, http.StatusNotFound, "not_found", "cat not found")
			default:
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			}
			return
		}
		c.JSON(http.StatusOK, dto.ToCatResponse(cat))
	}

}

// List spy cats
// @Summary      List spy cats
// @Description  ability to view the list of cats
// @Tags         cats
// @Param        limit      query int    false "Limit" minimum(1) maximum(200)
// @Param        offset     query int    false "Offset" minimum(0)
// @Param        name       query string false "Filter by name"
// @Param        breed      query string false "Filter by breed"
// @Param        min_years  query int    false "Min experience"
// @Param        max_years  query int    false "Max experience"
// @Produce      json
// @Success      200 {object} dto.GetCatsResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /cats [get]
func (h *CatHandler) GetCats() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var q dto.GetCatsQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			httperror.RespondError(c, http.StatusBadRequest, "invalid_query", err.Error())
			return
		}

		if err := h.validator.Validate(q); err != nil {
			httperror.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
			return
		}

		params := domain.ListCatsParams{
			Name:     q.Name,
			Breed:    q.Breed,
			MinYears: q.MinYears,
			MaxYears: q.MaxYears,
			Limit:    q.Limit,
			Offset:   q.Offset,
		}

		items, err := h.svc.ListCats(ctx, params)
		if err != nil {
			httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			return
		}

		next := 0
		if len(items) == q.Limit {
			next = q.Offset + q.Limit
		}
		c.JSON(http.StatusOK, dto.GetCatsResponse{
			Items:      dto.ToCatResponses(items),
			Limit:      q.Limit,
			Offset:     q.Offset,
			NextOffset: next,
		})

	}
}

// DeleteCat godoc
// @Summary      Delete cat
// @Description  Deletes a cat by id
// @Tags         cats
// @Produce      json
// @Param        id   path      int  true  "Cat ID"
// @Success      200  {object}  dto.DeleteCatResponse "OK"
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /cats/{id} [delete]
func (h *CatHandler) DeleteCat() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idStr := c.Param("id")
		if idStr == "" {
			httperror.RespondError(c, http.StatusBadRequest, "invalid id", "invalid id")
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			httperror.RespondError(c, http.StatusBadRequest, "invalid_path", "id must be a positive integer")
			return
		}

		_, err = h.svc.DeleteCat(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, serviceserrors.ErrCatNotFound):
				httperror.RespondError(c, http.StatusNotFound, "not found", "cat not found")
			default:
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")

			}
			return
		}
		c.JSON(http.StatusOK, dto.DeleteCatResponse{Deleted: true, ID: id})
	}

}

// Update salary
// @Summary      Update cat salary
// @Description  Updates salary for a specific cat
// @Tags         cats
// @Accept       json
// @Produce      json
// @Param        id   path  int true "Cat ID"
// @Param        body body  dto.UpdateSalaryRequest true "New salary"
// @Success      200  {object} dto.CatResponse
// @Failure      400  {object} dto.ErrorResponse
// @Failure      404  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /cats/{id}/salary [patch]
func (h *CatHandler) UpdateSalary() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil || id <= 0 {
			httperror.RespondError(c, http.StatusBadRequest, "invalid_path", "id must be a positive integer")
			return
		}

		req, err := validator.DecodeJSON[dto.UpdateSalaryRequest](h.validator, c.Request)
		if err != nil {
			if errors.Is(err, validator.ErrHandlerValidationFailed) {
				httperror.RespondError(c, http.StatusBadRequest, "invalid_body", err.Error())
				return
			}
			httperror.RespondError(c, http.StatusBadRequest, "invalid_json", "invalid json body")
			return
		}

		cat, err := h.svc.UpdateSalary(ctx, domain.UpdateSalaryParams{
			ID:     id,
			Salary: req.Salary,
		})

		if err != nil {
			switch {
			case errors.Is(err, serviceserrors.ErrCatNotFound):
				httperror.RespondError(c, http.StatusNotFound, "not_found", "cat not found")
			case errors.Is(err, serviceserrors.ErrInvalidSalary):
				httperror.RespondError(c, http.StatusBadRequest, "invalid_salary", "salary must be >= 0")
			default:
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			}
			return
		}

		c.JSON(http.StatusOK, dto.ToCatResponse(cat))
	}
}
