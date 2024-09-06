package mongo

import (
	"context"
	"fmt"
	"task_service/configs"
	"task_service/pkg/logger"
	"task_service/storage"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client *mongo.Client
	db     *mongo.Database
	cfg    configs.Config
	log    logger.ILogger
}

func NewStore(ctx context.Context, cfg configs.Config, log logger.ILogger) (*Store, error) {
    log.Info("Attempting to connect to MongoDB", logger.String("host", cfg.MongoHost), logger.String("db", cfg.MongoDB))
	log.Info(cfg.MongoDB)
    log.Info(cfg.MongoHost)	
	log.Info(cfg.MongoDB)
	log.Info(cfg.MongoPassword)
	log.Info(cfg.MongoUser)
    uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-1",
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort,
		cfg.MongoDB,
	)
	

    log.Info("Connecting to MongoDB", logger.String("uri", uri))

    clientOptions := options.Client().ApplyURI(uri).
        SetMaxPoolSize(100).
        SetConnectTimeout(10 * time.Second)

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Error("Failed to connect to MongoDB", logger.Error(err))
        return nil, err
    }

    defer func() {
        if err := client.Disconnect(ctx); err != nil {
            log.Error("Failed to disconnect MongoDB client", logger.Error(err))
        }
    }()

    if err := client.Ping(ctx, nil); err != nil { 
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
