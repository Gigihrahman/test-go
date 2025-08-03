package toko_handler

import (
	"log"
	"net/http"
	"strconv"
	toko_service "test-rakamin/internal/service/toko"
	"test-rakamin/utils"
	"test-rakamin/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

type TokoHandler interface {
	RegisterRoutes(app *fiber.App)
	GetMyToko(c *fiber.Ctx) error
	GetAllToko(c *fiber.Ctx) error
	GetTokoByID(c *fiber.Ctx) error
	UpdateToko(c *fiber.Ctx) error
}

type tokoHandlerImpl struct {
	tokoService toko_service.TokoService
}

func NewTokoHandler(service toko_service.TokoService) TokoHandler {
	return &tokoHandlerImpl{tokoService: service}
}

func (h *tokoHandlerImpl) RegisterRoutes(app *fiber.App) {

	tokoRoutes := app.Group("/api/toko")
	tokoRoutes.Get("/", h.GetAllToko)
	tokoRoutes.Get("/:id_toko", h.GetTokoByID)

	authTokoRoutes := app.Group("/api/toko", middleware.JWTMiddleware())
	authTokoRoutes.Get("/my", h.GetMyToko)
	authTokoRoutes.Put("/:id_toko", h.UpdateToko)
}

func (h *tokoHandlerImpl) GetMyToko(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Invalid user token", "User ID not found in token")
	}
	toko, err := h.tokoService.GetTokoByUserID(userID)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to get toko", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", toko)
}

func (h *tokoHandlerImpl) GetAllToko(c *fiber.Ctx) error {

	tokoList, err := h.tokoService.GetAllToko()
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to get toko list", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", tokoList)
}

func (h *tokoHandlerImpl) GetTokoByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id_toko"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid toko ID", err.Error())
	}
	toko, err := h.tokoService.GetTokoByID(uint(id))
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to get toko", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", toko)
}

func (h *tokoHandlerImpl) UpdateToko(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id_toko"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid toko ID", err.Error())
	}

	namaToko := c.FormValue("nama_toko")
	file, err := c.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("Failed to get form file: %v", err)
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Failed to parse form file", err.Error())
	}

	updatedToko, err := h.tokoService.UpdateToko(uint(id), namaToko, file)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to update toko", err.Error())
	}

	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to UPDATE data", updatedToko)
}
