package database

import (
	"context"
	"fmt"

	"github.com/SamsonGedefa/simulator/main.go/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg config.Config) (*pgxpool.Pool, error) {
	pgxConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPassword, cfg.PgDB, cfg.PgSSL)

	connPool, err := pgxpool.New(context.Background(), pgxConnString)

	if err != nil {
		panic(err)
	}

	return connPool, nil

}
