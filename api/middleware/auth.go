package middleware

import (
	"context"

	"github.com/AsaHero/chat_app/pkg/token"
	"github.com/gofiber/fiber/v2"
)

func Authorizer(secret string, c *fiber.Ctx) error {
	jwtToken := c.GetRespHeader("Authorization", "")

	if len(jwtToken) < 10 {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	if jwtToken[:6] != "Bearer" || jwtToken[7:] == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	jwtToken = jwtToken[7:]

	claims, err := token.ParseJWT(jwtToken, secret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"err": err.Error()})
	}

	id := claims["id"].(string)


	c.SetUserContext(context.WithValue(c.Context(), CtxKeyValue, id))	

	return nil
}
