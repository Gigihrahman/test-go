package utils

import (
	"github.com/gofiber/fiber/v2" // Menggunakan Fiber
)

// Response struct untuk konsistensi respons API
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponseFiber mengirim respons sukses untuk Fiber
func SuccessResponseFiber(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// ErrorResponseFiber mengirim respons error untuk Fiber
func ErrorResponseFiber(c *fiber.Ctx, status int, message string, err string) error {
	return c.Status(status).JSON(Response{
		Status:  status,
		Message: message,
		Error:   err,
	})
}
