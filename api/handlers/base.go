package handlers

import (
	"github.com/AsaHero/chat_app/api/middleware"
	"github.com/gofiber/fiber/v2"
)

type BaseHandler struct{}

func (BaseHandler) GetCtxData(c *fiber.Ctx) (map[string]string, bool) {
	data, ok := c.UserContext().Value(middleware.CtxKeyValue).(map[string]string)
	return data, ok
}
