package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/SamsonGedefa/simulator/main.go/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

func NewPGX(cfg config.Config) (*postgres, error) {
	pgxConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPassword, cfg.PgDB, cfg.PgSSL)

	pgOnce.Do(func() {
		db, err := pgxpool.New(context.Background(), pgxConnString)
		if err != nil {
			panic(err)
		}

		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}
