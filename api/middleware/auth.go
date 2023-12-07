package middleware

import (
	"context"

	"github.com/AsaHero/chat_app/pkg/token"
	"github.com/gofiber/fiber/v2"
)

func Authorizer(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwtToken := c.Get("Authorization")

		if len(jwtToken) < 10 {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"err": "token is invalid"})
		}

		if jwtToken[:6] != "Bearer" || jwtToken[7:] == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"err": "token is invalid (hint: add 'Bearer' at the beginning)"})
		}

		jwtToken = jwtToken[7:]

		claims, err := token.ParseJWT(jwtToken, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"err": err.Error()})
		}

		id := claims["id"].(string)

		c.SetUserContext(context.WithValue(c.Context(), CtxKeyValue, map[string]string{"id": id}))

		return c.Next()
	}
}
