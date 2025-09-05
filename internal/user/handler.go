package user

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, svc Service) {
	app.Post("/users", func(c *fiber.Ctx) error {
		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
		}
		resp, err := svc.Create(c.Context(), req)
		if err != nil {
			switch err {
			case ErrBadEmail, ErrBadPassword, ErrBadName:
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			case ErrDuplicateEmail:
				return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "email exists"})
			default:
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
			}
		}
		return c.Status(http.StatusCreated).JSON(resp)
	})
}
