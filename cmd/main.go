package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/AsaHero/chat_app/api"
	"github.com/AsaHero/chat_app/pkg/config"
	"github.com/AsaHero/chat_app/pkg/db/postgresql"
	"github.com/AsaHero/chat_app/repository"
	"github.com/AsaHero/chat_app/service"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.NewConfig()

	timeDuration, err := time.ParseDuration("30s")
	if err != nil {
		log.Fatalf("error while duration parse: %s", err.Error())
	}

	db, err := postgresql.NewPostgreSQLDatabase(cfg)
	if err != nil {
		log.Fatalf("error while db init: %s", err.Error())
	}

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(timeDuration, userRepo)

	router := api.NewRouter(api.RouterArgs{
		Cfg:         cfg,
		UserService: userService,
	})

	go func() {
		if err := router.Listen(fmt.Sprintf("%s%s", cfg.Host, cfg.Port)); err != nil {
			log.Fatalf("error while starting server: %s", err.Error())
		}
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	<-sign

	log.Println("Server stopped!")
	router.Shutdown()
	db.Close()
}
