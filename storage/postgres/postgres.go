package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"task_service/configs"
	"task_service/pkg/logger"
	"task_service/storage"
)

type Store struct {
	DB  *pgxpool.Pool
	log logger.ILogger
	cfg configs.Config
}

func NewStore(ctx context.Context, log logger.ILogger, cnf configs.Config) (*Store, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cnf.PostgresUser,
		cnf.PostgresPassword,
		cnf.PostgresHost,
		cnf.PostgresPort,
		cnf.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		fmt.Println("++++++++++", err)
		log.Error("this error is parse url -> can not parsing", logger.Error(err))
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Println("++++++++++", err)
		log.Error("this error ie new create config with pool", logger.Error(err))
		return nil, err
	}
	return &Store{
		DB:  pool,
		log: log,
		cfg: cnf,
	}, nil
}
func (s Store) Tasks() storage.ITaskStorage {
	return NewTaskRepository(s.DB, s.log)
}
