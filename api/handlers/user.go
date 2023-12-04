package handlers

import (
	"encoding/json"

	"github.com/AsaHero/chat_app/entity"
	"github.com/AsaHero/chat_app/pkg/config"
	"github.com/AsaHero/chat_app/pkg/token"
	"github.com/AsaHero/chat_app/service"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userHandler struct {
	cfg     *config.Config
	service service.User
}

func NewUserHandler(service service.User, cfg *config.Config) *fiber.App {
	handler := userHandler{
		cfg:     cfg,
		service: service,
	}

	app := fiber.New()

	app.Post("/sign-up", handler.SignUp)
	app.Post("/login", handler.Login)

	return app
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := json.Unmarshal(c.Request().Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "cannot parse request!"})
	}

	user, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
	}

	token, err := token.GenerateJWT(h.cfg.Token.Secret, map[string]any{"id": user.ID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *userHandler) SignUp(c *fiber.Ctx) error {
	var req SignUpRequest
	if err := json.Unmarshal(c.Request().Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "cannot parse request!"})
	}

	var isExists bool = true

	_, err := h.service.Get(c.Context(), map[string]string{"email": req.Email})
	if err != nil {
		if err == entity.ErrorNotFound {
			isExists = false
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}
	}
	if isExists {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "user already exists!"})
	}

	id, err := h.service.Create(c.Context(), &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}
