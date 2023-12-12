package api

import (
	"github.com/AsaHero/chat_app/api/handlers"
	"github.com/AsaHero/chat_app/api/middleware"
	"github.com/AsaHero/chat_app/api/ws"
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
	Hub         *ws.Hub
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

	public := app.Group("/")
	public.Mount("/", handlers.NewAuthHandler(args.UserService, args.Cfg))

	protected := app.Group("/", middleware.Authorizer(args.Cfg.Token.Secret))
	protected.Mount("/user", handlers.NewUserHandler(args.UserService, args.Cfg))
	protected.Mount("/ws", ws.NewHandler(args.Hub, args.UserService))

	return app
}
