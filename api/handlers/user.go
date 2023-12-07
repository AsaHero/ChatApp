package handlers

import (
	"time"

	"github.com/AsaHero/chat_app/pkg/config"
	"github.com/AsaHero/chat_app/service"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type userHandler struct {
	BaseHandler
	cfg     *config.Config
	service service.User
}

func NewUserHandler(service service.User, cfg *config.Config) *fiber.App {
	handler := userHandler{
		cfg:     cfg,
		service: service,
	}

	app := fiber.New()
	app.Get("/", handler.GetUser)

	return app
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
	data, ok := h.GetCtxData(c)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"err": "failed to parse context data"})
	}

	id, ok := data["id"]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"err": "user id is not identified in authorization token"})
	}

	user, err := h.service.Get(c.Context(), map[string]string{
		"id": id,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"err": err.Error()})
	}

	resp := User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC1123),
		UpdatedAt: user.UpdatedAt.Format(time.RFC1123),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
