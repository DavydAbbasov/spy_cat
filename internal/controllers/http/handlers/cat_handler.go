package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	"github.com/DavydAbbasov/spy-cat/internal/service"
	"github.com/gin-gonic/gin"

	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/dto"
	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"
)

type CatHandler struct {
	svc       service.CatService
	validator *validator.Validator
}

func NewCatHandler(svc service.CatService, validator *validator.Validator) *CatHandler {
	return &CatHandler{
		svc:       svc,
		validator: validator,
	}
}

// CreateCatHandler godoc
// @Summary      Создать Cat
// @Description  Добавляет нового кота в систему
// @Tags         cat
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateCatRequest   true  "Данные кота"
// @Success      201    {object}  dto.CreateCatResponse
// @Failure      400    {object}  renderer.BaseExceptions
// @Failure      500    {object}  renderer.BaseExceptions
// @Router       /cats [post]
func (h *CatHandler) CreateCatHandler(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := validator.DecodeJSON[dto.CreateCatRequest](h.validator, c.Request)
	if err != nil {
		if errors.Is(err, validator.ErrHandlerValidationFailed) {
			respondError(c, http.StatusBadRequest, "invalid_body", err.Error())
		}
		respondError(c, http.StatusInternalServerError, "internal", "internal server error")
		return
	}

	id, err := h.svc.CreateCat(ctx, &domain.Cat{
		Name:            req.Name,
		YearsExperience: req.YearsExperience,
		Breed:           req.Breed,
		Salary:          req.Salary,
	})

	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal", "internal server error")
		return
	}

	c.Header("Location", fmt.Sprintf("/cats/%d", id))
	c.JSON(http.StatusCreated, dto.CreateCatResponse{ID: id})
}
