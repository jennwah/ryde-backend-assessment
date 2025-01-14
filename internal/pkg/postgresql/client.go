package postgresql

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/jennwah/ryde-backend-engineer/internal/config"
)

type Client struct {
	*sqlx.DB
}

const (
	maxOpenConnection     = 50
	maxIdleConnection     = 25
	maxConnectionLifeTime = 10 * time.Minute
)

func New(cfg config.Postgres) (*Client, error) {
	connString := "user=%s password=%s host=%s port=%d dbname=%s search_path=%s sslmode=disable"
	dbDSN := fmt.Sprintf(connString,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SchemaName,
	)

	db, err := sqlx.Connect("pgx", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("connect to posgresql: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConnection)
	db.SetMaxIdleConns(maxIdleConnection)
	db.SetConnMaxLifetime(maxConnectionLifeTime)

	return &Client{
		DB: db,
	}, nil
}
