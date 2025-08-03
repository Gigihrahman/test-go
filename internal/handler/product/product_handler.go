package product_handler

import (
	"net/http"
	"strconv"

	"test-rakamin/internal/models"
	product_service "test-rakamin/internal/service/product"
	"test-rakamin/utils"
	"test-rakamin/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler interface {
	RegisterRoutes(app *fiber.App)
	GetAllProducts(c *fiber.Ctx) error
	GetProductByID(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productHandlerImpl struct {
	productService product_service.ProductService
}

func NewProductHandler(service product_service.ProductService) ProductHandler {
	return &productHandlerImpl{productService: service}
}

func (h *productHandlerImpl) RegisterRoutes(app *fiber.App) {
	productRoutes := app.Group("/api/product")
	productRoutes.Get("/", h.GetAllProducts)
	productRoutes.Get("/:id", h.GetProductByID)

	authProductRoutes := app.Group("/api/product", middleware.JWTMiddleware())
	authProductRoutes.Post("/", h.CreateProduct)
	authProductRoutes.Put("/:id", h.UpdateProduct)
	authProductRoutes.Delete("/:id", h.DeleteProduct)
}

func (h *productHandlerImpl) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts(
		c.Query("nama_produk"),
		c.Query("category_id"),
		c.Query("toko_id"),
		c.Query("min_harga"),
		c.Query("max_harga"),
	)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to get products", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", products)
}

func (h *productHandlerImpl) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid product ID", err.Error())
	}
	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to get product", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", product)
}

func (h *productHandlerImpl) CreateProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid form data", err.Error())
	}

	productPayload := models.Product{}
	if err := c.BodyParser(&productPayload); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid product payload", err.Error())
	}

	photos := form.File["photos"]

	newProduct, err := h.productService.CreateProduct(&productPayload, photos)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to create product", err.Error())
	}

	return utils.SuccessResponseFiber(c, http.StatusCreated, "Succeed to POST data", newProduct)
}

func (h *productHandlerImpl) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid product ID", err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid form data", err.Error())
	}

	productPayload := models.Product{}
	if err := c.BodyParser(&productPayload); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid product payload", err.Error())
	}

	photos := form.File["photos"]

	updatedProduct, err := h.productService.UpdateProduct(uint(id), &productPayload, photos)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to update product", err.Error())
	}

	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to UPDATE data", updatedProduct)
}

func (h *productHandlerImpl) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid product ID", err.Error())
	}
	err = h.productService.DeleteProduct(uint(id))
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to delete product", err.Error())
	}
	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to DELETE data", nil)
}
