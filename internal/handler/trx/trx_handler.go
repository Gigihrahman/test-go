package trx_handler

import (
	"net/http"
	"strconv"

	"test-rakamin/internal/models"
	trx_service "test-rakamin/internal/service/trx"
	"test-rakamin/utils"
	"test-rakamin/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

type TrxHandler interface {
	RegisterRoutes(app *fiber.App)
	GetAllTrx(c *fiber.Ctx) error
	GetTrxByID(c *fiber.Ctx) error
	CreateTrx(c *fiber.Ctx) error
}

type trxHandlerImpl struct {
	trxService trx_service.TrxService
}

func NewTrxHandler(service trx_service.TrxService) TrxHandler {
	return &trxHandlerImpl{trxService: service}
}

func (h *trxHandlerImpl) RegisterRoutes(app *fiber.App) {
	trxRoutes := app.Group("/api/trx", middleware.JWTMiddleware())
	trxRoutes.Get("/", h.GetAllTrx)
	trxRoutes.Get("/:id", h.GetTrxByID)
	trxRoutes.Post("/", h.CreateTrx)
}

func (h *trxHandlerImpl) GetAllTrx(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Invalid user token", "User ID not found")
	}
	trxList, err := h.trxService.GetAllTrxByUserID(userID)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to get transactions", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", trxList)
}

func (h *trxHandlerImpl) GetTrxByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Invalid user token", "User ID not found")
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
	}
	trx, err := h.trxService.GetTrxByID(uint(id), userID)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to get transaction", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", trx)
}

func (h *trxHandlerImpl) CreateTrx(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Invalid user token", "User ID not found")
	}

	var payload models.TrxPayload
	if err := c.BodyParser(&payload); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	newTrx, err := h.trxService.CreateTrx(userID, &payload)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to create transaction", err.Error())
	}

	return utils.SuccessResponseFiber(c, http.StatusCreated, "Succeed to POST data", newTrx)
}
