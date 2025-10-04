package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	servieserrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
	"github.com/gin-gonic/gin"

	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/dto"
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
// @Tags Cats
// @Description Used to create a new spy cat
// @Accept json
// @Produce json
// @Param CreateCatRequest body  dto.CreateCatRequest true "Request to create a cat"
// @Success 201 {object}  dto.CreateCatResponse
// @Failure 400
// @Failure 500
// @Router /cats/create [post]
func (h *CatHandler) CreateCat() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		req, err := validator.DecodeJSON[dto.CreateCatRequest](h.validator, c.Request)
		if err != nil {
			if errors.Is(err, validator.ErrHandlerValidationFailed) {
				httperror.RespondError(c, http.StatusBadRequest, "invalid_body", err.Error())
			}
			httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			return
		}

		id, err := h.svc.CreateCat(ctx, &domain.Cat{
			Name:            req.Name,
			YearsExperience: req.YearsExperience,
			Breed:           req.Breed,
			Salary:          req.Salary,
		})

		if err != nil {
			httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			return
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
// @Success      200 {object} domain.Cat
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
			case errors.Is(err, servieserrors.ErrCatNotFound):
				httperror.RespondError(c, http.StatusNotFound, "not_found", "cat not found")
			default:
				httperror.RespondError(c, http.StatusInternalServerError, "internal", "internal server error")
			}
			return
		}
		c.Header("Location", fmt.Sprintf("/cats/%d", cat.ID))
		c.JSON(http.StatusOK, cat)
	}

}

// List spy cats
// @Summary      List spy cats
// @Description  Возможность просматривать список котов
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
// func (h *CatHandler) GetCats() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()

// 		// 1) Парсим и валидируем query
// 		q, err := validator.DecodeQuery[dto.GetCatsQuery](h.validator, c.Request)
// 		if err != nil {
// 			if errors.Is(err, validator.ErrHandlerValidationFailed) {
// 				respondError(c, http.StatusBadRequest, "invalid_query", err.Error())
// 				return
// 			}
// 			respondError(c, http.StatusInternalServerError, "internal", "internal server error")
// 			return
// 		}
// 		// дефолты (если не выставляешь их в DecodeQuery)
// 		if q.Limit == 0 {
// 			q.Limit = 50
// 		}
// 		// 2) Готовим параметры домена
// 		params := domain.ListCatsParams{
// 			Name:     q.Name,
// 			Breed:    q.Breed,
// 			MinYears: q.MinYears,
// 			MaxYears: q.MaxYears,
// 			Limit:    q.Limit,
// 			Offset:   q.Offset,
// 		}
// 		// 3) Вызываем сервис
// 		items, total, err := h.svc.ListCats(ctx, params)
// 		if err != nil {
// 			respondError(c, http.StatusInternalServerError, "internal", "internal server error")
// 			return
// 		}
// 		// 4) Ответ (используем domain.Cat напрямую, как ты хотел)
// 		c.JSON(http.StatusOK, dto.GetCatsResponse{
// 			Items:      items,
// 			Total:      total,
// 			Limit:      q.Limit,
// 			Offset:     q.Offset,
// 			NextOffset: q.Offset + len(items),
// 		})
// 	}
// }
