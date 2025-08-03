package user_handler

import (
	"net/http"

	"test-rakamin/internal/models"
	user_service "test-rakamin/internal/service/user"
	"test-rakamin/utils"
	"test-rakamin/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetMyProfile(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	RegisterRoutes(app *fiber.App)
}

type userHandlerImpl struct {
	userService user_service.UserService
}

func NewUserHandler(service user_service.UserService) UserHandler {
	return &userHandlerImpl{userService: service}
}

func (h *userHandlerImpl) RegisterRoutes(app *fiber.App) {

	authRoutes := app.Group("/api/auth")
	authRoutes.Post("/register", h.Register)
	authRoutes.Post("/login", h.Login)

	userRoutes := app.Group("/api/user", middleware.JWTMiddleware())
	userRoutes.Get("/", h.GetMyProfile)
	userRoutes.Put("/", h.UpdateProfile)
}

func (h *userHandlerImpl) Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	createdUser, err := h.userService.RegisterUser(&user)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to register user", err.Error())
	}

	return utils.SuccessResponseFiber(c, http.StatusCreated, "Register Succeed", createdUser)
}

func (h *userHandlerImpl) Login(c *fiber.Ctx) error {
	var loginPayload struct {
		NoTelp    string `json:"no_telp"`
		KataSandi string `json:"kata_sandi"`
	}
	if err := c.BodyParser(&loginPayload); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	token, err := h.userService.LoginUser(loginPayload.NoTelp, loginPayload.KataSandi)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Failed to login", err.Error())
	}

	user, _ := h.userService.GetUserProfile(1)

	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to POST data", fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (h *userHandlerImpl) GetMyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Invalid user token", "User ID not found in token")
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusNotFound, "Failed to get user profile", err.Error())
	}

	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to GET data", user)
}

func (h *userHandlerImpl) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Invalid user token", "User ID not found in token")
	}

	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return utils.ErrorResponseFiber(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	user, err := h.userService.UpdateUserProfile(userID, &updatedUser)
	if err != nil {
		return utils.ErrorResponseFiber(c, http.StatusInternalServerError, "Failed to update user profile", err.Error())
	}

	return utils.SuccessResponseFiber(c, http.StatusOK, "Succeed to UPDATE data", user)
}
