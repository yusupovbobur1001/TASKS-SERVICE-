package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"task_service/configs"
	"task_service/pkg/logger"
	"task_service/storage"
	"time"
)

type Store struct {
	client *mongo.Client
	db     *mongo.Database
	cfg    configs.Config
	log    logger.ILogger
}

func NewStore(ctx context.Context, cfg configs.Config, log logger.ILogger) (*Store, error) {
	uri := fmt.Sprintf(
		`mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-256`,
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort,
		cfg.MongoDB,
	)

	log.Info("Attempting to connect to MongoDB", logger.String("uri", uri))

	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("Failed to connect to MongoDB", logger.Error(err))
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("Failed to ping MongoDB", logger.Error(err))
		return nil, err
	}

	log.Info("Successfully connected to MongoDB")

	db := client.Database(cfg.MongoDB)

	return &Store{
		client: client,
		db:     db,
		cfg:    cfg,
		log:    log,
	}, nil
}

func (s *Store) Close() {
	if err := s.client.Disconnect(context.Background()); err != nil {
		s.log.Error("Error disconnecting from MongoDB", logger.Error(err))
	}
}
func (s Store) Details() storage.IDetailStorage {
	return NewDetailMongo(s.db.Collection("detail"), s.log)
}
