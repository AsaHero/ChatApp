package repository

import (
	"context"
	"fmt"

	"github.com/AsaHero/chat_app/entity"
	"github.com/AsaHero/chat_app/pkg/db/postgresql"
	"github.com/Masterminds/squirrel"
)

const (
	UserTableName = "users"
)

type User interface {
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}

type userRepo struct {
	db        postgresql.PostgreSQL
	tableName string
}

func NewUserRepo(db *postgresql.PostgreSQL) User {
	return &userRepo{
		db:        *db,
		tableName: UserTableName,
	}
}

func (r *userRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	builder := r.db.Sq.Select(
		"id",
		"username",
		"email",
		"password",
		"created_at",
		"updated_at",
	).From(r.tableName)

	for k, v := range params {
		switch k {
		case "id":
			builder = builder.Where(squirrel.Eq{"id": v})
		case "email":
			builder = builder.Where(squirrel.Eq{"email": v})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error while building sql query: %s", err.Error())
	}

	fmt.Println(query, args)

	var user entity.User
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		fmt.Println(err)
		return nil, r.db.Error(err)
	}

	return &user, nil
}
func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	query, args, err := r.db.Sq.Insert(r.tableName).SetMap(map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"password":   user.PasswordHash,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}).ToSql()
	if err != nil {
		return fmt.Errorf("error while building sql query: %s", err.Error())
	}

	fmt.Println(query, args)

	if _, err = r.db.Pool.Exec(ctx, query, args...); err != nil {
		return r.db.Error(err)
	}

	return nil
}
func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	query, args, err := r.db.Sq.Update(r.tableName).SetMap(map[string]interface{}{
		"username":   user.Username,
		"email":      user.Email,
		"password":   user.PasswordHash,
		"updated_at": user.UpdatedAt,
	}).Where(squirrel.Eq{"id": user.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("error while building sql query: %s", err.Error())
	}
	if _, err = r.db.Pool.Exec(ctx, query, args...); err != nil {
		return r.db.Error(err)
	}

	return nil
}
func (r *userRepo) Delete(ctx context.Context, id string) error {
	query, args, err := r.db.Sq.Delete(r.tableName).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("error while building sql query: %s", err.Error())
	}
	if _, err = r.db.Pool.Exec(ctx, query, args...); err != nil {
		return r.db.Error(err)
	}

	return nil
}
