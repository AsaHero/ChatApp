package handlers

import (
	"github.com/AsaHero/chat_app/service"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	service service.User
}

func NewUserHandler(service service.User) *fiber.App {
	handler := userHandler{}

	app := fiber.New()

	app.Post("/sign-up", handler.SignUp)
	app.Post("/login", handler.Login)

	return app
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": "1"})
}

func (h *userHandler) SignUp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": "1"})
}
