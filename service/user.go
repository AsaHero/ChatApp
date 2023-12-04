package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AsaHero/chat_app/entity"
	"github.com/AsaHero/chat_app/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) (string, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error

	Login(ctx context.Context, email, password string) (*entity.User, error)
}

type userService struct {
	ctxTimeout time.Duration
	userRepo   repository.User
}

func NewUserService(ctxTimeout time.Duration, userRepo repository.User) User {
	return &userService{
		ctxTimeout: ctxTimeout,
		userRepo:   userRepo,
	}
}

func (userService) beforeCreate(user *entity.User) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (userService) beforeUpdate(user *entity.User) {
	user.UpdatedAt = time.Now()
}

func (s *userService) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.userRepo.Get(ctx, params)
}
func (s *userService) Create(ctx context.Context, user *entity.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error while hashing the password: %s", err.Error())
		return "", err
	}

	user.PasswordHash = string(hash)

	s.beforeCreate(user)

	return user.ID, s.userRepo.Create(ctx, user)
}
func (s *userService) Update(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error while hashing the password: %s", err.Error())
		return err
	}

	user.PasswordHash = string(hash)

	s.beforeUpdate(user)

	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.userRepo.Delete(ctx, id)
}

func (s *userService) Login(ctx context.Context, email, password string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	user, err := s.Get(ctx, map[string]string{"email": email})
	if err != nil {
		if err == entity.ErrorNotFound {
			return nil, fmt.Errorf("user doesn't exists!")
		} else {
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}
