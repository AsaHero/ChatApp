package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/AsaHero/chat_app/entity"
	"github.com/AsaHero/chat_app/pkg/config"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQL struct {
	Pool *pgxpool.Pool
	Sq   sq.StatementBuilderType
}

func NewPostgreSQLDatabase(cfg *config.Config) (*PostgreSQL, error) {
	dbConf, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", cfg.PostgresDB.User, cfg.PostgresDB.Password, cfg.PostgresDB.Host, cfg.PostgresDB.Port, cfg.PostgresDB.DBName, cfg.PostgresDB.SSLMode))
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), dbConf)
	if err != nil {
		return nil, err
	}

	return &PostgreSQL{
		Pool: db,
		Sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}, nil
}

func (p *PostgreSQL) Close() {
	p.Pool.Close()
}

func (p *PostgreSQL) Error(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return entity.ErrorConflict
		}
	}

	if err.Error() == pgx.ErrNoRows.Error() {
		return entity.ErrorNotFound
	}
	return err
}
