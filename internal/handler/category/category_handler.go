package category_handler

import (
	"net/http"
	"strconv"

	"test-rakamin/internal/models"
	category_service "test-rakamin/internal/service/category"
	"test-rakamin/utils"
	"test-rakamin/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler interface {
	RegisterRoutes(app *fiber.App)
	GetAllCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type categoryHandlerImpl struct {
	categoryService category_service.CategoryService
}

func NewCategoryHandler(service category_service.CategoryService) CategoryHandler {
	return &categoryHandlerImpl{categoryService: service}
}

func (h *categoryHandlerImpl) RegisterRoutes(app *fiber.App) {
	categoryRoutes := app.Group("/api/category")
	categoryRoutes.Get("/", h.GetAllCategories)
	categoryRoutes.Get("/:id", h.GetCategoryByID)

	authCategoryRoutes := app.Group("/api/category", middleware.JWTMiddleware())
	authCategoryRoutes.Post("/", h.CreateCategory)
	authCategoryRoutes.Put("/:id", h.UpdateCategory)
	authCategoryRoutes.Delete("/:id", h.DeleteCategory)
}

func (h *categoryHandlerImpl) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to get categories", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", categories)
}

func (h *categoryHandlerImpl) GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid category ID", err.Error())
	}
	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to get category", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", category)
}

func (h *categoryHandlerImpl) CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	newCategory, err := h.categoryService.CreateCategory(&category)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to create category", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusCreated, "Succeed to POST data", newCategory)
}

func (h *categoryHandlerImpl) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid category ID", err.Error())
	}
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	updatedCategory, err := h.categoryService.UpdateCategory(uint(id), &category)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to update category", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to PUT data", updatedCategory)
}

func (h *categoryHandlerImpl) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid category ID", err.Error())
	}
	err = h.categoryService.DeleteCategory(uint(id))
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to delete category", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to DELETE data", nil)
}
