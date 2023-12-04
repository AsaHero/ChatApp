package api

import (
	"github.com/AsaHero/chat_app/api/handlers"
	"github.com/AsaHero/chat_app/pkg/config"
	"github.com/AsaHero/chat_app/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type RouterArgs struct {
	Cfg         *config.Config
	UserService service.User
}

func NewRouter(args RouterArgs) *fiber.App {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(favicon.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-Token,X-Request-Id",
		AllowCredentials: true,
		ExposeHeaders:    "Link",
		MaxAge:           300,
	}))

	app.Route("/", func(router fiber.Router) {
		app.Mount("/user", handlers.NewUserHandler(args.UserService, args.Cfg))
	})

	return app
}
